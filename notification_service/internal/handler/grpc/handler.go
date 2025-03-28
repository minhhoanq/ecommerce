package grpc

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	pb "github.com/minhhoanq/ecommerce/notification_service/internal/generated/notification_service"
	"github.com/minhhoanq/ecommerce/notification_service/internal/service"
)

type Handler struct {
	pb.UnimplementedNotificationServiceServer
	l                   logger.Interface
	notificationService service.NotificationService
}

func NewHandler(l logger.Interface, notificationService service.NotificationService) (pb.NotificationServiceServer, error) {
	return &Handler{
		l:                   l,
		notificationService: notificationService,
	}, nil
}

func (h *Handler) CreateNotification(ctx context.Context, arg *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	return h.notificationService.SendNotification(ctx, arg)
}
