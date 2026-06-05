package csrf

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miiy/goc/gin/sessions"
)

const (
	FieldName  = "_csrf"
	HeaderName = "X-CSRF-Token"
	SessionKey = "_csrf_token"
)

const contextTokenKeyPrefix = "goc.csrf.token:"

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(opts *options) {
	f(opts)
}

type options struct {
	fieldName           string
	headerName          string
	sessionKey          string
	unauthorizedHandler func(*gin.Context)
	tokenGenerator      func() (string, error)
}

func defaultOptions() *options {
	return &options{
		fieldName:  FieldName,
		headerName: HeaderName,
		sessionKey: SessionKey,
		unauthorizedHandler: func(c *gin.Context) {
			c.AbortWithStatus(http.StatusForbidden)
		},
		tokenGenerator: randomToken,
	}
}

func newOptions(opts ...Option) *options {
	o := defaultOptions()
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func WithFieldName(fieldName string) Option {
	return optionFunc(func(opts *options) {
		if fieldName != "" {
			opts.fieldName = fieldName
		}
	})
}

func WithHeaderName(headerName string) Option {
	return optionFunc(func(opts *options) {
		opts.headerName = headerName
	})
}

func WithSessionKey(sessionKey string) Option {
	return optionFunc(func(opts *options) {
		if sessionKey != "" {
			opts.sessionKey = sessionKey
		}
	})
}

func WithUnauthorized(handler func(*gin.Context)) Option {
	return optionFunc(func(opts *options) {
		if handler != nil {
			opts.unauthorizedHandler = handler
		}
	})
}

func Middleware(opts ...Option) gin.HandlerFunc {
	o := newOptions(opts...)

	return func(c *gin.Context) {
		if isSafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		session := sessions.Default(c)
		expected, _ := session.Get(o.sessionKey).(string)
		actual := c.PostForm(o.fieldName)
		if actual == "" && o.headerName != "" {
			actual = c.GetHeader(o.headerName)
		}

		if expected == "" || actual == "" || subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) != 1 {
			o.unauthorizedHandler(c)
			if !c.IsAborted() {
				c.Abort()
			}
			return
		}

		c.Next()
	}
}

func Token(c *gin.Context, opts ...Option) string {
	o := newOptions(opts...)
	contextKey := contextTokenKey(o.sessionKey)
	if token, ok := c.Get(contextKey); ok {
		if value, ok := token.(string); ok {
			return value
		}
	}

	session := sessions.Default(c)
	if token, ok := session.Get(o.sessionKey).(string); ok && token != "" {
		c.Set(contextKey, token)
		return token
	}

	token, err := o.tokenGenerator()
	if err != nil {
		_ = c.Error(fmt.Errorf("csrf: generate token: %w", err))
		return ""
	}

	session.Set(o.sessionKey, token)
	if err := session.Save(); err != nil {
		_ = c.Error(fmt.Errorf("csrf: save token: %w", err))
		return ""
	}
	c.Set(contextKey, token)
	return token
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

func contextTokenKey(sessionKey string) string {
	return contextTokenKeyPrefix + sessionKey
}

func randomToken() (string, error) {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf[:]), nil
}
