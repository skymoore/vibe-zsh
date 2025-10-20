package formatter

import (
	"testing"

	"github.com/skymoore/vibe-zsh/internal/schema"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name            string
		resp            *schema.CommandResponse
		showExplanation bool
		showWarnings    bool
		want            string
	}{
		{
			name: "with explanation",
			resp: &schema.CommandResponse{
				Command:     "ls -la",
				Explanation: []string{"ls: list files", "-la: long format"},
			},
			showExplanation: true,
			showWarnings:    true,
			want:            "ls -la\n# ls: list files\n# -la: long format",
		},
		{
			name: "without explanation",
			resp: &schema.CommandResponse{
				Command:     "ls -la",
				Explanation: []string{"ls: list files", "-la: long format"},
			},
			showExplanation: false,
			showWarnings:    true,
			want:            "ls -la",
		},
		{
			name: "with warning",
			resp: &schema.CommandResponse{
				Command:     "rm -rf /tmp",
				Explanation: []string{"Removes files"},
				Warning:     "Dangerous command",
			},
			showExplanation: true,
			showWarnings:    true,
			want:            "rm -rf /tmp\n# Removes files\n# WARNING: Dangerous command",
		},
		{
			name: "warning disabled",
			resp: &schema.CommandResponse{
				Command:     "rm -rf /tmp",
				Explanation: []string{"Removes files"},
				Warning:     "Dangerous command",
			},
			showExplanation: true,
			showWarnings:    false,
			want:            "rm -rf /tmp\n# Removes files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.resp, tt.showExplanation, tt.showWarnings)
			if got != tt.want {
				t.Errorf("Format() = %q, want %q", got, tt.want)
			}
		})
	}
}
