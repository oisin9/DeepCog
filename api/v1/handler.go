package v1

import (
	"strings"

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

func ChatCompletion(c echo.Context) error {
	req := &ChatCompletionRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(400, err.Error())
	}
	cfg := config.GetConfig()
	model := cfg.GetModel(req.Model)
	if model == nil {
		return c.JSON(400, "Model not found")
	}
	if model.ApiKey != "" {
		authHeader := c.Request().Header.Get("Authorization")
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")
		if apiKey != model.ApiKey {
			return c.JSON(400, "Invalid API key")
		}
	}
	thinking_model := cfg.GetBaseModel(model.ThinkingModel)
	generate_model := cfg.GetBaseModel(model.GenerateModel)
	client := openai.NewClient(thinking_model, generate_model, model)
	client.SetMaxTokens(req.MaxTokens)
	client.SetTemperature(req.Temperature)
	return client.ChatStream(c, req.Messages)
}
