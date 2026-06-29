package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/skymoore/vibe-zsh/internal/cache"
	"github.com/skymoore/vibe-zsh/internal/config"
	"github.com/skymoore/vibe-zsh/internal/logger"
	"github.com/skymoore/vibe-zsh/internal/parser"
	"github.com/skymoore/vibe-zsh/internal/progress"
	"github.com/skymoore/vibe-zsh/internal/schema"
	"github.com/teilomillet/gollm"
)

// Client wraps a gollm LLM instance and the vibe response-parsing pipeline.
// All provider-specific transport, auth, and retry logic is handled by gollm;
// this layer is responsible only for prompt construction, the multi-strategy
// JSON parsing fallback, and caching.
type Client struct {
	config  *config.Config
	llm     gollm.LLM
	initErr error
	cache   *cache.Cache
}

// New constructs a Client. It builds the underlying gollm LLM from the
// resolved configuration (provider, model, key, generation params). If the
// LLM cannot be constructed, a Client with a nil llm is returned and the
// construction error is surfaced on the first GenerateCommand call.
func New(cfg *config.Config) *Client {
	client := &Client{config: cfg}

	llm, err := newLLM(cfg)
	if err != nil {
		client.initErr = err
		logger.Debug("Failed to initialize LLM provider: %v", err)
	} else {
		client.llm = llm
	}

	if cfg.EnableCache {
		if c, err := cache.New(cfg.CacheDir, cfg.CacheTTL); err == nil {
			client.cache = c
		}
	}

	return client
}

// newLLM translates the vibe Config into gollm options and builds the LLM.
func newLLM(cfg *config.Config) (gollm.LLM, error) {
	opts := []gollm.ConfigOption{
		gollm.SetProvider(cfg.Provider),
		gollm.SetModel(cfg.Model),
		gollm.SetTemperature(cfg.Temperature),
		gollm.SetMaxTokens(cfg.MaxTokens),
		gollm.SetTimeout(cfg.Timeout),
		gollm.SetMaxRetries(cfg.MaxRetries),
		gollm.SetRetryDelay(1 * time.Second),
		gollm.SetLogLevel(logLevel(cfg)),
	}

	if cfg.APIKey != "" {
		opts = append(opts, gollm.SetAPIKey(cfg.APIKey))
	}

	// Local providers reach a user-configured endpoint rather than a fixed
	// hosted URL. Pass VIBE_API_URL through so custom ports/hosts work.
	switch cfg.Provider {
	case "ollama":
		if cfg.APIURL != "" {
			// gollm's Ollama provider expects the native base URL, not the
			// OpenAI-compatible /v1 path, so strip a trailing /v1.
			opts = append(opts, gollm.SetOllamaEndpoint(ollamaEndpoint(cfg.APIURL)))
		}
	case "vllm":
		if cfg.APIURL != "" {
			opts = append(opts, gollm.SetVLLMEndpoint(cfg.APIURL))
		}
	}

	return gollm.NewLLM(opts...)
}

// isLocalProvider reports whether the provider runs locally and does not
// require API-key authentication.
func isLocalProvider(provider string) bool {
	switch provider {
	case "ollama", "lmstudio", "vllm":
		return true
	default:
		return false
	}
}

// notConfiguredError produces an actionable message explaining why the LLM
// could not be constructed. gollm validates the configuration up front: hosted
// providers require a correctly-formatted API key, while local providers must
// be reachable at construction time.
func (c *Client) notConfiguredError() error {
	hint := fmt.Sprintf("check VIBE_PROVIDER (%q), VIBE_MODEL (%q) and VIBE_API_KEY", c.config.Provider, c.config.Model)
	if isLocalProvider(c.config.Provider) {
		hint = fmt.Sprintf("ensure the %s server is running and reachable at %s", c.config.Provider, c.config.APIURL)
	}
	if c.initErr != nil {
		return fmt.Errorf("LLM provider %q is not configured correctly (%s): %w", c.config.Provider, hint, c.initErr)
	}
	return fmt.Errorf("LLM provider %q is not configured correctly - %s", c.config.Provider, hint)
}

func logLevel(cfg *config.Config) gollm.LogLevel {
	if cfg.EnableDebugLogs {
		return gollm.LogLevelDebug
	}
	return gollm.LogLevelError
}

// ollamaEndpoint normalizes an OpenAI-style Ollama URL (".../v1") to the
// native Ollama base URL that gollm expects.
func ollamaEndpoint(apiURL string) string {
	u := strings.TrimSuffix(apiURL, "/")
	u = strings.TrimSuffix(u, "/v1")
	return u
}

// generate runs a single completion through gollm and returns the raw text
// content. temperature lets individual strategies tune sampling (e.g. the
// explicit-JSON retry lowers it). Provider transport and retries are handled
// inside gollm.
func (c *Client) generate(ctx context.Context, systemPrompt, query string, temperature float64) (string, error) {
	if c.llm == nil {
		return "", c.notConfiguredError()
	}

	prompt := gollm.NewPrompt(
		query,
		gollm.WithSystemPrompt(systemPrompt, gollm.CacheTypeEphemeral),
	)

	// gollm exposes generation parameters as provider options rather than
	// per-call GenerateOptions, so set temperature on the instance before the
	// call. vibe runs a single sequential request per invocation, so mutating
	// the shared option here is safe.
	c.llm.SetOption("temperature", temperature)

	content, err := c.llm.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("empty response from provider %q", c.config.Provider)
	}

	return content, nil
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
	content, err := c.generate(ctx, schema.GetSystemPrompt(c.config.OSName, c.config.Shell), query, c.config.Temperature)
	if err != nil {
		return nil, err
	}

	var cmdResp schema.CommandResponse
	if err := json.Unmarshal([]byte(content), &cmdResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON content: %w", err)
	}

	return &cmdResp, nil
}

func (c *Client) generateWithEnhancedParsing(ctx context.Context, query string, spinner *progress.Spinner) (*schema.CommandResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
		if spinner != nil && c.config.MaxRetries > 1 {
			spinner.Update(fmt.Sprintf("Parsing response (attempt %d/%d)...", attempt, c.config.MaxRetries))
		}

		content, err := c.generate(ctx, schema.GetSystemPrompt(c.config.OSName, c.config.Shell), query, c.config.Temperature)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: request failed: %w", attempt, err)
			logger.LogParsingFailure(attempt, "enhanced_parsing_request", "", lastErr)
			continue
		}

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

	content, err := c.generate(ctx, explicitPrompt, query, c.config.Temperature*0.5)
	if err != nil {
		return nil, err
	}

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

func (c *Client) generateWithEmergencyFallback(_ context.Context, _ string, lastErr error) (*schema.CommandResponse, error) {
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
