package consumers

import (
	"context"
	"encoding/json"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/order_service/internal/dataaccess/kafka/consumer"
	"go.uber.org/zap"
)

type OrderServiceKafkaConsumer interface {
	Start(ctx context.Context) error
}

type orderServiceKafkaConsumer struct {
	kafkaConsumer               consumer.Consumer
	l                           logger.Interface
	paymentTransactionCompleted PaymentTransactionCompletedMessageHandler
}

func NewOrderServiceKafkaConsumer(kafkaConsumer consumer.Consumer, l logger.Interface, paymentTransactionCompleted PaymentTransactionCompletedMessageHandler) OrderServiceKafkaConsumer {
	return &orderServiceKafkaConsumer{
		kafkaConsumer:               kafkaConsumer,
		l:                           l,
		paymentTransactionCompleted: paymentTransactionCompleted,
	}
}

func (o orderServiceKafkaConsumer) Start(ctx context.Context) error {
	o.kafkaConsumer.RegisterHandler(
		TOPIC_NAME_PAYMENT_SERVICE_TRANSACTION_COMPLETED,
		func(ctx context.Context, topic string, message []byte) error {
			var payload PaymentTransactionCompleted
			if err := json.Unmarshal(message, &payload); err != nil {
				o.l.Error("failed to unmarshal message", zap.Error(err))
				return err
			}

			o.paymentTransactionCompleted.Handle(ctx, payload)
			return nil
		},
	)

	return o.kafkaConsumer.Start(ctx)
}
