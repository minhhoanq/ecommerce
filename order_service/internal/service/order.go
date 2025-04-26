package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/order_service/internal/dataaccess/database"
	"github.com/minhhoanq/ecommerce/order_service/internal/dataaccess/kafka/producer"
	"github.com/minhhoanq/ecommerce/order_service/internal/generated/catalog_service"
	pb "github.com/minhhoanq/ecommerce/order_service/internal/generated/order_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService interface {
	CreateOrder(ctx context.Context, arg *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
}

type orderService struct {
	l                    logger.Interface
	orderDataAccessor    database.OrderDataAccessor
	catalogServiceClient catalog_service.CatalogServiceClient
	kafkaProducer        producer.Producer
}

// catalogAccessor database.CatalogDataAccessor
func NewOrderService(l logger.Interface,
	orderDataAccessor database.OrderDataAccessor,
	catalogServiceClient catalog_service.CatalogServiceClient,
	kafkaProducer producer.Producer,
) OrderService {
	return &orderService{
		l:                    l,
		orderDataAccessor:    orderDataAccessor,
		catalogServiceClient: catalogServiceClient,
		kafkaProducer:        kafkaProducer,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, arg *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := &database.CreateOrderRequest{
		UserID:     arg.UserId,
		Address:    "",
		Status:     "PENDING",
		OrderItems: make([]database.CreateOrderItemRequest, 0, len(arg.CartItems)),
	}

	for _, item := range arg.CartItems {
		fmt.Println("item: ", item.SkuId)
		sku, err := o.catalogServiceClient.GetSKU(ctx, &catalog_service.GetSKURequest{
			SkuId: item.SkuId,
		})
		fmt.Println(sku)
		if err != nil {
			return nil, fmt.Errorf("failed to find sku", err.Error())
		}

		if sku.Sku.Id == "" || sku.Sku == nil {
			return nil, fmt.Errorf("invalid product", err)
		}

		inventory, err := o.catalogServiceClient.GetInventorySKU(ctx, &catalog_service.GetInventorySKURequest{
			SkuId: sku.Sku.Id,
		})

		if inventory.Inventory.Stock-item.Quantity < 0 {
			return nil, fmt.Errorf("quantity > stock")
		}

		order.OrderItems = append(order.OrderItems, database.CreateOrderItemRequest{
			SkuID:    item.SkuId,
			Quantity: item.Quantity,
		})
	}

	result, err := o.orderDataAccessor.CreateOrder(ctx, order)

	if err != nil {
		return nil, err
	}

	fmt.Println("rs ", result)

	message := &producer.OrderCreated{
		ID:     result.Order.ID,
		UserID: result.Order.UserID,
		Items:  make([]database.OrderItem, 0),
	}

	for _, item := range result.OrderItems {
		message.Items = append(message.Items, item)
	}

	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json message")
	}

	// send event order created to payment service
	o.kafkaProducer.Produce(ctx, producer.TOPIC_NAME_ORDER_SERVICE_ORDER_CREATED, jsonMessage)

	response := &pb.CreateOrderResponse{
		Order: &pb.Order{
			Id:            result.Order.ID.String(),
			UserId:        result.Order.UserID.String(),
			Status:        result.Order.Status,
			PaymentMethod: result.Order.PaymentMethod,
			CreatedAt:     timestamppb.New(result.Order.CreatedAt),
			UpdatedAt:     timestamppb.New(result.Order.UpdatedAt),
			Items:         make([]*pb.OrderItem, 0, len(arg.CartItems)),
		},
	}

	for _, item := range result.OrderItems {
		response.Order.Items = append(response.Order.Items, &pb.OrderItem{
			Id:       item.ID.String(),
			OrderId:  item.OrderID.String(),
			SkuId:    item.SkuID.String(),
			Quantity: item.Quantity,
		})
	}

	return response, nil
}
