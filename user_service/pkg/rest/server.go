package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/pkg/constants"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"
)

const (
	_defaultShutdownTimeout = constants.DefaultShutdownTimeout
	_defaultAddr            = constants.DefaultPort
	_defaultReadTimeout     = constants.DefaultReadTimeout
	_defaultWriteTimeout    = constants.DefaultWriteTimeout
)

type RestServer struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	waitGroup       *errgroup.Group
	ctx             context.Context
	l               logger.Interface
}

func NewRestServer(handler http.Handler, l logger.Interface, waitGroup *errgroup.Group, ctx context.Context, opts ...Option) {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &RestServer{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		waitGroup:       waitGroup,
		ctx:             ctx,
		l:               l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.Start()

	s.Shutdown()
}

func (s *RestServer) Start() {

	s.waitGroup.Go(func() error {
		s.l.Info("start HTTP server at",
			zap.String("Adr: ", s.server.Addr),
			zap.String("ReadTimeout: ", s.server.ReadTimeout.String()),
			zap.String("WriteTimeout: ", s.server.WriteTimeout.String()),
			zap.String("ShutdownTimeout: ", s.shutdownTimeout.String()))

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

func (s *RestServer) Notify() <-chan error {
	return s.notify
}

func (s *RestServer) Shutdown() {
	s.waitGroup.Go(func() error {
		s.l.Info("graceful shutdown HTTP gateway server")
		<-s.ctx.Done()
		err := s.server.Shutdown(context.Background())
		if err != nil {
			s.l.Error("faild to shutdown HTTP server", zap.String("Error", err.Error()))
		}

		s.l.Info("HTTP gateway server is stopped")
		return nil
	})
}
