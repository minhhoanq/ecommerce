package order

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/order_service"
)

type OrderManagementOperator interface {
	CreateOrder(ctx context.Context, arg any) (order_service.CreateOrderResponse, error)
}

type orderManagementOperator struct {
	orderServiceClient order_service.OrderServiceClient
	l                  logger.Interface
}

func NewOrderManagementOperator(orderServiceClient order_service.OrderServiceClient, l logger.Interface) OrderManagementOperator {
	return &orderManagementOperator{
		orderServiceClient: orderServiceClient,
		l:                  l,
	}
}

func (o *orderManagementOperator) CreateOrder(ctx context.Context, arg any) (order_service.CreateOrderResponse, error) {
	return order_service.CreateOrderResponse{}, nil
}
