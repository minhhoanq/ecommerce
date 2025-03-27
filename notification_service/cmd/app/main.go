package main

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/config"
	"github.com/minhhoanq/ecommerce/notification_service/internal/app/notification_service"
	"github.com/minhhoanq/ecommerce/notification_service/internal/initial"
	"go.uber.org/zap"
)

func main() {
	// init config
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// init logger
	logger.Setup(config.Environment, config.LogLevel)
	l := logger.NewWrapLogger(zap.DebugLevel, false)

	grpcServer, err := initial.Initial(config, l)
	if err != nil {
		l.Error("failed to initial server", zap.Error(err))
		panic(err)
	}

	// NewServer
	server := notification_service.NewServer(grpcServer, l)
	server.Start()
}
