package service

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/internal/dataaccess/database"
	pb "github.com/minhhoanq/ecommerce/notification_service/internal/generated/notification_service"
)

type NotificationService interface {
	SendNotification(ctx context.Context, arg *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error)
}

type notificationService struct {
	l                        logger.Interface
	notificationDataAccessor database.NotificationDataAccessor
}

func NewNotificationService(
	l logger.Interface,
	notificationDataAccessor database.NotificationDataAccessor,
) NotificationService {
	return &notificationService{
		l:                        l,
		notificationDataAccessor: notificationDataAccessor,
	}
}

func (n *notificationService) SendNotification(ctx context.Context, arg *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	payload := &database.SendNotificationRequest{
		UserID:   arg.UserId,
		Type:     arg.Type.String(),
		Title:    arg.Title,
		Message:  arg.Message,
		Status:   arg.Status.String(),
		Metadata: arg.Metadata,
	}

	notification, err := n.notificationDataAccessor.CreateNotification(ctx, payload)
	if err != nil {
		return nil, err
	}

	response := &pb.SendNotificationResponse{
		NotificationId: notification.ID,
		Status:         notification.Status,
	}

	return response, nil
}
