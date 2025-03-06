package userservice

import (
	"github.com/minhhoanq/ecommerce/catalog_service/configs"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/generated/user_service"
	"github.com/minhhoanq/ecommerce/common/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(cfg configs.Config, l logger.Interface) (user_service.UserServiceClient, error) {
	// var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(cfg.GRPCUserAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error("failed to create user service client", zap.Error(err))
		return nil, err
	}

	client := user_service.NewUserServiceClient(conn)
	l.Info("connect to user grpc client", zap.String("address", cfg.GRPCUserAddress))

	return client, nil
}
