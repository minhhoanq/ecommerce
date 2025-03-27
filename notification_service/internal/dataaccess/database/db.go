package database

import (
	"context"
	"fmt"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type Database struct {
	*mongo.Client
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
	}

	l.Info("connection to mongodb successfully", zap.Int("port", cfg.DBPort), zap.String("name", cfg.DBName))

	return Database{
		Client: client,
	}, nil
}
