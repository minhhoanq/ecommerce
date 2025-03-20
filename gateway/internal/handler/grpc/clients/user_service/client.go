package userservice

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/config"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/user_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(cfg config.Config, l logger.Interface) (user_service.UserServiceClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(cfg.GRPCUserAddress, opts...)
	if err != nil {
		return nil, err
	}

	client := user_service.NewUserServiceClient(conn)
	l.Info("connection to user service grpc client", zap.String("Address: ", cfg.GRPCUserAddress))

	return client, nil
}
