package initial

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/config"
	"github.com/minhhoanq/ecommerce/notification_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/notification_service/internal/handler/grpc"
	"github.com/minhhoanq/ecommerce/notification_service/internal/service"
)

func Initial(cfg config.Config, l logger.Interface) (grpc.Server, error) {
	db, err := database.New(cfg, l)
	if err != nil {
		return nil, err
	}

	notificationDataAccessor := database.NewNotificationDataAccessor(db, l)

	notificationService := service.NewNotificationService(l, notificationDataAccessor)

	handler, err := grpc.NewHandler(l, notificationService)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewGRPCServer(cfg, handler, l)

	return grpcServer, nil
}
