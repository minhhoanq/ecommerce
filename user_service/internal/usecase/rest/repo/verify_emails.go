package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/user_service/internal/entity"
)

type CreateVerifyEmailParams struct {
	UserId     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
}

func (q *SQLStore) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (*entity.VerifyEmail, error) {
	var verifyEmail *entity.VerifyEmail = &entity.VerifyEmail{
		UserId:     arg.UserId,
		Email:      arg.Email,
		SecretCode: arg.SecretCode,
	}
	result := q.db.WithContext(ctx).Create(&verifyEmail)
	if result.Error != nil {
		return nil, result.Error
	}

	return verifyEmail, nil
}
