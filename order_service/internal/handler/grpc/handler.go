package grpc

import (
	"context"

	"github.com/minhhoanq/ecommerce/common/logger"
	pb "github.com/minhhoanq/ecommerce/order_service/internal/generated/order_service"
	"github.com/minhhoanq/ecommerce/order_service/internal/service"
)

type Handler struct {
	pb.UnimplementedOrderServiceServer
	l            logger.Interface
	orderService service.OrderService
}

func NewHandler(l logger.Interface, orderService service.OrderService) (pb.OrderServiceServer, error) {
	return &Handler{
		l:            l,
		orderService: orderService,
	}, nil
}

func (h *Handler) CreateOrder(ctx context.Context, arg *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return h.orderService.CreateOrder(ctx, arg)
}
