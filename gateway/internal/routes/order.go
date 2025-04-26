package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/order_service"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/order"
)

type orderHandlerFunc struct {
	l                       logger.Interface
	orderManagementOperator order.OrderManagementOperator
}

func newOrderRouter(handler *echo.Group, l logger.Interface, orderManagementOperator order.OrderManagementOperator) {
	o := &orderHandlerFunc{l, orderManagementOperator}

	handler.POST("", o.createOrder)
}

type CartItems struct {
	SkuID    string `json:"sku_id"`
	Quantity int32  `json:"quantity"`
}

type CreateOrderParams struct {
	UserID    string      `json:"user_id"`
	CartItems []CartItems `json:"cart_items"`
}

func (o *orderHandlerFunc) createOrder(c echo.Context) error {
	fmt.Println("create order")
	var req CreateOrderParams

	if err := c.Bind(&req); err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	arg := &order_service.CreateOrderRequest{
		UserId:    req.UserID,
		CartItems: make([]*order_service.CartItems, 0, len(req.CartItems)),
	}

	for _, item := range req.CartItems {
		arg.CartItems = append(arg.CartItems, &order_service.CartItems{
			SkuId:    item.SkuID,
			Quantity: item.Quantity,
		})
	}

	order, err := o.orderManagementOperator.CreateOrder(c.Request().Context(), arg)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return Error(c, sql.ErrNoRows, http.StatusNotFound)
		}

		return Error(c, err, http.StatusBadRequest)
	}

	return Success(c, order, http.StatusCreated, nil)
}
