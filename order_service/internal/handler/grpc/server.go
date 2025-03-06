package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/order_service/config"
	pb "github.com/minhhoanq/ecommerce/order_service/internal/generated/order_service"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	cfg     config.Config
	handler pb.OrderServiceServer
	l       logger.Interface
}

func NewGRPCServer(
	cfg config.Config,
	handler pb.OrderServiceServer,
	l logger.Interface,
) Server {
	return &server{
		cfg:     cfg,
		handler: handler,
		l:       l,
	}
}

func (s *server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.GRPCServerAddress)
	if err != nil {
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
	pb.RegisterOrderServiceServer(grpcServer, s.handler)
	fmt.Println("registered services:", grpcServer.GetServiceInfo())
	s.l.Info("gRPC server is running on", zap.String("Address: ", s.cfg.GRPCServerAddress))
	return grpcServer.Serve(listener)
}
