package client

import (
	"context"
	"time"

	vibeErrors "github.com/skymoore/vibe-zsh/internal/errors"
)

const (
	maxRetries     = 3
	initialBackoff = 1 * time.Second
	maxBackoff     = 10 * time.Second
)

func (c *Client) withRetry(ctx context.Context, fn func() error) error {
	var lastErr error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return ctx.Err()
			}

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		if !vibeErrors.IsRetryable(err) {
			return err
		}
	}

	return lastErr
}
