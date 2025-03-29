package database

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	RoleId   int    `json:"role_id"`
}

type UserCreatedParams struct {
	UserID        uuid.UUID
	Email         string
	Username      string
	SecretCode    string
	VerifyEmailID int
}

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(params *UserCreatedParams) error
}

type CreateUserTxResult struct {
	User User
}

type CreateVerifyEmailParams struct {
	UserId     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
}

type CreateSessionParams struct {
	UserId       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
}
