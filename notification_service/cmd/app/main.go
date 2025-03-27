package main

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/config"
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

	l.Info("")
}
