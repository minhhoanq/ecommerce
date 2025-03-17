package main

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/config"
	rest "github.com/minhhoanq/ecommerce/user_service/internal/app/http"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// init logger
	logger.Setup(cfg.Environment, cfg.LogLevel)
	// start rest server
	rest.RunRestServer(cfg)
}
