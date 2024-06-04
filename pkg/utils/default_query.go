package utils

import "github.com/labstack/echo/v4"

func DefaultQuery(ctx echo.Context, key string, defaultValue string) string {
	if value := ctx.QueryParam(key); value != "" {
		return value
	}

	return defaultValue
}
