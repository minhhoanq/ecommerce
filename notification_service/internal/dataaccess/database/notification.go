package database

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
)

type NotificationDataAccessor interface {
	CreateNotification(ctx context.Context, arg *SendNotificationRequest) (*Notification, error)
}

type notificationDataAccessor struct {
	database Database
	l        logger.Interface
}

func NewNotificationDataAccessor(database Database, l logger.Interface) NotificationDataAccessor {
	return &notificationDataAccessor{
		database: database,
		l:        l,
	}
}

func (n *notificationDataAccessor) CreateNotification(ctx context.Context, arg *SendNotificationRequest) (*Notification, error) {
	return nil, nil
}
