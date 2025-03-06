package initial

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/order_service/config"
	"github.com/minhhoanq/ecommerce/order_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/order_service/internal/dataaccess/kafka/consumer"
	"github.com/minhhoanq/ecommerce/order_service/internal/dataaccess/kafka/producer"
	"github.com/minhhoanq/ecommerce/order_service/internal/dataaccess/redis"
	"github.com/minhhoanq/ecommerce/order_service/internal/handler/consumers"
	"github.com/minhhoanq/ecommerce/order_service/internal/handler/grpc"
	catalogservice "github.com/minhhoanq/ecommerce/order_service/internal/handler/grpc/clients/catalog_service"
	"github.com/minhhoanq/ecommerce/order_service/internal/service"
	"go.uber.org/zap"
)

func InitialServer(cfg config.Config, l logger.Interface) (grpc.Server, error) {
	db, err := database.New(cfg, l)
	if err != nil {
		return nil, err
	}

	redisInit := redis.NewRedis(cfg, l)
	// defer redisInit.Close()

	redisClient := redisInit.Connect()

	kafkaProducer, err := producer.NewProducer(cfg, l)
	if err != nil {
		return nil, err
	}

	kafkaConsumer, err := consumer.NewConsumer(cfg, l)
	if err != nil {
		return nil, err
	}

	catalogServiceClient, err := catalogservice.NewClient(cfg, l)

	orderDataAccessor := database.NewOrderDataAccessor(db, l, catalogServiceClient, redisInit, redisClient)
	// Initialize the service
	orderService := service.NewOrderService(l, orderDataAccessor, catalogServiceClient, kafkaProducer)
	// Initialize the handler
	handler, err := grpc.NewHandler(l, orderService)
	// Initialize the server
	grpcServer := grpc.NewGRPCServer(cfg, handler, l)

	// start consumer threads
	paymentTransactionCompleted := consumers.NewPaymentTransactionCompletedMessageHandler(orderService, l)
	newOrderServiceKafkaConsumer := consumers.NewOrderServiceKafkaConsumer(kafkaConsumer, l, paymentTransactionCompleted)

	go func(ctx context.Context) {
		err := newOrderServiceKafkaConsumer.Start(ctx)
		if err != nil {
			l.Error("failed to processing kafka consumer", zap.String("Error: ", err.Error()))
			return
		}
	}(context.Background())

	return grpcServer, nil
}
