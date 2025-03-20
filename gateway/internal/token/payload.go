package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is valid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	UserId    uuid.UUID `json:"user_id"`
	RoleId    int       `json:"role_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with specific UserId, role_id, duration
func NewPayload(user_id uuid.UUID, role_id int, duration time.Duration) *Payload {
	payload := &Payload{
		UserId:    user_id,
		RoleId:    role_id,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
