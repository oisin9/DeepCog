package openai

import (
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

type Message struct {
	Role    string `json:"role" query:"role" form:"role" valid:"role"`
	Content string `json:"content" query:"content" form:"content" valid:"content"`
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
		Content: m.Content,
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
