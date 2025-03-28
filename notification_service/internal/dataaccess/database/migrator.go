package database

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type Migrator interface {
	EnsureCollectionsAndIndexes(dbName string, collections []string)
}

type migrator struct {
	client *mongo.Client
	l      logger.Interface
}

func NewMigrator(client *mongo.Client, l logger.Interface) Migrator {
	return &migrator{
		client: client,
		l:      l,
	}
}

func (m *migrator) EnsureCollectionsAndIndexes(dbName string, collections []string) {
	db := m.client.Database(dbName)

	for _, collName := range collections {
		// coll := db.Collection(collName)
		if !m.collectionExists(db, collName) {
			db.CreateCollection(context.Background(), collName)
			m.l.Info("creating collection...", zap.String("name: ", collName))
		}
	}
}

func (m *migrator) collectionExists(db *mongo.Database, collName string) bool {
	filter := bson.M{"name": collName}
	cursor, err := db.ListCollections(context.Background(), filter)
	if err != nil {
		m.l.Fatal("failed to check exists of collection")
	}
	return cursor.Next(context.Background())
}
