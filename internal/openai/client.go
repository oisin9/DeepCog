package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"connor.run/deepcog/pkg/config"
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
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	w := c.Response().Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return c.JSON(500, "Streaming unsupported!")
	}
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
			return c.JSON(500, err)
		}
		reasoning_content := ""
		for {
			response := ChatCompletionStreamResponse{}
			rawLine, err := thinkingStream.RecvRaw()
			if err != nil {
				fmt.Fprint(w, "data: [DONE]\n\n")
				flusher.Flush()
				thinkingStream.Close()
				return err
			}
			err = json.Unmarshal(rawLine, &response)
			if errors.Is(err, io.EOF) {
				thinkingStream.Close()
				break
			}
			if err != nil {
				fmt.Fprint(w, "data: [DONE]\n\n")
				flusher.Flush()
				thinkingStream.Close()
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
				data, err := json.Marshal(response)
				if err != nil {
					fmt.Fprint(w, "data: [DONE]\n\n")
					flusher.Flush()
					thinkingStream.Close()
					return err
				}
				fmt.Fprintf(w, "data: %s\n\n", string(data))
				flusher.Flush()
			} else {
				data, err := json.Marshal(response)
				if err != nil {
					fmt.Fprint(w, "data: [DONE]\n\n")
					flusher.Flush()
					thinkingStream.Close()
					return err
				}
				fmt.Fprintf(w, "data: %s\n\n", string(data))
				flusher.Flush()
			}
		}
		fmt.Println(reasoning_content)
	}
	return nil
}
