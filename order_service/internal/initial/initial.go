package initial

import (
	"github.com/minhhoanq/lifeat/common/logger"
	"github.com/minhhoanq/lifeat/order_service/config"
	"github.com/minhhoanq/lifeat/order_service/internal/dataaccess/database"
	"github.com/minhhoanq/lifeat/order_service/internal/dataaccess/redis"
	"github.com/minhhoanq/lifeat/order_service/internal/handler/grpc"
	catalogservice "github.com/minhhoanq/lifeat/order_service/internal/handler/grpc/clients/catalog_service"
	"github.com/minhhoanq/lifeat/order_service/internal/service"
)

func InitialServer(cfg config.Config, l logger.Interface) (grpc.Server, error) {
	db, err := database.New(cfg, l)
	if err != nil {
		return nil, err
	}

	redisInit := redis.NewRedis(cfg, l)
	// defer redisInit.Close()

	redisClient := redisInit.Connect()

	catalogServiceClient, err := catalogservice.NewClient(cfg, l)

	orderDataAccessor := database.NewOrderDataAccessor(db, l, catalogServiceClient, redisInit, redisClient)
	// Initialize the service
	orderService := service.NewOrderService(l, orderDataAccessor, catalogServiceClient)
	// Initialize the handler
	handler, err := grpc.NewHandler(l, orderService)
	// Initialize the server
	grpcServer := grpc.NewGRPCServer(cfg, handler, l)

	return grpcServer, nil
}
