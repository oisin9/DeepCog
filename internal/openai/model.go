package openai

import (
	"encoding/json"
	"fmt"

	ai "github.com/sashabaranov/go-openai"
)

type Model struct {
	Id         string       `json:"id"`
	Object     string       `json:"object"`
	Created    int          `json:"created"`
	OwnedBy    string       `json:"owned_by"`
	Permission []Permission `json:"permission"`
}

type Permission struct {
	Id                 string   `json:"id"`
	Object             string   `json:"object"`
	Created            int      `json:"created"`
	AllowCreateEngine  bool     `json:"allow_create_engine"`
	AllowSampling      bool     `json:"allow_sampling"`
	AllowLogprobs      bool     `json:"allow_logprobs"`
	AllowSearchIndices bool     `json:"allow_search_indices"`
	AllowView          bool     `json:"allow_view"`
	AllowFineTuning    bool     `json:"allow_fine_tuning"`
	Organization       string   `json:"organization"`
	Group              []string `json:"group"`
	IsBlocking         bool     `json:"is_blocking"`
}

func NewModel(id string, owned_by string) *Model {
	permission := Permission{
		Id:                 "",
		Object:             "permission",
		Created:            1737331200,
		AllowCreateEngine:  false,
		AllowSampling:      true,
		AllowLogprobs:      true,
		AllowSearchIndices: false,
		AllowView:          true,
		AllowFineTuning:    false,
		Organization:       "*",
		Group:              nil,
		IsBlocking:         false,
	}
	return &Model{
		Id:         id,
		Object:     "model",
		Created:    1737331200,
		OwnedBy:    owned_by,
		Permission: []Permission{permission},
	}
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type MessageContent struct {
	StringValue string        `json:"-"`
	ArrayValue  []ContentItem `json:"-"`
	IsString    bool          `json:"-"`
}

func (m *MessageContent) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		m.StringValue = str
		m.IsString = true
		return nil
	}
	var arr []ContentItem
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("failed to unmarshal MessageContent: %w", err)
	}
	m.ArrayValue = arr
	m.IsString = false
	return nil
}

type Message struct {
	Role    string         `json:"role" query:"role" form:"role" valid:"role"`
	Content MessageContent `json:"content" query:"content" form:"content" valid:"content"`
}

func (m *Message) GetContent() string {
	if m.Content.IsString {
		return m.Content.StringValue
	}
	result := ""
	for _, item := range m.Content.ArrayValue {
		if item.Type == "text" {
			result += item.Text
		}
	}
	return result
}

func (m *Message) SetContent(content string) {
	m.Content.StringValue = content
	m.Content.IsString = true
}

type ChatChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index        int    `json:"index"`
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

func (m *Message) ChatMessage() *ai.ChatCompletionMessage {
	role := ai.ChatMessageRoleUser
	if m.Role == "system" {
		role = ai.ChatMessageRoleSystem
	}
	if m.Role == "assistant" {
		role = ai.ChatMessageRoleAssistant
	}
	return &ai.ChatCompletionMessage{
		Role:    role,
		Content: m.GetContent(),
	}
}

type ChatCompletionStreamChoiceDelta struct {
	ai.ChatCompletionStreamChoiceDelta
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

type ChatCompletionStreamChoice struct {
	ai.ChatCompletionStreamChoice
	Delta ChatCompletionStreamChoiceDelta `json:"delta"`
}

type ChatCompletionStreamResponse struct {
	ai.ChatCompletionStreamResponse
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

type ChatStream struct {
	Stream *ai.ChatCompletionStream
}

func (c *ChatStream) Recv() (*ChatCompletionStreamResponse, error) {
	resp := ChatCompletionStreamResponse{}
	rawLine, err := c.Stream.RecvRaw()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawLine, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *ChatStream) Close() {
	c.Stream.Close()
}
