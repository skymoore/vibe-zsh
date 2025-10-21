package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/skymoore/vibe-zsh/internal/logger"
	"github.com/skymoore/vibe-zsh/internal/schema"
)

func ParseTextResponse(text string) (*schema.CommandResponse, error) {
	text = strings.TrimSpace(text)

	text = cleanMarkdown(text)

	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil, nil
	}

	var command string
	var explanation []string
	var warning string

	inCodeBlock := false
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if command == "" && line != "" && !inCodeBlock {
			if !strings.Contains(line, ":") && !strings.HasPrefix(line, "#") &&
				!strings.HasPrefix(line, "*") && !strings.HasPrefix(line, "-") {
				command = line
				continue
			}
		}

		if inCodeBlock && command == "" && line != "" {
			command = line
			continue
		}

		if command != "" && line != "" {
			if strings.HasPrefix(strings.ToLower(line), "warning:") {
				warning = strings.TrimPrefix(strings.TrimPrefix(line, "warning:"), "WARNING:")
				warning = strings.TrimSpace(warning)
			} else if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") ||
				strings.HasPrefix(line, "*") || strings.Contains(line, ":") {
				cleaned := strings.TrimPrefix(line, "#")
				cleaned = strings.TrimPrefix(cleaned, "-")
				cleaned = strings.TrimPrefix(cleaned, "*")
				cleaned = strings.TrimSpace(cleaned)
				if cleaned != "" && i > 0 {
					explanation = append(explanation, cleaned)
				}
			}
		}
	}

	if command == "" {
		command = strings.Split(text, "\n")[0]
	}

	if len(explanation) == 0 {
		explanation = []string{"Generated command"}
	}

	return &schema.CommandResponse{
		Command:     command,
		Explanation: explanation,
		Warning:     warning,
	}, nil
}

func cleanMarkdown(text string) string {
	text = regexp.MustCompile("```(?:bash|sh|shell)?\n?").ReplaceAllString(text, "")
	text = regexp.MustCompile("```\n?").ReplaceAllString(text, "")

	text = regexp.MustCompile("`([^`]+)`").ReplaceAllString(text, "$1")

	text = regexp.MustCompile("(?i)^\\*\\*Command\\*\\*:?\\s*").ReplaceAllString(text, "")
	text = regexp.MustCompile("(?i)^\\*\\*Explanation\\*\\*:?\\s*").ReplaceAllString(text, "")

	return text
}

func ExtractJSON(corruptedText string) (string, error) {
	original := corruptedText

	firstBrace := strings.Index(corruptedText, "{")
	lastBrace := strings.LastIndex(corruptedText, "}")

	if firstBrace >= 0 && lastBrace > firstBrace {
		candidate := corruptedText[firstBrace : lastBrace+1]
		if json.Valid([]byte(candidate)) {
			trimmedPrefix := corruptedText[:firstBrace]
			trimmedSuffix := corruptedText[lastBrace+1:]
			logger.LogJSONExtraction(original, candidate, trimmedPrefix, trimmedSuffix)
			return candidate, nil
		}
	}

	cleaned := RemoveGarbagePatterns(corruptedText)
	if json.Valid([]byte(cleaned)) {
		logger.Debug("JSON extraction: RemoveGarbagePatterns succeeded")
		return cleaned, nil
	}

	fixed := AttemptJSONRepair(cleaned)
	if json.Valid([]byte(fixed)) {
		logger.Debug("JSON extraction: AttemptJSONRepair succeeded")
		return fixed, nil
	}

	return "", fmt.Errorf("unable to extract valid JSON from response")
}

func RemoveGarbagePatterns(text string) string {
	re1 := regexp.MustCompile(`\x1b\[[0-9]+~`)
	text = re1.ReplaceAllString(text, "")

	text = strings.ReplaceAll(text, "\u2028", "")
	text = strings.ReplaceAll(text, "\u2029", "")

	re2 := regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`)
	text = re2.ReplaceAllString(text, "")

	re3 := regexp.MustCompile(`[.…]{3,}`)
	text = re3.ReplaceAllString(text, "")

	re4 := regexp.MustCompile(`[^\x20-\x7E\n\t]`)
	text = re4.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

func AttemptJSONRepair(text string) string {
	text = strings.TrimSpace(text)

	if !strings.HasPrefix(text, "{") {
		if idx := strings.Index(text, "{"); idx >= 0 {
			text = text[idx:]
		}
	}

	if !strings.HasSuffix(text, "}") {
		if idx := strings.LastIndex(text, "}"); idx >= 0 {
			text = text[:idx+1]
		}
	}

	text = regexp.MustCompile(`,\s*}`).ReplaceAllString(text, "}")
	text = regexp.MustCompile(`,\s*]`).ReplaceAllString(text, "]")

	return text
}
