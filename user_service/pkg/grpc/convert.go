package grpc

import (
	pbuser "github.com/minhhoanq/ecommerce/user_service/internal/controller/grpc/v1/user_service"
	"github.com/minhhoanq/ecommerce/user_service/internal/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user *entity.User) *pbuser.User {
	return &pbuser.User{
		Id:                user.ID.String(),
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangeAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
