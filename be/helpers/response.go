package helpers

import "github.com/labstack/echo/v4"

func ResponseWithError(ctx echo.Context, code int, message string) {
	err := ctx.JSON(code, map[string]string{"error": message})
	if err != nil {
		ctx.Logger().Error(err)
	}
}

func ResponseWithSuccess(ctx echo.Context, code int, data interface{}) {
	err := ctx.JSON(code, map[string]interface{}{"data": data})
	if err != nil {
		ctx.Logger().Error(err)
	}

}
