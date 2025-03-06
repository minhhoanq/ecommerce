package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/user_service/internal/handler/producers"
	"github.com/minhhoanq/ecommerce/user_service/pkg/constants"
	"github.com/minhhoanq/ecommerce/user_service/pkg/kafka"
)

type UserServiceKafkaConsumer interface {
	Start(ctx context.Context) error
}

type userServiceKafkaConsumer struct {
	kafkaConsumer     kafka.Consumer
	userSignupHandler UserSignupHandler
}

func NewUserServiceKafkaConsumer(kafkaConsumer kafka.Consumer, userSignupHandler UserSignupHandler) UserServiceKafkaConsumer {
	return &userServiceKafkaConsumer{
		kafkaConsumer:     kafkaConsumer,
		userSignupHandler: userSignupHandler,
	}
}

func (u userServiceKafkaConsumer) Start(ctx context.Context) error {
	u.kafkaConsumer.RegisterHandler(constants.TopicVerifyEmailForSignup,
		func(ctx context.Context, topic string, message []byte) error {
			fmt.Println("start register")

			var payload string
			if err := json.Unmarshal(message, &payload); err != nil {
				fmt.Println("Error unmarshal", err)
				return err
			}
			fmt.Println("payload", payload)

			parsedUUID, err := uuid.Parse(payload)
			if err != nil {
				return err
			}

			return u.userSignupHandler.Handle(ctx, producers.UserSignupParams{
				UserId: parsedUUID,
			})
		},
	)

	return u.kafkaConsumer.Start(ctx)
}
