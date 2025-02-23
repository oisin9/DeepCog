package openai

import (
	"context"

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

func (c *Client) GetThinkingStream(msg []Message) (*ChatStream, error) {
	config := ai.DefaultConfig(c.ThinkingModel.ApiKey)
	config.BaseURL = c.ThinkingModel.BaseUrl
	aiClient := ai.NewClientWithConfig(config)
	ctx := context.Background()
	msgs := []ai.ChatCompletionMessage{}
	for _, message := range msg {
		msgs = append(msgs, *message.ChatMessage())
	}
	req := ai.ChatCompletionRequest{
		Model:     c.ThinkingModel.ModelName,
		Messages:  msgs,
		MaxTokens: c.MaxTokens,
		Stream:    true,
	}
	stream, err := aiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ChatStream{
		Stream: stream,
	}, nil
}

func (c *Client) GetGenerateStream(msg []Message, reasoningContent string) (*ChatStream, error) {
	config := ai.DefaultConfig(c.GenerateModel.ApiKey)
	config.BaseURL = c.GenerateModel.BaseUrl
	aiClient := ai.NewClientWithConfig(config)
	ctx := context.Background()
	// 把最后的一条role为user的消息替换成 reasoningContent
	if reasoningContent != "" {
		for i := len(msg) - 1; i >= 0; i-- {
			if msg[i].Role == "user" {
				msg[i].Content = reasoningContent
				break
			}
		}
	}
	msgs := []ai.ChatCompletionMessage{}
	for _, message := range msg {
		msgs = append(msgs, *message.ChatMessage())
	}
	req := ai.ChatCompletionRequest{
		Model:     c.GenerateModel.ModelName,
		Messages:  msgs,
		MaxTokens: c.MaxTokens,
		Stream:    true,
	}
	stream, err := aiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ChatStream{
		Stream: stream,
	}, nil
}

func (client *Client) ChatStream(c echo.Context, messages []Message) error {
	sseClient := utils.NewSSEClient(c)
	err := sseClient.Init()
	if err != nil {
		return err
	}
	defer sseClient.Close()
	id := "chatcmpl-" + uuid.New().String()
	reasoning_content := ""
	if client.ThinkingModel != nil {
		thinkingStream, err := client.GetThinkingStream(messages)
		if err != nil {
			return err
		}
		for {
			response, err := thinkingStream.Recv()
			if err != nil {
				thinkingStream.Close()
				break
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
				thinkingStream.Close()
			}
		}
	}
	generateStream, err := client.GetGenerateStream(messages, reasoning_content)
	if err != nil {
		return err
	}
	defer generateStream.Close()
	for {
		response, err := generateStream.Recv()
		if err != nil {
			break
		}
		response.ID = id
		response.Model = client.Model.Id
		if reasoning_content != "" && response.Choices[0].Delta.ReasoningContent != "" {
			continue
		}
		if response.Choices[0].FinishReason == "stop" {
			break
		}
		sseClient.SendResponse(response)
	}
	return nil
}
