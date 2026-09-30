package accesslog

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miiy/goc/logger"
	logzap "github.com/miiy/goc/logger/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body         *bytes.Buffer
	bodyLimit    int64
	bodyTooLarge bool
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.capture(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.capture([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

func (w *responseBodyWriter) capture(body []byte) {
	if w.bodyLimit <= 0 {
		_, _ = w.body.Write(body)
		return
	}
	remaining := w.bodyLimit - int64(w.body.Len())
	if remaining <= 0 {
		w.bodyTooLarge = w.bodyTooLarge || len(body) > 0
		return
	}
	if int64(len(body)) > remaining {
		_, _ = w.body.Write(body[:remaining])
		w.bodyTooLarge = true
		return
	}
	_, _ = w.body.Write(body)
}

const maxInt64 = int64(^uint64(0) >> 1)

type replayBody struct {
	io.Reader
	io.Closer
}

// RequestBodySanitizer returns the representation of a request body to include
// in logs. Returning nil omits the request-body field.
type RequestBodySanitizer func(c *gin.Context, body []byte) ([]byte, error)

// ResponseBodySanitizer returns the representation of a response body to
// include in logs. Returning nil omits the response-body field.
type ResponseBodySanitizer func(c *gin.Context, body []byte) ([]byte, error)

// PathSanitizer returns the request path representation to include in logs.
type PathSanitizer func(c *gin.Context, path string) string

type config struct {
	logRequestBody           bool
	logResponseBody          bool
	logErrors                bool
	maxBodyLogBytes          int64
	excludedBodyPathPrefixes []string
	requestBodySanitizer     RequestBodySanitizer
	responseBodySanitizer    ResponseBodySanitizer
	pathSanitizer            PathSanitizer
}

// Option configures access logging.
type Option func(*config)

// WithRequestBodyLogging controls request body logging.
func WithRequestBodyLogging(enabled bool) Option {
	return func(config *config) {
		config.logRequestBody = enabled
	}
}

// WithResponseBodyLogging controls response body logging.
func WithResponseBodyLogging(enabled bool) Option {
	return func(config *config) {
		config.logResponseBody = enabled
	}
}

// WithErrorLogging controls whether context errors are attached and logged at Error level.
func WithErrorLogging(enabled bool) Option {
	return func(config *config) {
		config.logErrors = enabled
	}
}

// WithMaxBodyLogBytes omits request or response bodies larger than the supplied limit.
// A non-positive value preserves unlimited body logging for backward compatibility.
func WithMaxBodyLogBytes(maxBytes int64) Option {
	return func(config *config) {
		config.maxBodyLogBytes = maxBytes
	}
}

// WithExcludedBodyPathPrefixes excludes request and response body logging for
// paths beginning with any supplied prefix.
func WithExcludedBodyPathPrefixes(prefixes ...string) Option {
	prefixes = append([]string(nil), prefixes...)
	return func(config *config) {
		config.excludedBodyPathPrefixes = append(config.excludedBodyPathPrefixes, prefixes...)
	}
}

// WithRequestBodySanitizer transforms the request body representation written
// to logs. The original body is always restored for downstream handlers.
func WithRequestBodySanitizer(sanitizer RequestBodySanitizer) Option {
	return func(config *config) {
		config.requestBodySanitizer = sanitizer
	}
}

// WithResponseBodySanitizer transforms the response body representation
// written to logs. The original response is returned to the client unchanged.
func WithResponseBodySanitizer(sanitizer ResponseBodySanitizer) Option {
	return func(config *config) {
		config.responseBodySanitizer = sanitizer
	}
}

// WithPathSanitizer transforms the request path representation written to logs.
func WithPathSanitizer(sanitizer PathSanitizer) Option {
	return func(config *config) {
		config.pathSanitizer = sanitizer
	}
}

// New returns middleware that records request and response access details.
func New(log logger.Logger, options ...Option) gin.HandlerFunc {
	config := config{
		logRequestBody:  false,
		logResponseBody: false,
		logErrors:       true,
		excludedBodyPathPrefixes: []string{
			"/uploads/",
		},
	}
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	return func(c *gin.Context) {
		startTime := time.Now()

		var requestLogBody []byte
		if config.shouldLogRequestBody(c.Request.URL.Path) {
			requestBody, tooLarge, restoredBody, err := captureRequestBody(c.Request.Body, config.maxBodyLogBytes)
			c.Request.Body = restoredBody
			if err != nil {
				log.Error("failed to read request body", logzap.Error(err))
			} else if !tooLarge {
				requestLogBody = requestBody
				if config.requestBodySanitizer != nil {
					requestLogBody, err = config.requestBodySanitizer(c, requestBody)
					if err != nil {
						log.Warn("failed to sanitize request body", logzap.Error(err))
						requestLogBody = nil
					}
				}
			}
		}

		var responseWriter *responseBodyWriter
		if config.shouldLogResponseBody(c.Request.URL.Path) {
			responseWriter = &responseBodyWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
				bodyLimit:      config.maxBodyLogBytes,
			}
			c.Writer = responseWriter
		}

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)
		path := c.Request.URL.Path
		if config.pathSanitizer != nil {
			path = config.pathSanitizer(c, path)
		}
		requestURI := path
		if c.Request.URL.RawQuery != "" {
			requestURI += "?" + c.Request.URL.RawQuery
		}
		fields := []logzap.Field{
			logzap.String("time", endTime.UTC().Format(time.RFC3339)),
			logzap.String("method", c.Request.Method),
			logzap.String("host", c.Request.Host),
			logzap.String("path", path),
			logzap.String("query", c.Request.URL.RawQuery),
			logzap.String("request_uri", requestURI),
			logzap.String("full-path", c.FullPath()),
			logzap.String("ip", c.ClientIP()),
			logzap.String("remote-addr", c.Request.RemoteAddr),
			logzap.String("user-agent", c.Request.UserAgent()),
			logzap.String("authorization", maskAuthorization(c.Request.Header.Get("Authorization"))),
			logzap.String("sign", c.Request.Header.Get("sign")),
			logzap.String("ts", c.Request.Header.Get("ts")),
			logzap.Int("status", c.Writer.Status()),
			logzap.Int("size", c.Writer.Size()),
			logzap.Duration("latency", latency),
			logzap.Int64("latency-us", latency.Microseconds()),
		}
		if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
			fields = append(fields, logzap.String("request-id", requestID))
		}
		if config.shouldLogRequestBody(c.Request.URL.Path) && requestLogBody != nil {
			fields = append(fields, logzap.ByteString("request-body", requestLogBody))
		}
		if responseWriter != nil && !responseWriter.bodyTooLarge {
			responseLogBody := responseWriter.body.Bytes()
			if config.responseBodySanitizer != nil {
				var err error
				responseLogBody, err = config.responseBodySanitizer(c, responseLogBody)
				if err != nil {
					log.Warn("failed to sanitize response body", logzap.Error(err))
					responseLogBody = nil
				}
			}
			if responseLogBody != nil {
				fields = append(fields, logzap.ByteString("response-body", responseLogBody))
			}
		}
		if config.logErrors && len(c.Errors) > 0 {
			fields = append(fields, logzap.String("errors", c.Errors.String()))
			log.Error(path, fields...)
			return
		}
		log.Info(path, fields...)
	}
}

func captureRequestBody(body io.ReadCloser, limit int64) ([]byte, bool, io.ReadCloser, error) {
	if body == nil {
		return nil, false, nil, nil
	}
	var reader io.Reader = body
	if limit > 0 {
		captureLimit := limit
		if captureLimit < maxInt64 {
			captureLimit++
		}
		reader = io.LimitReader(body, captureLimit)
	}
	captured, err := io.ReadAll(reader)
	restored := &replayBody{
		Reader: io.MultiReader(bytes.NewReader(captured), body),
		Closer: body,
	}
	return captured, limit > 0 && int64(len(captured)) > limit, restored, err
}

func (config config) shouldLogRequestBody(path string) bool {
	return config.logRequestBody && !config.excludesBodyPath(path)
}

func (config config) shouldLogResponseBody(path string) bool {
	return config.logResponseBody && !config.excludesBodyPath(path)
}

func (config config) excludesBodyPath(path string) bool {
	for _, prefix := range config.excludedBodyPathPrefixes {
		if prefix != "" && strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// maskAuthorization returns "Bearer ****" for Bearer tokens, or "****" for any non-empty value.
func maskAuthorization(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 {
		return parts[0] + " ****"
	}
	return "****"
}
