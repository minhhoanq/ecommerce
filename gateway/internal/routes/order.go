package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/order"
	"go.uber.org/zap"
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
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	o.l.Info("order params: ", zap.String("order: ", req.CartItems[0].SkuID))

	return nil
}
