package token

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(user_id uuid.UUID, role_id int, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
