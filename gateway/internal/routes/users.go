package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/modules/user"
	"github.com/minhhoanq/ecommerce/gateway/internal/routes/middleware"
	"github.com/minhhoanq/ecommerce/gateway/internal/token"
)

type userHandlerFunc struct {
	l                      logger.Interface
	userManagementOperator user.UserManagementOperator
}

func newUserRouter(handler *echo.Group, l logger.Interface, userManagementOperator user.UserManagementOperator, tokenMaker token.Maker) {
	u := &userHandlerFunc{
		l:                      l,
		userManagementOperator: userManagementOperator,
	}

	handler.POST("", u.signup)
	handler.POST("", u.signup, middleware.AuthMiddleware(tokenMaker))
}

func (userHandler *userHandlerFunc) signup(c echo.Context) error {
	return nil
}

func (userHandler *userHandlerFunc) signin(c echo.Context) error {
	return nil
}

func (userHanlder *userHandlerFunc) getUser(c echo.Context) error {

	return c.JSON(http.StatusOK, "get user")
}
