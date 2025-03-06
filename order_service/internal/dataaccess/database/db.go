package database

import (
	"context"
	"fmt"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/order_service/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
	l logger.Interface
}

func New(cfg config.Config, l logger.Interface) (Database, error) {
	migrator, err := NewMigrator(cfg)
	if err != nil {
		return Database{}, err
	}

	err = migrator.Up(context.Background())
	if err != nil {
		return Database{}, err
	}

	db, err := NewDatabase(cfg, l)
	if err != nil {
		return Database{}, err
	}

	return Database{
		DB: db,
		l:  l,
	}, nil
}

func NewDatabase(cfg config.Config, l logger.Interface) (*gorm.DB, error) {
	// Create data source name (DSN) string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// Open GORM database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	l.Info("database is running on",
		zap.String("Host: ", cfg.DBHost),
		zap.String("Name: ", cfg.DBName),
		zap.Int("Port: ", cfg.DBPort))

	return db, nil
}

func (p *Database) Close() {
	db, _ := p.DB.DB()
	db.Close()
}
