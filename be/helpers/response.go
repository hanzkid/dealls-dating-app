package helpers

import "github.com/labstack/echo/v4"

func ResponseWithError(ctx echo.Context, code int, message string) {
	ctx.JSON(code, map[string]string{"error": message})
}

func ResponseWithSuccess(ctx echo.Context, code int, data interface{}) {
	ctx.JSON(code, map[string]interface{}{"data": data})
}
