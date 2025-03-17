package grpc

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/config"
	pb "github.com/minhhoanq/ecommerce/user_service/internal/generated/user_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	cfg     config.Config
	l       logger.Interface
	handler pb.UserServiceServer
}

func NewServer(cfg config.Config, l logger.Interface, handler pb.UserServiceServer) Server {
	return &server{
		cfg:     cfg,
		l:       l,
		handler: handler,
	}
}

func (s *server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.GRPCServerAddress)
	if err != nil {
		s.l.Error("failed to connect gRPC server", zap.Error(err))
		return err
	}
	defer listener.Close()

	var opts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	s.l.Info("gRPC server is running on", zap.String("Address: ", s.cfg.GRPCServerAddress))

	return grpcServer.Serve(listener)
}
