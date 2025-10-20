package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skymoore/vibe-zsh/internal/config"
	vibeErrors "github.com/skymoore/vibe-zsh/internal/errors"
	"github.com/skymoore/vibe-zsh/internal/parser"
	"github.com/skymoore/vibe-zsh/internal/schema"
)

type Client struct {
	config     *config.Config
	httpClient *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

type JSONSchema struct {
	Name   string                 `json:"name"`
	Strict string                 `json:"strict"`
	Schema map[string]interface{} `json:"schema"`
}

type ChatCompletionRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	Temperature    float64         `json:"temperature,omitempty"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
	Stream         bool            `json:"stream"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func New(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *Client) GenerateCommand(ctx context.Context, query string) (*schema.CommandResponse, error) {
	if c.config.UseStructuredOutput {
		resp, err := c.generateWithStructuredOutput(ctx, query)
		if err == nil {
			return resp, nil
		}
	}

	return c.generateWithTextParsing(ctx, query)
}

func (c *Client) generateWithStructuredOutput(ctx context.Context, query string) (*schema.CommandResponse, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: schema.SystemPrompt,
		},
		{
			Role:    "user",
			Content: query,
		},
	}

	req := ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    messages,
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
		Stream:      false,
		ResponseFormat: &ResponseFormat{
			Type: "json_schema",
			JSONSchema: &JSONSchema{
				Name:   "shell_command_response",
				Strict: "true",
				Schema: schema.GetJSONSchema(),
			},
		},
	}

	chatResp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	content := chatResp.Choices[0].Message.Content

	var cmdResp schema.CommandResponse
	if err := json.Unmarshal([]byte(content), &cmdResp); err != nil {
		return nil, fmt.Errorf("%w: failed to parse JSON content", vibeErrors.ErrInvalidJSON)
	}

	return &cmdResp, nil
}

func (c *Client) generateWithTextParsing(ctx context.Context, query string) (*schema.CommandResponse, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: schema.SystemPrompt,
		},
		{
			Role:    "user",
			Content: query,
		},
	}

	req := ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    messages,
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
		Stream:      false,
	}

	chatResp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	content := chatResp.Choices[0].Message.Content

	cmdResp, err := parser.ParseTextResponse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse text response: %w", err)
	}

	return cmdResp, nil
}
