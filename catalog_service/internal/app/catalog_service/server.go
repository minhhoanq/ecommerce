package catalog_service

import (
	"context"

	"github.com/minhhoanq/ecommerce/catalog_service/internal/handler/grpc"
	"github.com/minhhoanq/ecommerce/common/logger"
	"go.uber.org/zap"
)

type Server struct {
	grpcServer grpc.Server
	l          logger.Interface
}

func NewServer(grpcServer grpc.Server, l logger.Interface) *Server {
	return &Server{
		grpcServer: grpcServer,
		l:          l,
	}
}

func (s *Server) Start() {
	s.l.Info("starting gRPC server")
	err := s.grpcServer.Start(context.Background())
	if err != nil {
		s.l.Error("failed start gRPC server", zap.Error(err))
	}
}
