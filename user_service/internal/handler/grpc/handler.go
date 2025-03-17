package grpc

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	pb "github.com/minhhoanq/ecommerce/user_service/internal/generated/user_service"
	"github.com/minhhoanq/ecommerce/user_service/internal/service"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	l           logger.Interface
	userService service.UserService
}

func NewHander(l logger.Interface, userService service.UserService) (pb.UserServiceServer, error) {
	return &Handler{
		l:           l,
		userService: userService,
	}, nil
}
