package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/order_service/internal/dataaccess/database"
	"go.uber.org/zap"
)

const (
	TOPIC_NAME_ORDER_SERVICE_ORDER_CREATED = "order_service_order_created"
)

type OrderCreated struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Items  []database.OrderItem
}

type OrderCreatedProducer interface {
	Produce(ctx context.Context, message OrderCreated) error
}

type orderCreatedProducer struct {
	producer Producer
	l        logger.Interface
}

func NewOrderCreatedProducer(producer Producer, l logger.Interface) OrderCreatedProducer {
	return &orderCreatedProducer{
		producer: producer,
		l:        l,
	}
}

func (o orderCreatedProducer) Produce(ctx context.Context, message OrderCreated) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		o.l.Error("failed to marshal message", zap.Error(err))
		return err
	}

	if err = o.producer.Produce(ctx, TOPIC_NAME_ORDER_SERVICE_ORDER_CREATED, messageBytes); err != nil {
		o.l.Error("failed to produce order created event", zap.Error(err))
		return err
	}

	return nil
}
