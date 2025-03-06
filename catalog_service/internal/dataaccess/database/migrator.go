package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/minhhoanq/ecommerce/catalog_service/configs"
)

var (
	//go:embed migrations/postgres
	migrationDirectoryPostgres embed.FS
)

type Migrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}

func NewMigrator(cfg configs.Config) (Migrator, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	sourceInstance, err := iofs.New(migrationDirectoryPostgres, "migrations/postgres")
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceInstance,
		cfg.DBName,
		driver,
	)
	if err != nil {
		return nil, err
	}

	return &migrator{
		instance: m,
	}, nil
}

type migrator struct {
	instance *migrate.Migrate
}

func (m *migrator) Down(ctx context.Context) error {
	err := m.instance.Down()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func (m *migrator) Up(ctx context.Context) error {
	err := m.instance.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
