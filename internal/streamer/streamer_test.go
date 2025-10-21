package streamer

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Mode != ModeWord {
		t.Errorf("DefaultConfig() Mode = %v, want %v", cfg.Mode, ModeWord)
	}
	if cfg.CharDelay != 10*time.Millisecond {
		t.Errorf("DefaultConfig() CharDelay = %v, want %v", cfg.CharDelay, 10*time.Millisecond)
	}
	if cfg.WordDelay != 30*time.Millisecond {
		t.Errorf("DefaultConfig() WordDelay = %v, want %v", cfg.WordDelay, 30*time.Millisecond)
	}
	if cfg.LineDelay != 50*time.Millisecond {
		t.Errorf("DefaultConfig() LineDelay = %v, want %v", cfg.LineDelay, 50*time.Millisecond)
	}
}

func TestStreamInstant(t *testing.T) {
	buf := &bytes.Buffer{}
	text := "Hello, World!"

	err := StreamInstant(buf, text)
	if err != nil {
		t.Fatalf("StreamInstant() error = %v", err)
	}

	if buf.String() != text {
		t.Errorf("StreamInstant() output = %q, want %q", buf.String(), text)
	}
}

func TestStreamCharacter(t *testing.T) {
	buf := &bytes.Buffer{}
	text := "Hello"

	start := time.Now()
	err := StreamCharacter(buf, text, 5*time.Millisecond)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("StreamCharacter() error = %v", err)
	}

	if buf.String() != text {
		t.Errorf("StreamCharacter() output = %q, want %q", buf.String(), text)
	}

	// Should take at least 4 delays (5 chars = 4 delays between them)
	minDuration := 4 * 5 * time.Millisecond
	if duration < minDuration {
		t.Errorf("StreamCharacter() took %v, expected at least %v", duration, minDuration)
	}
}

func TestStreamWord(t *testing.T) {
	buf := &bytes.Buffer{}
	text := "Hello World Test"

	start := time.Now()
	err := StreamWord(buf, text, 10*time.Millisecond)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("StreamWord() error = %v", err)
	}

	if buf.String() != text {
		t.Errorf("StreamWord() output = %q, want %q", buf.String(), text)
	}

	// Should take at least 2 delays (3 words = 2 delays between them)
	minDuration := 2 * 10 * time.Millisecond
	if duration < minDuration {
		t.Errorf("StreamWord() took %v, expected at least %v", duration, minDuration)
	}
}

func TestStreamLine(t *testing.T) {
	buf := &bytes.Buffer{}
	text := "Line 1\nLine 2\nLine 3"

	start := time.Now()
	err := StreamLine(buf, text, 10*time.Millisecond)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("StreamLine() error = %v", err)
	}

	if buf.String() != text {
		t.Errorf("StreamLine() output = %q, want %q", buf.String(), text)
	}

	// Should take at least 2 delays (3 lines = 2 delays between them)
	minDuration := 2 * 10 * time.Millisecond
	if duration < minDuration {
		t.Errorf("StreamLine() took %v, expected at least %v", duration, minDuration)
	}
}

func TestStreamOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	text := "Test output"

	err := StreamOutput(buf, text, 5*time.Millisecond)
	if err != nil {
		t.Fatalf("StreamOutput() error = %v", err)
	}

	if buf.String() != text {
		t.Errorf("StreamOutput() output = %q, want %q", buf.String(), text)
	}
}

