package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	pbuser "github.com/minhhoanq/ecommerce/user_service/internal/controller/grpc/v1/user_service"
	"github.com/minhhoanq/ecommerce/user_service/internal/entity"
	"github.com/minhhoanq/ecommerce/user_service/internal/usecase/rest/repo"
	"github.com/minhhoanq/ecommerce/user_service/pkg/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *GrpcServer) Signup(ctx context.Context, arg *pbuser.SignupRequest) (*pbuser.SignupResponse, error) {
	hashedPassword, err := entity.HashPassword(arg.Password)
	if err != nil {
		return nil, err
	}

	args := repo.CreateUserTxParams{
		CreateUserParams: repo.CreateUserParams{
			Username: arg.Username,
			Email:    arg.Email,
			Password: hashedPassword,
			RoleId:   1,
		},
		AfterCreate: func(user *entity.User) error {
			message, err := json.Marshal(user.ID)
			if err != nil {
				return err
			}

			return server.kafkaProducer.Produce(ctx, constants.TopicVerifyEmailForSignup, message)
		},
	}

	txResult, err := server.q.CreateUserTx(ctx, args)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists %s", err)

			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err.Error())
	}

	rsp := &pbuser.SignupResponse{
		User: convertUser(txResult.User),
	}

	return rsp, nil
}

func (server *GrpcServer) Signin(ctx context.Context, arg *pbuser.SigninRequest) (*pbuser.SigninResponse, error) {
	user, err := server.q.GetUserByUsername(ctx, arg.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, err
	}

	// password of user is hashed password
	err = entity.ComparePassword(arg.Password, user.Password)
	if err != nil {
		return nil, err
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, user.RoleId, server.cfg.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("create access token error: %v", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, user.RoleId, server.cfg.RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("create refresh token error: %v", err)
	}

	session, err := server.q.CreateSession(ctx, repo.CreateSessionParams{
		UserId:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, fmt.Errorf("faild create session: %v", err)
	}

	return &pbuser.SigninResponse{
		SessionId:             int32(session.ID),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}, nil
}

func (server *GrpcServer) GetUser(ctx context.Context, arg *pbuser.GetUserRequest) (*pbuser.GetUserResponse, error) {
	userID, err := uuid.Parse(arg.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %v", err)
	}
	user, err := server.q.GetUserByID(ctx, userID)
	if err != nil {
		fmt.Println("error get user by id", err)
		return nil, err
	}

	return &pbuser.GetUserResponse{
		User: convertUser(user),
	}, nil
}
