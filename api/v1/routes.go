package v1

import (
	"github.com/labstack/echo/v4"
)

func SetupV1Routes(e *echo.Group) {
	e.GET("/models", GetModels)
	e.POST("/chat/completions", ChatCompletion)
}
