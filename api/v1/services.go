package v1

import (
	"connor.run/deepcog/internal/openai"
)

type GetModelsResponse struct {
	Object string          `json:"object"`
	Data   []*openai.Model `json:"data"`
}

type ChatCompletionRequest struct {
	Model       string           `json:"model" query:"model" form:"model" valid:"model"`
	Messages    []openai.Message `json:"messages" query:"messages" form:"messages" valid:"messages"`
	Temperature float32          `json:"temperature" query:"temperature" form:"temperature" valid:"temperature"`
	Stream      bool             `json:"stream" query:"stream" form:"stream" valid:"stream"`
	Stop        string           `json:"stop" query:"stop" form:"stop" valid:"stop"`
	MaxTokens   int              `json:"max_tokens" query:"max_tokens" form:"max_tokens" valid:"max_tokens"`
}
