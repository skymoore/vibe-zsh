package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skymoore/vibe-zsh/internal/cache"
	"github.com/skymoore/vibe-zsh/internal/config"
	vibeErrors "github.com/skymoore/vibe-zsh/internal/errors"
	"github.com/skymoore/vibe-zsh/internal/logger"
	"github.com/skymoore/vibe-zsh/internal/parser"
	"github.com/skymoore/vibe-zsh/internal/progress"
	"github.com/skymoore/vibe-zsh/internal/schema"
)

type Client struct {
	config     *config.Config
	httpClient *http.Client
	cache      *cache.Cache
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
	client := &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}

	if cfg.EnableCache {
		c, err := cache.New(cfg.CacheDir, cfg.CacheTTL)
		if err == nil {
			client.cache = c
		}
	}

	return client
}

func (c *Client) GenerateCommand(ctx context.Context, query string) (*schema.CommandResponse, error) {
	// Initialize spinner if progress is enabled and stderr is a terminal
	var spinner *progress.Spinner
	if c.config.ShowProgress && progress.IsStderrTerminal() {
		spinner = progress.NewSpinner(c.config.ProgressStyle)
		defer spinner.Stop()
	}

	// Check cache first
	if spinner != nil {
		spinner.Start(ctx, "Checking cache...")
	}

	if c.cache != nil {
		if cached, ok := c.cache.Get(query); ok {
			// Cache hit - stop spinner immediately
			if spinner != nil {
				spinner.Stop()
			}
			return cached, nil
		}
	}

	// Update spinner for API call
	if spinner != nil {
		spinner.Update("Contacting API...")
	}

	var resp *schema.CommandResponse
	var err error

	if c.config.UseStructuredOutput {
		if spinner != nil {
			spinner.Update("Generating command...")
		}
		resp, err = c.generateWithStructuredOutput(ctx, query)
		if err == nil && c.config.StrictValidation {
			if validErr := resp.Validate(); validErr == nil {
				c.cacheIfEnabled(query, resp)
				logger.LogLayerSuccess("structured_output", 1)
				return resp, nil
			}
		} else if err == nil {
			c.cacheIfEnabled(query, resp)
			logger.LogLayerSuccess("structured_output", 1)
			return resp, nil
		}
		logger.LogParsingFailure(1, "structured_output", "", err)
	}

	if spinner != nil {
		spinner.Update("Parsing response...")
	}
	resp, err = c.generateWithEnhancedParsing(ctx, query, spinner)
	if err == nil {
		c.cacheIfEnabled(query, resp)
		logger.LogLayerSuccess("enhanced_parsing", 2)
		return resp, nil
	}
	logger.LogParsingFailure(2, "enhanced_parsing", "", err)

	if spinner != nil {
		spinner.Update("Retrying with explicit JSON...")
	}
	resp, err = c.generateWithExplicitJSONPrompt(ctx, query)
	if err == nil {
		c.cacheIfEnabled(query, resp)
		logger.LogLayerSuccess("explicit_json_prompt", 3)
		return resp, nil
	}
	logger.LogParsingFailure(3, "explicit_json_prompt", "", err)

	if spinner != nil {
		spinner.Update("Using fallback...")
	}

	// Pass the last error to the fallback so it can provide better feedback
	resp, fallbackErr := c.generateWithEmergencyFallback(ctx, query, err)
	if fallbackErr == nil {
		logger.LogLayerSuccess("emergency_fallback", 4)
		return resp, nil
	}
	logger.LogParsingFailure(4, "emergency_fallback", "", fallbackErr)

	return nil, fmt.Errorf("all parsing strategies failed: %w", err)
}

func (c *Client) cacheIfEnabled(query string, resp *schema.CommandResponse) {
	if c.cache != nil {
		if err := c.cache.Set(query, resp); err != nil {
			logger.Debug("Failed to cache response: %v", err)
		}
	}
}

