package producers

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/pkg/constants"
	"github.com/minhhoanq/ecommerce/user_service/pkg/kafka"
	"go.uber.org/zap"
)

type UserSignupParams struct {
	UserId uuid.UUID `json:"user_id"`
}

type UserSignupProducer interface {
	Produce(ctx context.Context, message UserSignupParams) error
}

type userSignupProducer struct {
	producer kafka.Producer
	l        logger.Interface
}

func NewUserSignupProducer(producer kafka.Producer, l logger.Interface) UserSignupProducer {
	return &userSignupProducer{
		producer: producer,
		l:        l,
	}
}

func (u userSignupProducer) Produce(ctx context.Context, message UserSignupParams) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		u.l.Error("failed to marshal message event", zap.String("Error", err.Error()))
		return err
	}

	err = u.producer.Produce(ctx, constants.TopicVerifyEmailForSignup, messageBytes)
	if err != nil {
		u.l.Error("failed to produce user signup message event", zap.String("Error", err.Error()))
		return err
	}

	return nil
}
