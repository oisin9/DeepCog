package v1

import (
	"connor.run/deepcog/internal/openai"
	"connor.run/deepcog/pkg/config"
	"github.com/labstack/echo/v4"
)

func GetModels(c echo.Context) error {
	cfg := config.GetConfig()
	models := make([]*openai.Model, 0)

	for _, model := range cfg.Models {
		models = append(models, openai.NewModel(model.Id, model.OwnedBy))
	}
	response := &GetModelsResponse{
		Object: "list",
		Data:   models,
	}
	return c.JSON(200, response)
}
