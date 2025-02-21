package api

import (
	"github.com/labstack/echo/v4"
)

func NewServer() *echo.Echo {
	e := echo.New()

	SetupRouter(e)

	return e
}
