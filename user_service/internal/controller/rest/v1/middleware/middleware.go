package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/minhhoanq/ecommerce/user_service/internal/token"

	"github.com/labstack/echo/v4"
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
				return c.JSON(http.StatusUnauthorized, "Authorization header is missing")
			}

			fields := strings.Fields(authorization)
			if len(fields) < 2 || strings.ToLower(fields[0]) != authorizationTypeBearer {
				return c.JSON(http.StatusUnauthorized, "Invalid authorization header format")
			}

			token := fields[1]
			payload, err := tokenMaker.VerifyToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err.Error())
			}

			fmt.Println("IssuedAt", payload.IssuedAt)
			fmt.Println("ExpiredAt", payload.ExpiredAt)

			if time.Now().After(payload.ExpiredAt) {
				return c.JSON(http.StatusUnauthorized, "expired session")
			}

			c.Set(authorizationPayloadKey, payload)

			return next(c)
		}
	}
}
