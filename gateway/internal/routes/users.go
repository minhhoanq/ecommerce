package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/gateway/internal/generated/user_service"
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
	handler.POST("/signin", u.signin)
	handler.GET("/:id", u.getUser, middleware.AuthMiddleware(tokenMaker))
}

func (userHandler *userHandlerFunc) signup(c echo.Context) error {
	var req user_service.SignupRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	user, err := userHandler.userManagementOperator.Signup(c.Request().Context(), &req)
	if err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	return Success(c, user, http.StatusCreated, nil)
}

func (userHandler *userHandlerFunc) signin(c echo.Context) error {
	var req user_service.SigninRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	user, err := userHandler.userManagementOperator.Signin(c.Request().Context(), &req)
	if err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	return Success(c, user, http.StatusOK, nil)
}

func (userHanlder *userHandlerFunc) getUser(c echo.Context) error {
	var req user_service.GetUserRequest
	id, err := parsePathParam(c, "id")
	if err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	if id != "" {
		req.Id = id
	}
	user, err := userHanlder.userManagementOperator.GetUser(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Error(c, err, http.StatusNotFound)
		}

		return Error(c, err, http.StatusBadRequest)

	}
	return Success(c, user, http.StatusOK, nil)
}
