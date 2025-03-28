package database

import (
	"context"
	"fmt"
	"time"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/pkg/constants"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationDataAccessor interface {
	CreateNotification(ctx context.Context, arg *SendNotificationRequest) (*SendNotificationResponse, error)
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

func (n *notificationDataAccessor) CreateNotification(ctx context.Context, arg *SendNotificationRequest) (*SendNotificationResponse, error) {
	collection := n.database.returnCollectionPointer(constants.NOTIFICATION_COLLECTION)

	result, err := collection.InsertOne(ctx,
		Notification{
			UserID:    arg.UserID,
			Type:      arg.Type,
			Title:     arg.Title,
			Message:   arg.Message,
			Status:    arg.Status,
			Metadata:  arg.Metadata,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert notify to database, err: ", err)
	}

	// convert insertedID type
	insertedID := result.InsertedID.(bson.ObjectID).Hex()

	response := &SendNotificationResponse{
		ID:     insertedID,
		Status: arg.Status,
	}

	return response, nil
}
