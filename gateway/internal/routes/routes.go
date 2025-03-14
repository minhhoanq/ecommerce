package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/order"
)

func NewRouter(e *echo.Echo, l logger.Interface, ordeManagementOperator order.OrderManagementOperator) {
	h := e.Group("/v1")
	{
		orderHandler := h.Group("/orders")
		{
			newOrderRouter(orderHandler, l, ordeManagementOperator)
		}
	}
}
