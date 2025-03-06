package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/config"
	"github.com/minhhoanq/ecommerce/user_service/internal/entity"
	"github.com/minhhoanq/ecommerce/user_service/internal/token"
	"github.com/minhhoanq/ecommerce/user_service/internal/usecase/rest/repo"
	"github.com/minhhoanq/ecommerce/user_service/internal/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type UserUsecase interface {
	GetUserByID(context.Context, uuid.UUID) (*entity.User, error)
	CreateUser(context.Context, CreateUserUsecaseParams) (*entity.User, error)
	Login(context.Context, LoginUsecaseParams) (*LoginUscaseResponse, error)
	RenewAccessToken(ctx context.Context, arg RenewAccessTokenUsecaseParams) (*renewAccessTokenUsecaseResponse, error)
}

type userUsecase struct {
	tokenMaker      token.Maker
	cfg             config.Config
	taskDistributor worker.TaskDistributor
	q               repo.Querier
	l               logger.Interface
}

func New(q repo.Querier, tokenMaker token.Maker, cfg config.Config, taskDistributor worker.TaskDistributor, l logger.Interface) UserUsecase {
	return &userUsecase{
		tokenMaker:      tokenMaker,
		cfg:             cfg,
		taskDistributor: taskDistributor,
		q:               q,
		l:               l,
	}
}

func (uc userUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type CreateUserUsecaseParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	RoleId   int    `json:"role_id"`
}

func (uc userUsecase) CreateUser(ctx context.Context, arg CreateUserUsecaseParams) (*entity.User, error) {
	hashedPassword, err := entity.HashPassword(arg.Password)
	if err != nil {
		return nil, err
	}

	args := repo.CreateUserTxParams{
		CreateUserParams: repo.CreateUserParams{
			Username: arg.Username,
			Email:    arg.Email,
			Password: hashedPassword,
			RoleId:   arg.RoleId,
		},
		AfterCreate: func(user *entity.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{UserId: user.ID}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.Queue(worker.QueueCritial),
			}

			return uc.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := uc.q.CreateUserTx(ctx, args)

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

	return txResult.User, nil
}

type LoginUsecaseParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUscaseResponse struct {
	SessionID             int          `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  *entity.User `json:"user"`
}

func (uc userUsecase) Login(ctx context.Context, arg LoginUsecaseParams) (*LoginUscaseResponse, error) {
	user, err := uc.q.GetUserByUsername(ctx, arg.Username)
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

	accessToken, accessPayload, err := uc.tokenMaker.CreateToken(user.ID, user.RoleId, uc.cfg.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("create access token error: %v", err)
	}

	refreshToken, refreshPayload, err := uc.tokenMaker.CreateToken(user.ID, user.RoleId, uc.cfg.RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("create refresh token error: %v", err)
	}

	session, err := uc.q.CreateSession(ctx, repo.CreateSessionParams{
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

type renewAccessTokenUsecaseResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (uc userUsecase) RenewAccessToken(ctx context.Context, arg RenewAccessTokenUsecaseParams) (*renewAccessTokenUsecaseResponse, error) {
	refreshTokenPayload, err := uc.tokenMaker.VerifyToken(arg.RefreshToken)
	if err != nil {
		return nil, err
	}

	session, err := uc.q.GetSessionByUserId(ctx, refreshTokenPayload.UserId)
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

	user, err := uc.q.GetUserByID(ctx, session.UserId)
	if err != nil {
		return nil, err
	}

	accessToken, payload, err := uc.tokenMaker.CreateToken(user.ID, user.RoleId, uc.cfg.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	rsp := &renewAccessTokenUsecaseResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: payload.ExpiredAt,
	}

	return rsp, nil
}
