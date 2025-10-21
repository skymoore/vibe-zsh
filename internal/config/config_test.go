package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	os.Setenv("VIBE_API_URL", "http://test:1234/v1")
	os.Setenv("VIBE_MODEL", "test-model")
	os.Setenv("VIBE_TEMPERATURE", "0.5")
	os.Setenv("VIBE_MAX_TOKENS", "1000")
	os.Setenv("VIBE_TIMEOUT", "60s")
	os.Setenv("VIBE_SHOW_EXPLANATION", "false")
	defer func() {
		os.Unsetenv("VIBE_API_URL")
		os.Unsetenv("VIBE_MODEL")
		os.Unsetenv("VIBE_TEMPERATURE")
		os.Unsetenv("VIBE_MAX_TOKENS")
		os.Unsetenv("VIBE_TIMEOUT")
		os.Unsetenv("VIBE_SHOW_EXPLANATION")
	}()

	cfg := Load()

	if cfg.APIURL != "http://test:1234/v1" {
		t.Errorf("APIURL = %v, want http://test:1234/v1", cfg.APIURL)
	}
	if cfg.Model != "test-model" {
		t.Errorf("Model = %v, want test-model", cfg.Model)
	}
	if cfg.Temperature != 0.5 {
		t.Errorf("Temperature = %v, want 0.5", cfg.Temperature)
	}
	if cfg.MaxTokens != 1000 {
		t.Errorf("MaxTokens = %v, want 1000", cfg.MaxTokens)
	}
	if cfg.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want 60s", cfg.Timeout)
	}
	if cfg.ShowExplanation != false {
		t.Errorf("ShowExplanation = %v, want false", cfg.ShowExplanation)
	}
}

func TestLoadDefaults(t *testing.T) {
	cfg := Load()

	if cfg.APIURL != "http://localhost:11434/v1" {
		t.Errorf("Default APIURL = %v, want http://localhost:11434/v1", cfg.APIURL)
	}
	if cfg.Model != "llama3:8b" {
		t.Errorf("Default Model = %v, want llama3:8b", cfg.Model)
	}
	if cfg.Temperature != 0.2 {
		t.Errorf("Default Temperature = %v, want 0.2", cfg.Temperature)
	}
	if cfg.ShowExplanation != true {
		t.Errorf("Default ShowExplanation = %v, want true", cfg.ShowExplanation)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("Default MaxRetries = %v, want 3", cfg.MaxRetries)
	}
	if cfg.EnableJSONExtraction != true {
		t.Errorf("Default EnableJSONExtraction = %v, want true", cfg.EnableJSONExtraction)
	}
	if cfg.StrictValidation != true {
		t.Errorf("Default StrictValidation = %v, want true", cfg.StrictValidation)
	}
}
