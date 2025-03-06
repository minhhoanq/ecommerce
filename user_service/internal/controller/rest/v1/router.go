package v1

import (
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/internal/token"
	usecase "github.com/minhhoanq/ecommerce/user_service/internal/usecase/rest"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, l logger.Interface, u usecase.UserUsecase, tokenMaker token.Maker) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := e.Group("/v1")
	{
		usersHandler := h.Group("/users")
		{
			newUserRoutes(usersHandler, u, l, tokenMaker)
		}
	}
}

func errorResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}
