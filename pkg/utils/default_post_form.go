package utils

import "github.com/labstack/echo/v4"

func DefaultPostForm(ctx echo.Context, key, defaultValue string) string {
	value := ctx.Request().Form.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}
