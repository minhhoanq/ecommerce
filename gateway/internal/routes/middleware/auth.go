package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/minhhoanq/ecommerce/gateway/internal/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get(authorizationHeaderKey)
			if len(authorization) == 0 {
				return c.JSON(http.StatusUnauthorized, "Authorization header is not provided")
			}

			fields := strings.Fields(authorization)
			if len(fields) < 2 {
				return c.JSON(http.StatusUnauthorized, "Invalid authorization header format")
			}

			token := fields[1]
			payload, err := tokenMaker.VerifyToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err.Error())
			}

			if time.Now().After(payload.ExpiredAt) {
				return c.JSON(http.StatusUnauthorized, "Token is expired")
			}

			c.Set(authorizationPayloadKey, payload)

			return next(c)
		}
	}
}
