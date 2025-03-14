package initial

import (
	"context"

	"github.com/minhhoanq/ecommerce/catalog_service/configs"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/dataaccess/redis"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/dataaccess/s3"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/handler/grpc"
	userservice "github.com/minhhoanq/ecommerce/catalog_service/internal/handler/grpc/clients/user_service"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/service"
	"github.com/minhhoanq/ecommerce/common/logger"
)

func InitialServer(cfg configs.Config, l logger.Interface) (grpc.Server, error) {

	db, err := database.New(cfg, l)
	if err != nil {
		return nil, err
	}

	// s3 client
	s3Client, err := s3.NewClient(cfg, l)
	if err != nil {
		return nil, err
	}

	// connect to redis
	redis := redis.NewRedis(cfg, l)
	redisClient := redis.Connect()
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	userServiceClient, err := userservice.NewClient(cfg, l)
	if err != nil {
		return nil, err
	}
	catalogAccessor := database.NewCatalogDataAccessor(db, l)
	catalogService := service.NewCatalogService(db.DB, l, catalogAccessor, userServiceClient, s3Client)
	handler, err := grpc.NewHandler(catalogService, l, redisClient)
	if err != nil {
		return nil, err
	}
	server := grpc.NewGRPCServer(cfg, handler, l)

	return server, nil
}
