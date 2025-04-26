package routes

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

// parseQueryIntParam gets a query parameter from Echo context, parses it to int32,
// and falls back to a default value if it's missing or invalid.
func parseQueryIntParam(c echo.Context, key string, defaultVal, min, max int) int32 {
	valStr := c.QueryParam(key)
	val, err := strconv.Atoi(valStr)
	if err != nil || val < min || val > max {
		return int32(defaultVal)
	}
	return int32(val)
}

func parsePathIntParam(c echo.Context, key string) (int32, error) {
	valStr := c.Param(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid path param: %s", key)
	}
	return int32(val), nil
}

func parsePathParam(c echo.Context, key string) (string, error) {
	valStr := c.Param(key)
	if valStr == "" {
		return "", fmt.Errorf("invalid path param: %s", key)
	}
	return valStr, nil
}

func parseJSONFormField[T any](c echo.Context, field string, out *T) error {
	jsonStr := c.FormValue(field)
	if jsonStr == "" {
		return fmt.Errorf("form field %q is missing", field)
	}

	if err := json.Unmarshal([]byte(jsonStr), out); err != nil {
		return fmt.Errorf("invalid JSON in form field %q: %w", field, err)
	}

	return nil
}

// Bind body nếu Content-Type là application/json
func bodyBind[T any](c echo.Context) error {
	var req T

	if c.Request().Header.Get("Content-Type") == "application/json" {
		if err := c.Bind(&req); err != nil {
			return fmt.Errorf("failed to bind JSON body: %w", err)
		}
	}

	return nil
}
