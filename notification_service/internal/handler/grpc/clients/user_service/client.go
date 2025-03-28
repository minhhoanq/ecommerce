package userservice

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/config"
	"github.com/minhhoanq/ecommerce/notification_service/internal/generated/user_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(cfg config.Config, l logger.Interface) (user_service.UserServiceClient, error) {
	var otps = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(cfg.GRPCUserAddress, otps...)
	if err != nil {
		l.Error("failed to connect user grpc client", zap.Error(err))
		return nil, err
	}

	client := user_service.NewUserServiceClient(conn)

	return client, nil
}
