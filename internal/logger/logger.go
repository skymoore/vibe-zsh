package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	debugEnabled bool
	logger       *log.Logger
)

func Init(enabled bool) {
	debugEnabled = enabled
	logger = log.New(os.Stderr, "[VIBE] ", log.LstdFlags)
}

func LogParsingFailure(attempt int, layer string, rawResponse string, err error) {
	if !debugEnabled {
		return
	}
	logger.Printf("[ATTEMPT %d][%s] Parsing failed: %v\nRaw response (first 500 chars):\n%s\n",
		attempt, layer, err, truncate(rawResponse, 500))
}

func LogJSONExtraction(original string, extracted string, trimmedPrefix, trimmedSuffix string) {
	if !debugEnabled {
		return
	}
	logger.Printf("[JSON_EXTRACT] Success\nTrimmed prefix: %q\nTrimmed suffix: %q\nExtracted length: %d bytes\n",
		truncate(trimmedPrefix, 100), truncate(trimmedSuffix, 100), len(extracted))
}

func LogLayerSuccess(layer string, attempt int) {
	if !debugEnabled {
		return
	}
	logger.Printf("[SUCCESS] Layer '%s' succeeded on attempt %d\n", layer, attempt)
}

func Debug(format string, args ...interface{}) {
	if !debugEnabled {
		return
	}
	logger.Printf("[DEBUG] "+format+"\n", args...)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return fmt.Sprintf("%s... (truncated, total: %d chars)", s[:maxLen], len(s))
}
