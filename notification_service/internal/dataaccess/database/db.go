package database

import (
	"context"
	"fmt"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/config"
	"github.com/minhhoanq/ecommerce/notification_service/pkg/constants"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

var collection *mongo.Collection

type Database struct {
	*mongo.Client
	cfg config.Config
}

func New(cfg config.Config, l logger.Interface) (Database, error) {
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	otps := options.Client().ApplyURI(connString)

	client, err := mongo.Connect(otps)
	if err != nil {
		return Database{}, err
	}

	// check connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		l.Fatal("cannot connection to mongodb: ", zap.Error(err))
		return Database{}, err
	}

	l.Info("connection to mongodb successfully", zap.Int("port", cfg.DBPort), zap.String("name", cfg.DBName))

	// Migrations
	migrator := NewMigrator(client, l)
	// migration collections if not exists
	migrator.EnsureCollectionsAndIndexes(cfg.DBName, []string{constants.NOTIFICATION_COLLECTION})

	return Database{
		Client: client,
		cfg:    cfg,
	}, nil
}

func (db Database) returnCollectionPointer(collection string) *mongo.Collection {
	return db.Client.Database(db.cfg.DBName).Collection(collection)
}

func (db Database) Disconnect(ctx context.Context) error {
	return db.Disconnect(ctx)
}
