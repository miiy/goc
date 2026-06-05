package server

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const defaultShutdownTimeout = 5 * time.Second

type Server struct {
	debug           bool
	addr            string
	shutdownTimeout time.Duration
	engine          *gin.Engine
	logger          *zap.Logger
}

func New(opts ...Option) *Server {
	s := &Server{
		debug:           false,
		addr:            "127.0.0.1:8080",
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	if s.debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	s.engine = gin.New()
	return s
}

func (s *Server) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.Use(middleware...)
}

func (s *Server) RegisterRouter(funcs ...func(r *gin.Engine)) {
	for _, f := range funcs {
		f(s.engine)
	}
}

func (s *Server) Handler() http.Handler {
	return s.engine
}

// Run
// https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
func (s *Server) Run(addr ...string) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := s.RunContext(ctx, addr...); err != nil {
		s.zapLogger().Fatal("HTTP server stopped with error", zap.Error(err))
	}
}

func (s *Server) RunContext(ctx context.Context, addr ...string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	address := resolveAddress(addr, s.addr)
	logger := s.logger
	if logger == nil {
		logger = zap.NewNop()
	}
	logger.Info("listening and serving HTTP", zap.String("addr", address))

	srv := &http.Server{
		Addr:    address,
		Handler: s.engine,
	}

	serveErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	logger.Info("shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	if err := <-serveErr; err != nil {
		return err
	}

	logger.Info("HTTP server exited")
	return nil
}

func (s *Server) zapLogger() *zap.Logger {
	if s.logger != nil {
		return s.logger
	}
	return zap.NewNop()
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
