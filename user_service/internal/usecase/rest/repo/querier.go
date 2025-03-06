package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/user_service/internal/entity"
)

type Querier interface {
	GetUserByID(ctx context.Context, user_id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (*entity.VerifyEmail, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (*entity.Session, error)
	GetSessionByUserId(ctx context.Context, user_id uuid.UUID) (*entity.Session, error)
	//Tx
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}
