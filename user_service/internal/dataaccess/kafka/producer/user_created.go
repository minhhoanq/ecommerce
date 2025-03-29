package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
)

const (
	TOPIC_NAME_USER_SERVICE_USER_CREATED = "user_service_user_created"
)

type UserCreated struct {
	UserID        uuid.UUID
	Email         string
	Username      string
	SecretCode    string
	VerifyEmailID int
}

type UserCreatedProducer interface {
	Produce(ctx context.Context, message UserCreated) error
}

type userCreatedProducer struct {
	producer Producer
	l        logger.Interface
}

func NewUserCreatedProducer(
	producer Producer,
	l logger.Interface,
) UserCreatedProducer {
	return &userCreatedProducer{
		producer: producer,
		l:        l,
	}
}

func (u *userCreatedProducer) Produce(ctx context.Context, message UserCreated) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message", err)
	}

	if err = u.producer.Produce(ctx, TOPIC_NAME_USER_SERVICE_USER_CREATED, messageBytes); err != nil {
		return fmt.Errorf("failed to send message to consumer, ", err)
	}

	return nil
}
