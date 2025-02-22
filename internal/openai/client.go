package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"connor.run/deepcog/pkg/config"
	"connor.run/deepcog/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	ai "github.com/sashabaranov/go-openai"
)

type Client struct {
	Model         *config.Model
	ThinkingModel *config.BaseModel
	GenerateModel *config.BaseModel
	MaxTokens     int
	Temperature   float32
}

func NewClient(thinking_model, generate_model *config.BaseModel, model *config.Model) *Client {
	return &Client{
		Model:         model,
		ThinkingModel: thinking_model,
		GenerateModel: generate_model,
		MaxTokens:     8192,
		Temperature:   0.7,
	}
}

func (c *Client) SetMaxTokens(max_tokens int) {
	c.MaxTokens = max_tokens
}

func (c *Client) SetTemperature(temperature float32) {
	c.Temperature = temperature
}

func (client *Client) ChatStream(c echo.Context, messages []Message) error {
	sseClient := utils.NewSSEClient(c)
	err := sseClient.Init()
	if err != nil {
		return err
	}
	defer sseClient.Close()
	id := "chatcmpl-" + uuid.New().String()
	if client.ThinkingModel != nil {
		config := ai.DefaultConfig(client.ThinkingModel.ApiKey)
		config.BaseURL = client.ThinkingModel.BaseUrl
		aiClient := ai.NewClientWithConfig(config)
		ctx := context.Background()
		msgs := []ai.ChatCompletionMessage{}
		for _, message := range messages {
			msgs = append(msgs, *message.ChatMessage())
		}
		req := ai.ChatCompletionRequest{
			Model:     client.ThinkingModel.ModelName,
			Messages:  msgs,
			MaxTokens: client.MaxTokens,
			Stream:    true,
		}
		thinkingStream, err := aiClient.CreateChatCompletionStream(ctx, req)
		if err != nil {
			return err
		}
		defer thinkingStream.Close()
		reasoning_content := ""
		for {
			response := ChatCompletionStreamResponse{}
			rawLine, err := thinkingStream.RecvRaw()
			if err != nil {
				break
			}
			err = json.Unmarshal(rawLine, &response)
			if err != nil {
				return err
			}
			response.ID = id
			response.Model = client.Model.Id
			response.Choices[0].FinishReason = ""
			if response.Choices[0].Delta.Content == "" && response.Choices[0].Delta.ReasoningContent == "" {
				continue
			}
			if response.Choices[0].Delta.ReasoningContent != "" {
				reasoning_content += response.Choices[0].Delta.ReasoningContent
				sseClient.SendResponse(response)
			} else {
				sseClient.SendResponse(response)
			}
		}
		fmt.Println(reasoning_content)
	}
	return nil
}
