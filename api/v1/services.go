package v1

import (
	"connor.run/deepcog/internal/openai"
)

type GetModelsResponse struct {
	Object string          `json:"object"`
	Data   []*openai.Model `json:"data"`
}
