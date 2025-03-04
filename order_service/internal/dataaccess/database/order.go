package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/minhhoanq/lifeat/common/logger"
	redisOrder "github.com/minhhoanq/lifeat/order_service/internal/dataaccess/redis"
	"github.com/minhhoanq/lifeat/order_service/internal/generated/catalog_service"
	"github.com/minhhoanq/lifeat/order_service/pkg/helper"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderDataAccessor interface {
	CreateOrder(ctx context.Context, arg *CreateOrderRequest) (*CreateOrderResponse, error)
	createOrderItems(ctx context.Context, order_id uuid.UUID, arg []CreateOrderItemRequest, tx *gorm.DB) ([]OrderItem, error)
}

type orderDataAccessor struct {
	database             Database
	l                    logger.Interface
	catalogServiceClient catalog_service.CatalogServiceClient
	redisInit            *redisOrder.Redis
	redisClient          *redis.Client
}

func NewOrderDataAccessor(database Database,
	l logger.Interface,
	catalogServiceClient catalog_service.CatalogServiceClient,
	redisInit *redisOrder.Redis,
	redisClient *redis.Client,
) OrderDataAccessor {
	return &orderDataAccessor{
		database:             database,
		l:                    l,
		catalogServiceClient: catalogServiceClient,
		redisInit:            redisInit,
		redisClient:          redisClient,
	}
}

func (o *orderDataAccessor) CreateOrder(ctx context.Context, arg *CreateOrderRequest) (*CreateOrderResponse, error) {
	var response *CreateOrderResponse

	err := o.database.Transaction(func(tx *gorm.DB) error {
		userID, err := uuid.Parse(arg.UserID)
		if err != nil {
			return err
		}

		order := &Order{
			UserID:  userID,
			Address: arg.Address,
			Status:  arg.Status,
		}
		if err := tx.WithContext(ctx).Create(&order).Error; err != nil {
			return err
		}

		orderItems, err := o.createOrderItems(ctx, order.ID, arg.OrderItems, tx)

		response = &CreateOrderResponse{
			Order: Order{
				ID:            order.ID,
				UserID:        order.UserID,
				Address:       order.Address,
				Status:        order.Status,
				PaymentMethod: order.PaymentMethod,
				CreatedAt:     order.CreatedAt,
				UpdatedAt:     order.UpdatedAt,
			},
			OrderItems: orderItems,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, err
}

func (o *orderDataAccessor) createOrderItems(ctx context.Context, order_id uuid.UUID, arg []CreateOrderItemRequest, tx *gorm.DB) ([]OrderItem, error) {
	orderItems := make([]OrderItem, 0, len(arg))

	for _, item := range arg {
		lockKey := fmt.Sprintf("lock_order_%w", item.SkuID)
		lockTimeout := 10 * time.Second
		lockValue, err := helper.GenerateRandomValue()
		if err != nil {
			return nil, fmt.Errorf("failed to generate lock value", err)
		}

		o.l.Info("lock info", zap.String("lock_key", lockKey), zap.String("lock_value", lockValue), zap.String("time_out", lockTimeout.String()))

		o.l.Info("start acquirelock")
		if ok := o.redisInit.AcquireLock(o.redisClient, lockKey, lockValue, lockTimeout); !ok {
			o.l.Error("failed to acquire lock")
			return nil, fmt.Errorf("failed to acquire lock")
		} else {
			o.l.Info("acquired lock success")
			// process order
			skuID, err := uuid.Parse(item.SkuID)
			if err != nil {
				return nil, fmt.Errorf("failed to parse uuid with sku id")
			}

			orderItems = append(orderItems, OrderItem{
				SkuID:    skuID,
				OrderID:  order_id,
				Quantity: item.Quantity,
			})
			o.catalogServiceClient.UpdateInventorySKU(ctx, &catalog_service.UpdateInventorySKURequest{
				SkuId:    item.SkuID,
				Quantity: item.Quantity,
			})

			// send event payment service => send event verify order

			o.l.Info("start release lock")
			// get lock value
			val, err := o.redisClient.Get(ctx, lockKey).Result()
			if err != nil {
				o.l.Error("failed to get value lock")
				return nil, fmt.Errorf("failed to get value lock")
			}

			// if val == lockKey => release lock
			o.l.Info("macth val && lockKey", zap.String("val", val), zap.String("lock_val", lockValue))
			if val == lockValue {
				if err := o.redisInit.ReleaseLock(o.redisClient, lockKey); err != nil {
					o.l.Error("failed to delete lock ", zap.Error(err))
					return nil, fmt.Errorf("failed to delete lock ", err)
				}

				o.l.Info("release lock successfully")
			} else {
				o.l.Info("lock is not match")
				// TODO
			}
		}
	}

	if err := tx.WithContext(ctx).Create(&orderItems).Error; err != nil {
		return nil, err
	}

	return orderItems, nil
}
