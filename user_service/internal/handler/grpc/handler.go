package grpc

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	pb "github.com/minhhoanq/ecommerce/user_service/internal/generated/user_service"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	l logger.Interface
}

func NewHander(l logger.Interface) (pb.UserServiceServer, error) {
	return &Handler{
		l: l,
	}, nil
}
