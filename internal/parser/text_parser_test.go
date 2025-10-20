package parser

import (
	"testing"
)

func TestParseTextResponse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCommand string
		wantExplain int
	}{
		{
			name:        "simple command",
			input:       "ls -la\n# ls: list files\n# -la: long format with hidden",
			wantCommand: "ls -la",
			wantExplain: 2,
		},
		{
			name:        "markdown code block",
			input:       "```bash\nls -la\n```\n# Explanation\n- Lists files",
			wantCommand: "ls -la",
			wantExplain: 1,
		},
		{
			name:        "command only",
			input:       "docker ps",
			wantCommand: "docker ps",
			wantExplain: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTextResponse(tt.input)
			if err != nil {
				t.Fatalf("ParseTextResponse() error = %v", err)
			}
			if got.Command != tt.wantCommand {
				t.Errorf("Command = %v, want %v", got.Command, tt.wantCommand)
			}
			if len(got.Explanation) < tt.wantExplain {
				t.Errorf("Explanation count = %v, want at least %v", len(got.Explanation), tt.wantExplain)
			}
		})
	}
}

func TestCleanMarkdown(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "remove code blocks",
			input: "```bash\nls -la\n```",
			want:  "ls -la\n",
		},
		{
			name:  "remove inline code",
			input: "Use `ls -la` to list files",
			want:  "Use ls -la to list files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanMarkdown(tt.input)
			if got != tt.want {
				t.Errorf("cleanMarkdown() = %q, want %q", got, tt.want)
			}
		})
	}
}
