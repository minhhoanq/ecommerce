package main

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/config"
	"github.com/minhhoanq/ecommerce/gateway/internal/initial"
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
	initial.Initial(config, l)
	if err != nil {
		l.Error("failed to initial server", zap.Error(err))
	}
}
