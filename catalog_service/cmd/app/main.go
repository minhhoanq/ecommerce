package main

import (
	"fmt"

	"github.com/minhhoanq/ecommerce/catalog_service/configs"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/app/catalog_service"
	"github.com/minhhoanq/ecommerce/catalog_service/internal/initial"
	"github.com/minhhoanq/ecommerce/common/logger"
	"go.uber.org/zap"
)

func main() {
	// Config
	config, err := configs.LoadConfig(".")
	if err != nil {
		fmt.Println("Error", err.Error())
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
	server := catalog_service.NewServer(grpcServer, l)

	server.Start()
}
