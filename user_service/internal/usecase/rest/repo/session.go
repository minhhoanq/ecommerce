package repo

import (
	"context"
	"time"

	"github.com/minhhoanq/ecommerce/user_service/internal/entity"

	"github.com/google/uuid"
)

type CreateSessionParams struct {
	UserId       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (q *SQLStore) CreateSession(ctx context.Context, arg CreateSessionParams) (*entity.Session, error) {
	var session *entity.Session = &entity.Session{
		UserId:       arg.UserId,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		IsBlocked:    arg.IsBlocked,
		ExpiredAt:    arg.ExpiredAt,
	}

	result := q.db.WithContext(ctx).Create(&session)

	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}

func (q *SQLStore) GetSessionByUserId(ctx context.Context, user_id uuid.UUID) (*entity.Session, error) {
	var session *entity.Session
	result := q.db.WithContext(ctx).First(&session).Where("user_id = ?", user_id)
	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}
