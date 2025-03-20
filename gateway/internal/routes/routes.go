package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/catalog"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/order"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/user"
	"github.com/minhhoanq/ecommerce/gateway/internal/token"
)

func NewRouter(e *echo.Echo, l logger.Interface, ordeManagementOperator order.OrderManagementOperator, catalogManagementOperator catalog.CatalogManagementOperator, userManagementOperator user.UserManagementOperator, tokenMaker token.Maker) {
	h := e.Group("/v1")
	{
		orderHandler := h.Group("/orders")
		{
			newOrderRouter(orderHandler, l, ordeManagementOperator)
		}
		catalogHandler := h.Group("/catalogs")
		{
			newCatalogRouter(catalogHandler, l, catalogManagementOperator)
		}

		userHandler := h.Group("/users")
		{
			newUserRouter(userHandler, l, userManagementOperator, tokenMaker)
		}
	}
}
