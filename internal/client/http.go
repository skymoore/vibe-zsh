package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	vibeErrors "github.com/skymoore/vibe-zsh/internal/errors"
)

func (c *Client) doRequest(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	var chatResp *ChatCompletionResponse
	var lastErr error

	err := c.withRetry(ctx, func() error {
		reqBody, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		url := strings.TrimSuffix(c.config.APIURL, "/") + "/chat/completions"
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		httpReq.Header.Set("Content-Type", "application/json")
		if c.config.APIKey != "" {
			httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return fmt.Errorf("%w: %v", vibeErrors.ErrTimeout, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return vibeErrors.NewAPIError(resp.StatusCode, string(body))
		}

		var tempResp ChatCompletionResponse
		if err := json.Unmarshal(body, &tempResp); err != nil {
			return fmt.Errorf("%w: %v", vibeErrors.ErrInvalidJSON, err)
		}

		if len(tempResp.Choices) == 0 {
			return vibeErrors.ErrNoResponse
		}

		chatResp = &tempResp
		return nil
	})

	if err != nil {
		lastErr = err
	}

	if chatResp == nil {
		if lastErr != nil {
			return nil, lastErr
		}
		return nil, vibeErrors.ErrEmptyResponse
	}

	return chatResp, nil
}
