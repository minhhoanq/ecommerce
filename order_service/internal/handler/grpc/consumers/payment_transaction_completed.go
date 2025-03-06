package consumers

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/order_service/internal/service"
	"go.uber.org/zap"
)

const (
	TOPIC_NAME_PAYMENT_SERVICE_TRANSACTION_COMPLETED = "payment_service_payment_transaction_completed"
)

type PaymentTransactionCompleted struct {
	OrderID                  uuid.UUID
	PaymentTransactionStatus string
}

type PaymentTransactionCompletedMessageHandler interface {
	Handle(ctx context.Context, message PaymentTransactionCompleted) error
}

type paymentTransactionCompletedMessageHandler struct {
	orderService service.OrderService
	l            logger.Interface
}

func NewPaymentTransactionCompletedMessageHandler(orderService service.OrderService, l logger.Interface) PaymentTransactionCompletedMessageHandler {
	return &paymentTransactionCompletedMessageHandler{
		orderService: orderService,
		l:            l,
	}
}

func (p paymentTransactionCompletedMessageHandler) Handle(ctx context.Context, message PaymentTransactionCompleted) error {
	p.l.Info("payment transaction completed received")
	p.l.Info("message", zap.Any("message", message))
	return nil
}