func TestStreamerModes(t *testing.T) {
	tests := []struct {
		name string
		mode StreamMode
		text string
		want string
	}{
		{"instant mode", ModeInstant, "Hello World", "Hello World"},
		{"character mode", ModeCharacter, "Hi", "Hi"},
		{"word mode", ModeWord, "Hello World", "Hello World"},
		{"line mode", ModeLine, "Line1\nLine2", "Line1\nLine2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			config := &Config{
				Mode:      tt.mode,
				CharDelay: 1 * time.Millisecond,
				WordDelay: 1 * time.Millisecond,
				LineDelay: 1 * time.Millisecond,
			}
			s := New(config)

			err := s.Stream(buf, tt.text)
			if err != nil {
				t.Fatalf("Stream() error = %v", err)
			}

			if buf.String() != tt.want {
				t.Errorf("Stream() output = %q, want %q", buf.String(), tt.want)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple words",
			input: "hello world",
			want:  []string{"hello", " ", "world"},
		},
		{
			name:  "multiple spaces",
			input: "hello  world",
			want:  []string{"hello", "  ", "world"},
		},
		{
			name:  "leading space",
			input: " hello",
			want:  []string{" ", "hello"},
		},
		{
			name:  "trailing space",
			input: "hello ",
			want:  []string{"hello", " "},
		},
		{
			name:  "tabs and newlines",
			input: "hello\tworld\n",
			want:  []string{"hello", "\t", "world", "\n"},
		},
		{
			name:  "single word",
			input: "hello",
			want:  []string{"hello"},
		},
		{
			name:  "empty string",
			input: "",
			want:  []string{},
		},
		{
			name:  "only spaces",
			input: "   ",
			want:  []string{"   "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenize(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("tokenize() length = %v, want %v", len(got), len(tt.want))
				t.Errorf("got: %#v", got)
				t.Errorf("want: %#v", tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("tokenize()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestIsWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"single space", " ", true},
		{"multiple spaces", "   ", true},
		{"tab", "\t", true},
		{"newline", "\n", true},
		{"mixed whitespace", " \t\n", true},
		{"word", "hello", false},
		{"word with space", "hello ", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWhitespace(tt.input)
			if got != tt.want {
				t.Errorf("isWhitespace(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStreamerWithNilWriter(t *testing.T) {
	s := New(nil)
	err := s.Stream(nil, "test")

	if err == nil {
		t.Error("Stream() with nil writer should return error")
	}
}

func TestStreamerPreservesFormatting(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"command with flags", "docker ps -a --format json"},
		{"multiline", "line1\nline2\nline3"},
		{"mixed whitespace", "hello\tworld  test\n"},
		{"special chars", "echo 'hello world' | grep test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			config := &Config{
				Mode:      ModeWord,
				WordDelay: 0, // No delay for faster test
			}
			s := New(config)

			err := s.Stream(buf, tt.text)
			if err != nil {
				t.Fatalf("Stream() error = %v", err)
			}

			if buf.String() != tt.text {
				t.Errorf("Stream() output = %q, want %q", buf.String(), tt.text)
			}
		})
	}
}

func TestStreamerZeroDelay(t *testing.T) {
	buf := &bytes.Buffer{}
	text := "Hello World"

	config := &Config{
		Mode:      ModeCharacter,
		CharDelay: 0, // Zero delay should still work
	}
	s := New(config)

	start := time.Now()
	err := s.Stream(buf, text)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	if buf.String() != text {
		t.Errorf("Stream() output = %q, want %q", buf.String(), text)
	}

	// Should be very fast with zero delay
	if duration > 10*time.Millisecond {
		t.Errorf("Stream() with zero delay took %v, expected < 10ms", duration)
	}
}

func TestNewWithWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	s := NewWithWriter(config, buf)

	text := "Test"
	err := s.Stream(nil, text) // Pass nil, should use stored writer

	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	if buf.String() != text {
		t.Errorf("Stream() output = %q, want %q", buf.String(), text)
	}
}

func TestStreamLinePreservesEmptyLines(t *testing.T) {
	buf := &bytes.Buffer{}
	text := "line1\n\nline3"

	err := StreamLine(buf, text, 0)
	if err != nil {
		t.Fatalf("StreamLine() error = %v", err)
	}

	if buf.String() != text {
		t.Errorf("StreamLine() output = %q, want %q", buf.String(), text)
	}

	lines := strings.Split(buf.String(), "\n")
	if len(lines) != 3 {
		t.Errorf("StreamLine() produced %d lines, want 3", len(lines))
	}
	if lines[1] != "" {
		t.Errorf("StreamLine() line[1] = %q, want empty string", lines[1])
	}
}
