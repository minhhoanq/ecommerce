package user

import (
	"context"
	"fmt"

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
	user, err := u.userServiceClient.Signup(ctx, arg)
	if err != nil {
		u.l.Error("failed to signup", zap.Error(err))
		return nil, fmt.Errorf("failed to signup ", err)
	}

	response := &user_service.SignupResponse{
		User: &user_service.User{
			Id:                user.User.Id,
			Username:          user.User.Username,
			Email:             user.User.Email,
			PasswordChangedAt: user.User.PasswordChangedAt,
			CreatedAt:         user.User.CreatedAt,
		},
	}

	return response, nil
}

func (u *userManagementOperator) Signin(ctx context.Context, arg *user_service.SigninRequest) (*user_service.SigninResponse, error) {
	user, err := u.userServiceClient.Signin(ctx, arg)
	if err != nil {
		u.l.Error("failed to signin", zap.Error(err))
		return nil, fmt.Errorf("failed to signin ", err)
	}

	response := &user_service.SigninResponse{
		User: &user_service.User{
			Id:                user.User.Id,
			Username:          user.User.Username,
			Email:             user.User.Email,
			PasswordChangedAt: user.User.PasswordChangedAt,
			CreatedAt:         user.User.CreatedAt,
		},
		AccessToken:           user.AccessToken,
		AccessTokenExpiresAt:  user.AccessTokenExpiresAt,
		RefreshToken:          user.RefreshToken,
		RefreshTokenExpiresAt: user.RefreshTokenExpiresAt,
	}

	return response, nil
}

func (u *userManagementOperator) GetUser(ctx context.Context, arg *user_service.GetUserRequest) (*user_service.GetUserResponse, error) {
	user, err := u.userServiceClient.GetUser(ctx, arg)
	fmt.Println("id", arg.Id)
	if err != nil {
		u.l.Error("failed to get user", zap.Error(err))
		return nil, err
	}

	return user, nil
}
