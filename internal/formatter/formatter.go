package formatter

import (
	"strings"

	"github.com/skymoore/vibe-zsh/internal/schema"
)

func Format(resp *schema.CommandResponse, showExplanation bool) string {
	var buf strings.Builder

	buf.WriteString(resp.Command)

	if showExplanation {
		buf.WriteString("\n")

		for _, line := range resp.Explanation {
			buf.WriteString("# ")
			buf.WriteString(line)
			buf.WriteString("\n")
		}

		if resp.Warning != "" {
			buf.WriteString("# WARNING: ")
			buf.WriteString(resp.Warning)
			buf.WriteString("\n")
		}

		return strings.TrimSuffix(buf.String(), "\n")
	}

	return buf.String()
}
