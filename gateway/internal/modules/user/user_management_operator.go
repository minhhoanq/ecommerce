package user

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/user_service"
	"go.uber.org/zap"
)

type UserManagementOperator interface {
	Signup(ctx context.Context, arg *user_service.SignupRequest) (*user_service.SignupResponse, error)
	Signin(ctx context.Context, arg *user_service.SigninRequest) (*user_service.SigninResponse, error)
	GetUser(ctx context.Context, arg *user_service.GetUserRequest) (*user_service.GetUserResponse, error)
}

type userManagementOperator struct {
	userServiceClient user_service.UserServiceClient
	l                 logger.Interface
}

func NewUserManagementOperator(
	userServiceClient user_service.UserServiceClient,
	l logger.Interface) UserManagementOperator {
	return &userManagementOperator{
		userServiceClient: userServiceClient,
		l:                 l,
	}
}

func (u *userManagementOperator) Signup(ctx context.Context, arg *user_service.SignupRequest) (*user_service.SignupResponse, error) {
	return nil, nil
}

func (u *userManagementOperator) Signin(ctx context.Context, arg *user_service.SigninRequest) (*user_service.SigninResponse, error) {
	return nil, nil
}

func (u *userManagementOperator) GetUser(ctx context.Context, arg *user_service.GetUserRequest) (*user_service.GetUserResponse, error) {
	user, err := u.userServiceClient.GetUser(ctx, arg)
	if err != nil {
		u.l.Error("failed to get user", zap.Error(err))
		return nil, err
	}

	return user, nil
}
