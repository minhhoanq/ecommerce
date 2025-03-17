package initial

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/config"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/user_service/internal/handler/grpc"
)

func InitialServer(cfg config.Config, l logger.Interface) (grpc.Server, error) {
	db, err := database.New(cfg, l)
	if err != nil {
		return nil, err
	}

	// initialize handler
	handler, err := grpc.NewHander(l)
	if err != nil {
		return nil, err
	}
	// initialize the server
	grpcServer := grpc.NewServer(cfg, l, handler)

	return grpcServer, nil
}
