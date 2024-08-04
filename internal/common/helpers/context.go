package helpers

import (
	"context"

	"github.com/labstack/echo/v4"
)

func Context(ctx echo.Context) context.Context {
	return ctx.Request().Context()
}
