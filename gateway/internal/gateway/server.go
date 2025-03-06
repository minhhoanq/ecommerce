package gateway

import (
	"context"
	"errors"
	"net/http"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/pkg/constants"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultShutdownTimeout = constants.DefaultShutdownTimeout
	_defaultAddr            = constants.DefaultPort
	_defaultReadTimeout     = constants.DefaultReadTimeout
	_defaultWriteTimeout    = constants.DefaultWriteTimeout
)

type Server struct {
	server    *http.Server
	l         logger.Interface
	waitGroup *errgroup.Group
	notify    chan error
	ctx       context.Context
}

func NewServer(handler http.Handler, l logger.Interface, waitGroup *errgroup.Group, ctx context.Context, opts ...Option) {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &Server{
		server:    httpServer,
		l:         l,
		waitGroup: waitGroup,
		notify:    make(chan error, 1),
		ctx:       ctx,
	}

	// custom options
	for _, opt := range opts {
		opt(s)
	}

	s.Start()

	s.Shutdown()
}

func (s *Server) Start() {
	s.waitGroup.Go(func() error {
		s.l.Info("start HTTP server at",
			zap.String("Adr: ", s.server.Addr),
			zap.String("ReadTimeout: ", s.server.ReadTimeout.String()),
			zap.String("WriteTimeout: ", s.server.WriteTimeout.String()))

		err := s.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			s.l.Error("HTTP server faild to serve", zap.String("Error", err.Error()))
			return err
		}

		return nil
	})
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() {
	s.waitGroup.Go(func() error {
		s.l.Info("graceful shutdown HTTP gateway server")
		<-s.ctx.Done()
		err := s.server.Shutdown(context.Background())
		if err != nil {
			s.l.Error("failed to shutdown HTTP server", zap.Error(err))
		}

		s.l.Info("HTTP gateway server is stopped")
		return nil
	})
}
