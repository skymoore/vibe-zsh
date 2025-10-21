package config

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/skymoore/vibe-zsh/internal/progress"
)

type Config struct {
	APIURL               string
	APIKey               string
	Model                string
	Temperature          float64
	MaxTokens            int
	Timeout              time.Duration
	UseStructuredOutput  bool
	ShowExplanation      bool
	EnableCache          bool
	CacheDir             string
	CacheTTL             time.Duration
	InteractiveMode      bool
	ShowWarnings         bool
	MaxRetries           int
	EnableJSONExtraction bool
	StrictValidation     bool
	EnableDebugLogs      bool
	ShowRetryStatus      bool
	ShowProgress         bool
	ProgressStyle        progress.SpinnerStyle
	StreamOutput         bool
	StreamDelay          time.Duration
	OSName               string
	Shell                string
	EnableHistory        bool
	HistorySize          int
	HistoryKey           string
}

func Load() *Config {
	return &Config{
		APIURL:               getEnv("VIBE_API_URL", "http://localhost:11434/v1"),
		APIKey:               getEnv("VIBE_API_KEY", ""),
		Model:                getEnv("VIBE_MODEL", "llama3:8b"),
		Temperature:          getEnvFloat("VIBE_TEMPERATURE", 0.2),
		MaxTokens:            getEnvInt("VIBE_MAX_TOKENS", 1000),
		Timeout:              getEnvDuration("VIBE_TIMEOUT", 30*time.Second),
		UseStructuredOutput:  getEnvBool("VIBE_USE_STRUCTURED_OUTPUT", true),
		ShowExplanation:      getEnvBool("VIBE_SHOW_EXPLANATION", true),
		EnableCache:          getEnvBool("VIBE_ENABLE_CACHE", true),
		CacheDir:             getEnv("VIBE_CACHE_DIR", ""),
		CacheTTL:             getEnvDuration("VIBE_CACHE_TTL", 24*time.Hour),
		InteractiveMode:      getEnvBool("VIBE_INTERACTIVE", false),
		ShowWarnings:         getEnvBool("VIBE_SHOW_WARNINGS", true),
		MaxRetries:           getEnvInt("VIBE_MAX_RETRIES", 3),
		EnableJSONExtraction: getEnvBool("VIBE_ENABLE_JSON_EXTRACTION", true),
		StrictValidation:     getEnvBool("VIBE_STRICT_VALIDATION", true),
		EnableDebugLogs:      getEnvBool("VIBE_DEBUG_LOGS", false),
		ShowRetryStatus:      getEnvBool("VIBE_SHOW_RETRY_STATUS", true),
		ShowProgress:         getEnvBool("VIBE_SHOW_PROGRESS", true),
		ProgressStyle:        getEnvProgressStyle("VIBE_PROGRESS_STYLE", progress.StyleDots),
		StreamOutput:         getEnvBool("VIBE_STREAM_OUTPUT", true),
		StreamDelay:          getEnvDuration("VIBE_STREAM_DELAY", 20*time.Millisecond),
		OSName:               getOSName(),
		Shell:                getShell(),
		EnableHistory:        getEnvBool("VIBE_ENABLE_HISTORY", true),
		HistorySize:          getEnvInt("VIBE_HISTORY_SIZE", 100),
		HistoryKey:           getEnv("VIBE_HISTORY_KEY", "^Xh"),
	}
}

func getOSName() string {
	osName := runtime.GOOS
	// Normalize OS names for better LLM understanding
	switch osName {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return osName
	}
}

func getShell() string {
	// Try to detect shell from environment
	shell := os.Getenv("SHELL")
	if shell == "" {
		// Default to zsh since this is a zsh plugin
		return "zsh"
	}
	// Extract just the shell name (e.g., /bin/zsh -> zsh)
	if idx := len(shell) - 1; idx >= 0 {
		for i := idx; i >= 0; i-- {
			if shell[i] == '/' {
				return shell[i+1:]
			}
		}
	}
	return shell
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

func getEnvProgressStyle(key string, defaultValue progress.SpinnerStyle) progress.SpinnerStyle {
	if value := os.Getenv(key); value != "" {
		style := progress.SpinnerStyle(value)
		// Validate it's a known style
		switch style {
		case progress.StyleDots, progress.StyleLine, progress.StyleCircle, progress.StyleBounce, progress.StyleArrow:
			return style
		}
	}
	return defaultValue
}
