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
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	if count, err := w.body.Write(b); err != nil {
		return count, err
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	if _, err := w.body.WriteString(s); err != nil {
		return 0, err
	}
	return w.ResponseWriter.WriteString(s)
}

// RequestBodySanitizer returns the representation of a request body to include
// in logs. Returning nil omits the request-body field.
type RequestBodySanitizer func(c *gin.Context, body []byte) ([]byte, error)

// ResponseBodySanitizer returns the representation of a response body to
// include in logs. Returning nil omits the response-body field.
type ResponseBodySanitizer func(c *gin.Context, body []byte) ([]byte, error)

type config struct {
	logRequestBody           bool
	logResponseBody          bool
	excludedBodyPathPrefixes []string
	requestBodySanitizer     RequestBodySanitizer
	responseBodySanitizer    ResponseBodySanitizer
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

// New returns middleware that records request and response access details.
func New(log logger.Logger, options ...Option) gin.HandlerFunc {
	config := config{
		logRequestBody:  false,
		logResponseBody: false,
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
			requestBody, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Error("failed to read request body", logzap.Error(err))
			} else {
				requestLogBody = requestBody
				if config.requestBodySanitizer != nil {
					requestLogBody, err = config.requestBodySanitizer(c, requestBody)
					if err != nil {
						log.Warn("failed to sanitize request body", logzap.Error(err))
					}
				}
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
		}

		var responseWriter *responseBodyWriter
		if config.shouldLogResponseBody(c.Request.URL.Path) {
			responseWriter = &responseBodyWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
			}
			c.Writer = responseWriter
		}

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)
		fields := []logzap.Field{
			logzap.String("time", endTime.UTC().Format(time.RFC3339)),
			logzap.String("method", c.Request.Method),
			logzap.String("host", c.Request.Host),
			logzap.String("path", c.Request.URL.Path),
			logzap.String("query", c.Request.URL.RawQuery),
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
		}
		if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
			fields = append(fields, logzap.String("request-id", requestID))
		}
		if config.shouldLogRequestBody(c.Request.URL.Path) && requestLogBody != nil {
			fields = append(fields, logzap.ByteString("request-body", requestLogBody))
		}
		if responseWriter != nil {
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
		if len(c.Errors) > 0 {
			fields = append(fields, logzap.String("errors", c.Errors.String()))
			log.Error(c.Request.URL.Path, fields...)
			return
		}
		log.Info(c.Request.URL.Path, fields...)
	}
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
