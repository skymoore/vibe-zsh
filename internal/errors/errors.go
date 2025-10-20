package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrTimeout       = errors.New("request timeout")
	ErrRateLimit     = errors.New("rate limit exceeded")
	ErrUnauthorized  = errors.New("unauthorized - check API key")
	ErrBadRequest    = errors.New("bad request")
	ErrServerError   = errors.New("server error")
	ErrNoResponse    = errors.New("no response from API")
	ErrInvalidJSON   = errors.New("invalid JSON response")
	ErrEmptyResponse = errors.New("empty response")
)

type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

func (e *APIError) IsRetryable() bool {
	return e.StatusCode == http.StatusTooManyRequests ||
		e.StatusCode >= 500
}

func NewAPIError(statusCode int, body string) error {
	err := &APIError{
		StatusCode: statusCode,
		Body:       body,
	}

	switch statusCode {
	case http.StatusBadRequest:
		err.Message = "Bad request - check your configuration"
		return fmt.Errorf("%w: %s", ErrBadRequest, err.Error())
	case http.StatusUnauthorized, http.StatusForbidden:
		err.Message = "Authentication failed"
		return fmt.Errorf("%w: %s", ErrUnauthorized, err.Error())
	case http.StatusTooManyRequests:
		err.Message = "Rate limit exceeded - please wait before retrying"
		return fmt.Errorf("%w: %s", ErrRateLimit, err.Error())
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		err.Message = "Server error - please try again"
		return fmt.Errorf("%w: %s", ErrServerError, err.Error())
	default:
		err.Message = "Unknown error"
		return err
	}
}

func IsRetryable(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsRetryable()
	}

	return errors.Is(err, ErrTimeout) ||
		errors.Is(err, ErrRateLimit) ||
		errors.Is(err, ErrServerError)
}
