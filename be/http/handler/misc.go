package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Healthz(c echo.Context) error {
	return c.String(http.StatusOK, "OK!")
}