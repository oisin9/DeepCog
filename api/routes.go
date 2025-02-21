package api

import (
	v1 "connor.run/deepcog/api/v1"
	"github.com/labstack/echo/v4"
)

type ServerHealthStatus struct {
	Status string `json:"status"`
}

func SetupRouter(e *echo.Echo) {
	v1Group := e.Group("/api/v1")
	v1.SetupV1Routes(v1Group)

	// Add health check route
	e.GET("/health", func(c echo.Context) error {
		status := &ServerHealthStatus{
			Status: "OK",
		}
		return c.JSON(200, status)
	})
}
