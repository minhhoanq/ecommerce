package migrations

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/pkg/constants"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

// Migration generate collection "notifications"
func CreateNotificationCollection(db *mongo.Database, l logger.Interface) error {
	collections, err := db.ListCollectionNames(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, name := range collections {
		if name == constants.NOTIFICATION_COLLECTION {
			return nil
		}
	}

	l.Info("creating collection", zap.String("name: ", constants.NOTIFICATION_COLLECTION))
	return db.CreateCollection(context.Background(), constants.NOTIFICATION_COLLECTION)
}
