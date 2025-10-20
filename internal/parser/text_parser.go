package parser

import (
	"regexp"
	"strings"

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
