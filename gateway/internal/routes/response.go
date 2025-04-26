package routes

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Success gửi phản hồi thành công với dữ liệu JSON, status code và headers tuỳ chọn.
func Success(c echo.Context, data interface{}, status int, headers map[string]string) error {
	for k, v := range headers {
		c.Response().Header().Set(k, v)
	}
	return c.JSON(status, &Response{
		Data: data,
	})
}

// Error gửi phản hồi lỗi với message JSON chuẩn hóa.
func Error(c echo.Context, err error, status int) error {
	return c.JSON(status, &Response{
		Message: err.Error(),
	})
}
