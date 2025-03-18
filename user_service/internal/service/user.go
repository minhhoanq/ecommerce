package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/config"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/user_service/internal/token"
	"github.com/minhhoanq/ecommerce/user_service/internal/util"
	"github.com/minhhoanq/ecommerce/user_service/internal/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*database.User, error)
	CreateUser(context.Context, CreateUserParams) (*database.User, error)
	Login(context.Context, LoginUsecaseParams) (*LoginUscaseResponse, error)
	RenewAccessToken(ctx context.Context, arg RenewAccessTokenUsecaseParams) (*RenewAccessTokenUsecaseResponse, error)
}

type userService struct {
	tokenMaker       token.Maker
	cfg              config.Config
	l                logger.Interface
	userDataAccessor database.UserDataAccessor
	taskDistributor  worker.TaskDistributor
}

func NewUserService(tokenMaker token.Maker,
	cfg config.Config,
	l logger.Interface,
	userDataccessor database.UserDataAccessor,
	taskDistributor worker.TaskDistributor) UserService {
	return &userService{
		tokenMaker:       tokenMaker,
		cfg:              cfg,
		l:                l,
		userDataAccessor: userDataccessor,
		taskDistributor:  taskDistributor,
	}
}

func (us *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*database.User, error) {
	user, err := us.userDataAccessor.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	RoleId   int    `json:"role_id"`
}

func (us *userService) CreateUser(ctx context.Context, arg CreateUserParams) (*database.User, error) {
	hashedPassword, err := util.HashPassword(arg.Password)
	if err != nil {
		return nil, err
	}

	args := database.CreateUserTxParams{
		CreateUserParams: database.CreateUserParams{
			Username: arg.Username,
			Email:    arg.Email,
			Password: hashedPassword,
			RoleId:   arg.RoleId,
		},
		AfterCreate: func(user *database.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{UserId: user.ID}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.Queue(worker.QueueCritial),
			}

			return us.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := us.userDataAccessor.CreateUserTx(ctx, args)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists %s", err)

			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err.Error())
	}

	// producer, err := kafka.NewProducer(uc.cfg, uc.l)
	// producer.Produce(ctx, "VerifyEmailSignup", )

	return &txResult.User, nil
}

type LoginUscaseResponse struct {
	SessionID             int            `json:"session_id"`
	AccessToken           string         `json:"access_token"`
	RefreshToken          string         `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time      `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time      `json:"refresh_token_expires_at"`
	User                  *database.User `json:"user"`
}

type LoginUsecaseParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (us *userService) Login(ctx context.Context, arg LoginUsecaseParams) (*LoginUscaseResponse, error) {
	user, err := us.userDataAccessor.GetUserByUsername(ctx, arg.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, err
	}

	// password of user is hashed password
	err = util.ComparePassword(arg.Password, user.Password)
	if err != nil {
		return nil, err
	}

	accessToken, accessPayload, err := us.tokenMaker.CreateToken(user.ID, user.RoleId, us.cfg.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("create access token error: %v", err)
	}

	refreshToken, refreshPayload, err := us.tokenMaker.CreateToken(user.ID, user.RoleId, us.cfg.RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("create refresh token error: %v", err)
	}

	session, err := us.userDataAccessor.CreateSession(ctx, database.CreateSessionParams{
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

	return &LoginUscaseResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  user,
	}, nil
}

type RenewAccessTokenUsecaseParams struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenUsecaseResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (us *userService) RenewAccessToken(ctx context.Context, arg RenewAccessTokenUsecaseParams) (*RenewAccessTokenUsecaseResponse, error) {
	refreshTokenPayload, err := us.tokenMaker.VerifyToken(arg.RefreshToken)
	if err != nil {
		return nil, err
	}

	session, err := us.userDataAccessor.GetSessionByUserId(ctx, refreshTokenPayload.UserId)
	if err != nil {
		return nil, err
	}

	if session.IsBlocked {
		return nil, fmt.Errorf("token is blocked")
	}

	if session.UserId != refreshTokenPayload.UserId {
		return nil, fmt.Errorf("incorrect session user")
	}

	if session.RefreshToken != arg.RefreshToken {
		return nil, fmt.Errorf("mismatch session token")
	}

	if time.Now().After(session.ExpiredAt) {
		return nil, fmt.Errorf("expired session")
	}

	user, err := us.userDataAccessor.GetUserByID(ctx, session.UserId)
	if err != nil {
		return nil, err
	}

	accessToken, payload, err := us.tokenMaker.CreateToken(user.ID, user.RoleId, us.cfg.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	rsp := &RenewAccessTokenUsecaseResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: payload.ExpiredAt,
	}

	return rsp, nil
}
