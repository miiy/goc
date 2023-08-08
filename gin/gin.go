package gin

import (
	"bytes"
	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type (
	Engine  = gin.Engine
	Context = gin.Context
	H       = gin.H
)

type Server struct {
	debug  bool
	addr   string
	engine *gin.Engine
	logger zap.Logger
}

func New(opts ...Option) *Server {
	s := &Server{
		debug: false,
		addr:  "127.0.0.1:8080",
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	if s.debug {
		gin.SetMode(gin.DebugMode)
		//gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	s.engine = gin.New()
	return s
}

func (s *Server) UseZap(logger *zap.Logger) {
	s.engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	s.engine.Use(ginzap.RecoveryWithZap(logger, true))
}

func (s *Server) UseZapWithConfig(logger *zap.Logger) {
	s.engine.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) (fields []zapcore.Field) {
			// log response ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request-id", requestID))
			}

			// log request body
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			fields = append(fields, zap.String("body", string(body)))

			// log response body

			return
		}),
	}))
}

func (s *Server) RegisterRouter(funcs ...func(r *gin.Engine)) {
	for _, f := range funcs {
		f(s.engine)
	}
}

// Run
// https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
func (s *Server) Run(addr ...string) {
	address := resolveAddress(addr, s.addr)
	s.logger.Info("Listening and serving HTTP on %s\n", zap.String("address", address))

	srv := &http.Server{
		Addr:    address,
		Handler: s.engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server Shutdown:", zap.Error(err))
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		s.logger.Info("timeout of 5 seconds.")
	}
	s.logger.Info("Server exiting")
}

func resolveAddress(addr []string, defaultAddr string) string {
	switch len(addr) {
	case 0:
		return defaultAddr
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