func (c *Client) generateWithStructuredOutput(ctx context.Context, query string) (*schema.CommandResponse, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: schema.GetSystemPrompt(c.config.OSName, c.config.Shell),
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

func (c *Client) generateWithEnhancedParsing(ctx context.Context, query string, spinner *progress.Spinner) (*schema.CommandResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
		if spinner != nil && c.config.MaxRetries > 1 {
			spinner.Update(fmt.Sprintf("Parsing response (attempt %d/%d)...", attempt, c.config.MaxRetries))
		}
		messages := []Message{
			{
				Role:    "system",
				Content: schema.GetSystemPrompt(c.config.OSName, c.config.Shell),
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
			lastErr = fmt.Errorf("attempt %d: request failed: %w", attempt, err)
			logger.LogParsingFailure(attempt, "enhanced_parsing_request", "", lastErr)
			continue
		}

		content := chatResp.Choices[0].Message.Content

		if c.config.EnableJSONExtraction {
			cleanedJSON, err := parser.ExtractJSON(content)
			if err != nil {
				lastErr = fmt.Errorf("attempt %d: JSON extraction failed: %w", attempt, err)
				logger.LogParsingFailure(attempt, "enhanced_parsing_extraction", content, lastErr)
				continue
			}

			var cmdResp schema.CommandResponse
			if err := json.Unmarshal([]byte(cleanedJSON), &cmdResp); err != nil {
				lastErr = fmt.Errorf("attempt %d: unmarshal failed: %w", attempt, err)
				logger.LogParsingFailure(attempt, "enhanced_parsing_unmarshal", cleanedJSON, lastErr)
				continue
			}

			if c.config.StrictValidation {
				if err := cmdResp.Validate(); err != nil {
					lastErr = fmt.Errorf("attempt %d: validation failed: %w", attempt, err)
					logger.LogParsingFailure(attempt, "enhanced_parsing_validation", cleanedJSON, lastErr)
					continue
				}
			}

			return &cmdResp, nil
		}

		cmdResp, err := parser.ParseTextResponse(content)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: text parsing failed: %w", attempt, err)
			logger.LogParsingFailure(attempt, "enhanced_parsing_text", content, lastErr)
			continue
		}

		if c.config.StrictValidation {
			if err := cmdResp.Validate(); err != nil {
				lastErr = fmt.Errorf("attempt %d: validation failed: %w", attempt, err)
				logger.LogParsingFailure(attempt, "enhanced_parsing_validation", content, lastErr)
				continue
			}
		}

		return cmdResp, nil
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", c.config.MaxRetries, lastErr)
}

func (c *Client) generateWithExplicitJSONPrompt(ctx context.Context, query string) (*schema.CommandResponse, error) {
	explicitPrompt := schema.GetSystemPrompt(c.config.OSName, c.config.Shell) + "\n\nREMINDER: Your response must START with { and END with }. Nothing else."

	messages := []Message{
		{
			Role:    "system",
			Content: explicitPrompt,
		},
		{
			Role:    "user",
			Content: query,
		},
	}

	req := ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    messages,
		Temperature: c.config.Temperature * 0.5,
		MaxTokens:   c.config.MaxTokens,
		Stream:      false,
	}

	chatResp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	content := chatResp.Choices[0].Message.Content

	if c.config.EnableJSONExtraction {
		cleanedJSON, err := parser.ExtractJSON(content)
		if err != nil {
			logger.LogParsingFailure(1, "explicit_json_extraction", content, err)
			return nil, fmt.Errorf("JSON extraction failed: %w", err)
		}

		var cmdResp schema.CommandResponse
		if err := json.Unmarshal([]byte(cleanedJSON), &cmdResp); err != nil {
			logger.LogParsingFailure(1, "explicit_json_unmarshal", cleanedJSON, err)
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		if c.config.StrictValidation {
			if err := cmdResp.Validate(); err != nil {
				logger.LogParsingFailure(1, "explicit_json_validation", cleanedJSON, err)
				return nil, fmt.Errorf("validation failed: %w", err)
			}
		}

		return &cmdResp, nil
	}

	var cmdResp schema.CommandResponse
	if err := json.Unmarshal([]byte(content), &cmdResp); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	if c.config.StrictValidation {
		if err := cmdResp.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}

	return &cmdResp, nil
}

func (c *Client) generateWithEmergencyFallback(ctx context.Context, query string, lastErr error) (*schema.CommandResponse, error) {
	explanation := []string{
		fmt.Sprintf("Vibe failed to generate a valid command after %d attempts.", c.config.MaxRetries),
	}

	// Add specific error information if available
	if lastErr != nil {
		explanation = append(explanation, fmt.Sprintf("Error: %v", lastErr))
	}

	explanation = append(explanation, "Try rephrasing your request or report at: https://github.com/skymoore/vibe-zsh/issues")

	return &schema.CommandResponse{
		Command:     "",
		Explanation: explanation,
		Warning:     "Failed to generate command",
	}, nil
}
