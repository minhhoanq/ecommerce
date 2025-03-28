package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/notification_service/internal/dataaccess/kafka/consumer"
	"go.uber.org/zap"
)

type NotificationServiceKafkaConsumer interface {
	Start(ctx context.Context) error
}

type notificationServiceKafkaConsumer struct {
	kafkaConsumer        consumer.Consumer
	emailNotifyCompleted EmailNotifyCompletedMessageHandler
	l                    logger.Interface
}

func NewNotificaitonServiceKafkaConsumer(
	kafkaConsumer consumer.Consumer,
	l logger.Interface,
	emailNotifyCompleted EmailNotifyCompletedMessageHandler,
) NotificationServiceKafkaConsumer {
	return &notificationServiceKafkaConsumer{
		kafkaConsumer: kafkaConsumer,
		l:             l,
	}
}

func (n notificationServiceKafkaConsumer) Start(ctx context.Context) error {
	n.kafkaConsumer.RegisterHandler(
		TOPIC_NAME_NOTIFICATION_SERVICE_EMAIL_NOTIFY_COMPLETED,
		func(ctx context.Context, topic string, message []byte) error {
			var payload EmailNotifyCompleted
			if err := json.Unmarshal(message, &payload); err != nil {
				n.l.Error("failed to unmarshal message", zap.Error(err))
				return fmt.Errorf("failed to unmarshal message", err)
			}

			n.emailNotifyCompleted.Handle(ctx, payload)
			return nil
		},
	)

	return n.Start(ctx)
}
