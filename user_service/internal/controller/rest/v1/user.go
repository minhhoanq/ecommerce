package v1

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/minhhoanq/ecommerce/user_service/internal/controller/rest/v1/middleware"
	"github.com/minhhoanq/ecommerce/user_service/internal/entity"
	"github.com/minhhoanq/ecommerce/user_service/internal/token"
	usecase "github.com/minhhoanq/ecommerce/user_service/internal/usecase/rest"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type userRoutes struct {
	u usecase.UserUsecase
	l logger.Interface
}

func newUserRoutes(handler *echo.Group, u usecase.UserUsecase, l logger.Interface, tokenMaker token.Maker) {
	r := &userRoutes{u: u, l: l}

	handler.POST("/login", r.login)
	handler.POST("/renew_access", r.createUser)
	handler.POST("", r.createUser)
	handler.Use(middleware.AuthMiddleware(tokenMaker))
	handler.GET("/:id", r.getUserByID)

}

func newUserResponse(user *entity.User) userResponse {
	return userResponse{
		ID:               user.ID.String(),
		Username:         user.Username,
		Email:            user.Email,
		RoleId:           user.RoleId,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
	}
}

type userResponse struct {
	ID               string    `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	RoleId           int       `json:"role_id"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func (r *userRoutes) getUserByID(c echo.Context) error {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := r.u.GetUserByID(c.Request().Context(), id)
	if err != nil {
		r.l.Error(err.Error())
		return c.JSON(http.StatusNotFound, errorResponse(err))
	}

	res := newUserResponse(user)

	return c.JSON(http.StatusOK, res)
}

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	RoleId   int    `json:"role_id"`
}

func (r *userRoutes) createUser(c echo.Context) error {
	var req CreateUserParams

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := r.u.CreateUser(c.Request().Context(), usecase.CreateUserUsecaseParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		RoleId:   req.RoleId,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResponse(err))
	}

	res := newUserResponse(user)

	return c.JSON(http.StatusOK, res)
}

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	SessionID             int          `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (r *userRoutes) login(c echo.Context) error {
	var req loginParams

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := r.u.Login(c.Request().Context(), usecase.LoginUsecaseParams{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		r.l.Error(err.Error())
		return c.JSON(http.StatusNotFound, errorResponse(err))

	}

	res := loginUserResponse{
		SessionID:             result.SessionID,
		AccessToken:           result.AccessToken,
		RefreshToken:          result.RefreshToken,
		AccessTokenExpiresAt:  result.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: result.RefreshTokenExpiresAt,
		User:                  newUserResponse(result.User),
	}

	return c.JSON(http.StatusOK, res)
}

type renewAccessTokenParams struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (r *userRoutes) renewAccessToken(c echo.Context) error {
	var req renewAccessTokenParams

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	result, err := r.u.RenewAccessToken(c.Request().Context(), usecase.RenewAccessTokenUsecaseParams{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, renewAccessTokenResponse{
		AccessToken:          result.AccessToken,
		AccessTokenExpiresAt: result.AccessTokenExpiresAt,
	})
}
