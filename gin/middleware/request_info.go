package middleware

import (
	"bytes"
	"docxlib.com/pkg/log"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func RequestInfo(l *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		reqBody, err := c.GetRawData()
		if err != nil {
			log.Error(err)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		c.Next()

		l.Info(
			"",
			log.String("method", c.Request.Method),
			log.String("host", c.Request.Host),
			log.String("request_uri", c.Request.RequestURI),
			log.String("full_path", c.FullPath()),
			//log.String("header", c.Request.Header.Values()),
			log.String("remote-addr", c.Request.RemoteAddr),
			log.String("user-agent", c.Request.UserAgent()),
			log.ByteString("request-body", reqBody),
			log.String("x-token", c.Request.Header.Get("x-token")),
			log.String("sign", c.Request.Header.Get("sign")),
			log.String("ts", c.Request.Header.Get("ts")),
			log.String("app-version", c.Request.Header.Get("app-version")),
			log.Int("status", c.Writer.Status()),
			log.String("response-body", w.body.String()),
			log.Int("size", c.Writer.Size()),
			log.Int64("latency", time.Since(startTime).Microseconds()),
		)
	}
}
