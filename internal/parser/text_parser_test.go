package parser

import (
	"encoding/json"
	"testing"

	"github.com/skymoore/vibe-zsh/internal/schema"
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

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValid   bool
		wantCommand string
	}{
		{
			name:        "clean JSON",
			input:       `{"command":"ls -la","explanation":["list files"]}`,
			wantValid:   true,
			wantCommand: "ls -la",
		},
		{
			name:        "JSON with escape sequences",
			input:       "\x1b[200~\x1b[201~{\"command\":\"aws ecr describe-repositories\",\"explanation\":[\"list repos\"]}\x1b[201~\x1b[200~",
			wantValid:   true,
			wantCommand: "aws ecr describe-repositories",
		},
		{
			name:        "JSON with Unicode pollution",
			input:       "AWS\u2028…\u2029{\"command\":\"docker ps\",\"explanation\":[\"list containers\"]}…\u2028",
			wantValid:   true,
			wantCommand: "docker ps",
		},
		{
			name:        "JSON buried in text",
			input:       "Here is the command:\n\n{\"command\":\"git status\",\"explanation\":[\"show status\"]}\n\nLet me know if you need help!",
			wantValid:   true,
			wantCommand: "git status",
		},
		{
			name:        "malformed with trailing comma",
			input:       `{"command":"ls","explanation":["test"],}`,
			wantValid:   true,
			wantCommand: "ls",
		},
		{
			name:      "completely corrupted - no JSON",
			input:     "AWS?…..??……… ...…..??……....??…... …………â[201~[200~¦………",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractJSON(tt.input)

			if tt.wantValid {
				if err != nil {
					t.Errorf("ExtractJSON() error = %v, want nil", err)
					return
				}

				if !json.Valid([]byte(got)) {
					t.Errorf("ExtractJSON() returned invalid JSON: %q", got)
					return
				}

				var cmdResp schema.CommandResponse
				if err := json.Unmarshal([]byte(got), &cmdResp); err != nil {
					t.Errorf("ExtractJSON() returned unparseable JSON: %v", err)
					return
				}

				if cmdResp.Command != tt.wantCommand {
					t.Errorf("Command = %q, want %q", cmdResp.Command, tt.wantCommand)
				}
			} else {
				if err == nil {
					t.Errorf("ExtractJSON() should have failed but got: %q", got)
				}
			}
		})
	}
}

func TestRemoveGarbagePatterns(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "escape sequences",
			input: "\x1b[200~hello\x1b[201~world",
			want:  "helloworld",
		},
		{
			name:  "unicode separators",
			input: "line1\u2028line2\u2029line3",
			want:  "line1line2line3",
		},
		{
			name:  "ellipsis",
			input: "text…..…with…….ellipsis",
			want:  "textwithellipsis",
		},
		{
			name:  "control characters",
			input: "normal\x00text\x1fhere\x7f",
			want:  "normaltexthere",
		},
		{
			name:  "combined garbage",
			input: "\x1b[201~AWS…..{\"command\":\"test\"}\u2028",
			want:  "AWS{\"command\":\"test\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveGarbagePatterns(tt.input)
			if got != tt.want {
				t.Errorf("RemoveGarbagePatterns() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAttemptJSONRepair(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{
			name:      "trailing comma in object",
			input:     `{"command":"ls","explanation":["test"],}`,
			wantValid: true,
		},
		{
			name:      "trailing comma in array",
			input:     `{"command":"ls","explanation":["test",]}`,
			wantValid: true,
		},
		{
			name:      "extra text before brace",
			input:     `some text {"command":"ls","explanation":["test"]}`,
			wantValid: true,
		},
		{
			name:      "extra text after brace",
			input:     `{"command":"ls","explanation":["test"]} extra`,
			wantValid: true,
		},
		{
			name:      "multiple issues",
			input:     `text {"command":"ls","explanation":["test"],} more`,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AttemptJSONRepair(tt.input)

			if tt.wantValid {
				if !json.Valid([]byte(got)) {
					t.Errorf("AttemptJSONRepair() returned invalid JSON: %q", got)
				}
			}
		})
	}
}
