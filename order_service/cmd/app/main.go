package main

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/order_service/config"
	"github.com/minhhoanq/ecommerce/order_service/internal/app/order_service"
	"github.com/minhhoanq/ecommerce/order_service/internal/initial"
	"go.uber.org/zap"
)

func main() {
	//init config
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// Inittial logger

	logger.Setup(config.Environment, config.LogLevel)
	l := logger.NewWrapLogger(zap.DebugLevel, false)

	// Initial server
	grpcServer, err := initial.InitialServer(config, l)
	if err != nil {
		l.Error("failed to initial server", zap.Error(err))
	}

	// NewServer
	server := order_service.NewServer(grpcServer, l)

	server.Start()
}
