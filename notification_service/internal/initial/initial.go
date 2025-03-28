package initial

import (
	"context"
	"fmt"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/config"
	"github.com/minhhoanq/ecommerce/notification_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/notification_service/internal/dataaccess/kafka/consumer"
	"github.com/minhhoanq/ecommerce/notification_service/internal/email"
	"github.com/minhhoanq/ecommerce/notification_service/internal/handler/consumers"
	"github.com/minhhoanq/ecommerce/notification_service/internal/handler/grpc"
	userservice "github.com/minhhoanq/ecommerce/notification_service/internal/handler/grpc/clients/user_service"
	"github.com/minhhoanq/ecommerce/notification_service/internal/service"
	"go.uber.org/zap"
)

func Initial(cfg config.Config, l logger.Interface) (grpc.Server, error) {
	db, err := database.New(cfg, l)
	if err != nil {
		return nil, err
	}

	kafkaConsumer, err := consumer.NewConsumer(cfg, l)
	if err != nil {
		return nil, err
	}

	notificationDataAccessor := database.NewNotificationDataAccessor(db, l)
	// user service grpc client
	userServiceClient, err := userservice.NewClient(cfg, l)
	// EmailSender
	emailSender := email.NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)

	notificationService := service.NewNotificationService(l, notificationDataAccessor, emailSender, userServiceClient)

	handler, err := grpc.NewHandler(l, notificationService)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get instance user service grpc client, err: ", err)
	}

	grpcServer := grpc.NewGRPCServer(cfg, handler, l)

	// start consumer threads
	emailNotifyCompleted := consumers.NewEmailNotifyCompletedMessageHandler(notificationService, l)
	newNotificationServiceKafkaConsumer := consumers.NewNotificaitonServiceKafkaConsumer(kafkaConsumer, l, emailNotifyCompleted)

	go func(ctx context.Context) {
		err := newNotificationServiceKafkaConsumer.Start(ctx)
		if err != nil {
			l.Error("failed to processing kafka consumer", zap.String("Error: ", err.Error()))
			return
		}
	}(context.Background())

	return grpcServer, nil
}
