package config

import (
	"os"
	"strconv"
	"time"
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
}

func Load() *Config {
	return &Config{
		APIURL:               getEnv("VIBE_API_URL", "http://localhost:11434/v1"),
		APIKey:               getEnv("VIBE_API_KEY", ""),
		Model:                getEnv("VIBE_MODEL", "llama3:8b"),
		Temperature:          getEnvFloat("VIBE_TEMPERATURE", 0.2),
		MaxTokens:            getEnvInt("VIBE_MAX_TOKENS", 500),
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
	}
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
