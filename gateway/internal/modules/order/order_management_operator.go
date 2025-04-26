package order

import (
	"context"
	"fmt"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/order_service"
)

type OrderManagementOperator interface {
	CreateOrder(ctx context.Context, arg *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error)
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

func (o *orderManagementOperator) CreateOrder(ctx context.Context, arg *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error) {
	order, err := o.orderServiceClient.CreateOrder(ctx, arg)
	fmt.Println(arg)
	if err != nil {
		return nil, err
	}

	return order, nil
}
