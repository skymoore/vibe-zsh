package streamer

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode"
)

// StreamMode defines how content should be streamed
type StreamMode string

const (
	// ModeCharacter streams one character at a time
	ModeCharacter StreamMode = "character"
	// ModeWord streams one word at a time
	ModeWord StreamMode = "word"
	// ModeLine streams one line at a time
	ModeLine StreamMode = "line"
	// ModeInstant outputs everything immediately (no streaming)
	ModeInstant StreamMode = "instant"
)

// Config holds configuration for the streamer
type Config struct {
	// Mode determines how content is streamed
	Mode StreamMode
	// CharDelay is the delay between characters (used in ModeCharacter)
	CharDelay time.Duration
	// WordDelay is the delay between words (used in ModeWord)
	WordDelay time.Duration
	// LineDelay is the delay between lines (used in ModeLine)
	LineDelay time.Duration
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		Mode:      ModeWord,
		CharDelay: 10 * time.Millisecond,
		WordDelay: 30 * time.Millisecond,
		LineDelay: 50 * time.Millisecond,
	}
}

// Streamer handles streaming output to a writer
type Streamer struct {
	config *Config
	writer io.Writer
}

// New creates a new Streamer with the given configuration
func New(config *Config) *Streamer {
	if config == nil {
		config = DefaultConfig()
	}
	return &Streamer{
		config: config,
		writer: nil,
	}
}

// NewWithWriter creates a new Streamer with a specific writer
func NewWithWriter(config *Config, w io.Writer) *Streamer {
	s := New(config)
	s.writer = w
	return s
}

// Stream outputs text to the writer according to the configured mode
func (s *Streamer) Stream(w io.Writer, text string) error {
	if w == nil {
		w = s.writer
	}
	if w == nil {
		return fmt.Errorf("no writer specified")
	}

	switch s.config.Mode {
	case ModeInstant:
		return s.streamInstant(w, text)
	case ModeCharacter:
		return s.streamCharacter(w, text)
	case ModeWord:
		return s.streamWord(w, text)
	case ModeLine:
		return s.streamLine(w, text)
	default:
		return s.streamWord(w, text) // Default to word mode
	}
}

// streamInstant outputs everything immediately
func (s *Streamer) streamInstant(w io.Writer, text string) error {
	_, err := fmt.Fprint(w, text)
	return err
}

// streamCharacter outputs one character at a time
func (s *Streamer) streamCharacter(w io.Writer, text string) error {
	for _, ch := range text {
		if _, err := fmt.Fprintf(w, "%c", ch); err != nil {
			return err
		}
		if s.config.CharDelay > 0 {
			time.Sleep(s.config.CharDelay)
		}
	}
	return nil
}

// streamWord outputs one word at a time, preserving whitespace
func (s *Streamer) streamWord(w io.Writer, text string) error {
	// Split into tokens (words and whitespace)
	tokens := tokenize(text)

	for i, token := range tokens {
		if _, err := fmt.Fprint(w, token); err != nil {
			return err
		}

		// Add delay after words (not after whitespace)
		if i < len(tokens)-1 && !isWhitespace(token) && s.config.WordDelay > 0 {
			time.Sleep(s.config.WordDelay)
		}
	}
	return nil
}

// streamLine outputs one line at a time
func (s *Streamer) streamLine(w io.Writer, text string) error {
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if _, err := fmt.Fprint(w, line); err != nil {
			return err
		}

		// Add newline back (except for last line if it didn't have one)
		if i < len(lines)-1 {
			if _, err := fmt.Fprint(w, "\n"); err != nil {
				return err
			}
			if s.config.LineDelay > 0 {
				time.Sleep(s.config.LineDelay)
			}
		}
	}
	return nil
}

// StreamOutput is a convenience function that streams text with default word mode
func StreamOutput(w io.Writer, text string, charDelay time.Duration) error {
	config := &Config{
		Mode:      ModeWord,
		CharDelay: charDelay,
		WordDelay: charDelay * 3, // Words slightly slower than chars
	}
	s := New(config)
	return s.Stream(w, text)
}

// StreamInstant is a convenience function for instant output
func StreamInstant(w io.Writer, text string) error {
	config := &Config{
		Mode: ModeInstant,
	}
	s := New(config)
	return s.Stream(w, text)
}

// StreamCharacter is a convenience function for character-by-character streaming
func StreamCharacter(w io.Writer, text string, delay time.Duration) error {
	config := &Config{
		Mode:      ModeCharacter,
		CharDelay: delay,
	}
	s := New(config)
	return s.Stream(w, text)
}

// StreamWord is a convenience function for word-by-word streaming
func StreamWord(w io.Writer, text string, delay time.Duration) error {
	config := &Config{
		Mode:      ModeWord,
		WordDelay: delay,
	}
	s := New(config)
	return s.Stream(w, text)
}

// StreamLine is a convenience function for line-by-line streaming
func StreamLine(w io.Writer, text string, delay time.Duration) error {
	config := &Config{
		Mode:      ModeLine,
		LineDelay: delay,
	}
	s := New(config)
	return s.Stream(w, text)
}

// tokenize splits text into words and whitespace tokens
func tokenize(text string) []string {
	if text == "" {
		return []string{}
	}

	var tokens []string
	var current strings.Builder

	// Determine if we're starting in a word or whitespace
	firstChar := rune(text[0])
	inWord := !unicode.IsSpace(firstChar)

	for _, ch := range text {
		isSpace := unicode.IsSpace(ch)

		// Check if we're transitioning between word and whitespace
		if isSpace == inWord {
			// Transition: save current token and start new one
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			inWord = !isSpace
		}

		current.WriteRune(ch)
	}

	// Add final token
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// isWhitespace checks if a string contains only whitespace
func isWhitespace(s string) bool {
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			return false
		}
	}
	return len(s) > 0
}
