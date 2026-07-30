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
	trustedOrigins      []string
}

// Checker validates CSRF tokens using options prepared at startup.
type Checker struct {
	options          *options
	originProtection *http.CrossOriginProtection
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

// WithTrustedOrigins allows credentialed requests from the listed origins.
// Origins must include an http or https scheme and must not contain a path.
func WithTrustedOrigins(origins ...string) Option {
	return optionFunc(func(opts *options) {
		opts.trustedOrigins = append(opts.trustedOrigins, origins...)
	})
}

// NewChecker prepares reusable CSRF validation state.
func NewChecker(opts ...Option) (*Checker, error) {
	o := newOptions(opts...)
	originProtection := http.NewCrossOriginProtection()
	for _, origin := range o.trustedOrigins {
		if err := originProtection.AddTrustedOrigin(origin); err != nil {
			return nil, fmt.Errorf("csrf: add trusted origin: %w", err)
		}
	}
	return &Checker{options: o, originProtection: originProtection}, nil
}

func Middleware(opts ...Option) gin.HandlerFunc {
	checker, err := NewChecker(opts...)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		if !checker.Check(c) {
			return
		}
		c.Next()
	}
}

// Check validates CSRF for the current request without continuing the Gin chain.
func (checker *Checker) Check(c *gin.Context) bool {
	return check(c, checker.options, checker.originProtection)
}

func check(c *gin.Context, o *options, originProtection *http.CrossOriginProtection) bool {
	if isSafeMethod(c.Request.Method) {
		return true
	}
	if err := originProtection.Check(c.Request); err != nil {
		o.reject(c)
		return false
	}

	session := sessions.Default(c)
	expected, _ := session.Get(o.sessionKey).(string)
	actual := c.PostForm(o.fieldName)
	if actual == "" && o.headerName != "" {
		actual = c.GetHeader(o.headerName)
	}

	if expected == "" || actual == "" || subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) != 1 {
		o.reject(c)
		return false
	}
	return true
}

func (o *options) reject(c *gin.Context) {
	o.unauthorizedHandler(c)
	if !c.IsAborted() {
		c.Abort()
	}
}

func Token(c *gin.Context, opts ...Option) string {
	o := newOptions(opts...)
	session := sessions.Default(c)
	if token, ok := session.Get(o.sessionKey).(string); ok && token != "" {
		return token
	}

	token, err := generateToken(c, o)
	if err != nil {
		return ""
	}
	if !saveToken(c, o, token) {
		return ""
	}
	return token
}

// RotateToken replaces the current session token. Call it after login,
// privilege elevation, or another security-sensitive session change.
func RotateToken(c *gin.Context, opts ...Option) string {
	o := newOptions(opts...)
	token, err := generateToken(c, o)
	if err != nil {
		return ""
	}
	if !saveToken(c, o, token) {
		return ""
	}
	return token
}

func generateToken(c *gin.Context, o *options) (string, error) {
	token, err := o.tokenGenerator()
	if err != nil {
		_ = c.Error(fmt.Errorf("csrf: generate token: %w", err))
		return "", err
	}
	return token, nil
}

func saveToken(c *gin.Context, o *options, token string) bool {
	session := sessions.Default(c)
	previous := session.Get(o.sessionKey)
	session.Set(o.sessionKey, token)
	if err := session.Save(); err != nil {
		if previous == nil {
			session.Delete(o.sessionKey)
		} else {
			session.Set(o.sessionKey, previous)
		}
		_ = c.Error(fmt.Errorf("csrf: save token: %w", err))
		return false
	}
	return true
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func randomToken() (string, error) {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf[:]), nil
}
