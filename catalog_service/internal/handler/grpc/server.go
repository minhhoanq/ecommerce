package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/minhhoanq/ecommerce/catalog_service/configs"
	pb "github.com/minhhoanq/ecommerce/catalog_service/internal/generated/catalog_service"
	"github.com/minhhoanq/ecommerce/common/logger"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	cfg     configs.Config
	handler pb.CatalogServiceServer
	l       logger.Interface
}

func NewGRPCServer(
	cfg configs.Config,
	handler pb.CatalogServiceServer,
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
	pb.RegisterCatalogServiceServer(grpcServer, s.handler)
	fmt.Println("Registered services:", grpcServer.GetServiceInfo())
	s.l.Info("gRPC server is running on", zap.String("Address: ", s.cfg.GRPCServerAddress))
	return grpcServer.Serve(listener)
}
