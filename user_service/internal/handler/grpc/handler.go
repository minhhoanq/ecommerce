package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
	pb "github.com/minhhoanq/ecommerce/user_service/internal/generated/user_service"
	"github.com/minhhoanq/ecommerce/user_service/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (h *Handler) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	arg := service.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		RoleId:   1,
	}

	user, err := h.userService.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	response := &pb.SignupResponse{
		User: &pb.User{
			Id:                user.ID.String(),
			Username:          user.Username,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangeAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}

	return response, nil
}

func (h *Handler) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	arg := service.LoginParams{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := h.userService.Login(ctx, arg)
	if err != nil {
		return nil, err
	}

	response := &pb.SigninResponse{
		AccessToken:           result.AccessToken,
		AccessTokenExpiresAt:  timestamppb.New(result.AccessTokenExpiresAt),
		RefreshToken:          result.RefreshToken,
		RefreshTokenExpiresAt: timestamppb.New(result.RefreshTokenExpiresAt),
		User: &pb.User{
			Id:                result.User.ID.String(),
			Email:             result.User.Email,
			Username:          result.User.Username,
			PasswordChangedAt: timestamppb.New(result.User.PasswordChangeAt),
			CreatedAt:         timestamppb.New(result.User.CreatedAt),
		},
	}

	return response, nil
}

func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userId, err := uuid.Parse(req.Id)
	if err != nil {
		h.l.Error("failed to parse user id", zap.Error(err))
		return nil, err
	}

	user, err := h.userService.GetUserByID(ctx, userId)
	if err != nil {
		h.l.Error("failed to get user", zap.Error(err))
		return nil, err
	}

	response := &pb.GetUserResponse{
		User: &pb.User{
			Id:                user.ID.String(),
			Username:          user.Username,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangeAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}

	return response, nil
}
