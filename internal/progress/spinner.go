package progress

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// SpinnerStyle defines the animation frames for different spinner styles
type SpinnerStyle string

const (
	StyleDots   SpinnerStyle = "dots"
	StyleLine   SpinnerStyle = "line"
	StyleCircle SpinnerStyle = "circle"
	StyleBounce SpinnerStyle = "bounce"
	StyleArrow  SpinnerStyle = "arrow"
)

var spinnerFrames = map[SpinnerStyle][]string{
	StyleDots:   {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	StyleLine:   {"-", "\\", "|", "/"},
	StyleCircle: {"◐", "◓", "◑", "◒"},
	StyleBounce: {"⠁", "⠂", "⠄", "⠂"},
	StyleArrow:  {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
}

// Spinner provides an animated progress indicator that writes to stderr
type Spinner struct {
	frames   []string
	message  string
	active   bool
	updateCh chan string
	doneCh   chan struct{}
	mu       sync.Mutex
	stderr   io.Writer
	ctx      context.Context
}

// NewSpinner creates a new spinner with the specified style
// If style is empty or invalid, defaults to StyleDots
func NewSpinner(style SpinnerStyle) *Spinner {
	frames, ok := spinnerFrames[style]
	if !ok || style == "" {
		frames = spinnerFrames[StyleDots]
	}

	return &Spinner{
		frames:   frames,
		updateCh: make(chan string, 10), // Buffered to prevent blocking
		doneCh:   make(chan struct{}),
		stderr:   os.Stderr,
		active:   false,
	}
}

// NewSpinnerWithWriter creates a spinner that writes to a custom writer
// Useful for testing
func NewSpinnerWithWriter(style SpinnerStyle, w io.Writer) *Spinner {
	s := NewSpinner(style)
	s.stderr = w
	return s
}

// Start begins the spinner animation with an initial message
// The spinner runs in a separate goroutine and respects context cancellation
func (s *Spinner) Start(ctx context.Context, initialMsg string) {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return // Already running
	}
	s.message = initialMsg
	s.active = true
	s.ctx = ctx
	s.mu.Unlock()

	go s.run()
}

// Update changes the spinner's message without stopping the animation
// This is thread-safe and non-blocking
func (s *Spinner) Update(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	// Non-blocking send
	select {
	case s.updateCh <- msg:
	default:
		// Channel full, update directly (rare case)
		s.message = msg
	}
}

// Stop terminates the spinner animation and clears the line
// This is safe to call multiple times
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	// Signal the goroutine to stop
	close(s.doneCh)
}

// run is the main animation loop that runs in a goroutine
// It handles cursor hiding/showing and frame animation
func (s *Spinner) run() {
	// Ensure cursor is always restored and line is cleared, even on panic
	defer func() {
		if r := recover(); r != nil {
			// Recover from panic, still cleanup
			fmt.Fprint(s.stderr, "\033[?25h\033[2K\r")
		}
	}()
	defer fmt.Fprint(s.stderr, "\033[?25h\033[2K\r") // Show cursor, clear line

	// Hide cursor for clean animation
	fmt.Fprint(s.stderr, "\033[?25l")

	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	frameIdx := 0

	for {
		select {
		case <-s.ctx.Done():
			// Context cancelled (e.g., Ctrl+C)
			return

		case <-s.doneCh:
			// Normal stop
			return

		case msg := <-s.updateCh:
			// Update message
			s.mu.Lock()
			s.message = msg
			s.mu.Unlock()

		case <-ticker.C:
			// Animate frame
			s.mu.Lock()
			frame := s.frames[frameIdx]
			msg := s.message
			s.mu.Unlock()

			// Clear line, write spinner + message
			// \033[2K clears entire line
			// \r returns cursor to start of line
			fmt.Fprintf(s.stderr, "\033[2K\r%s %s", frame, msg)

			frameIdx = (frameIdx + 1) % len(s.frames)
		}
	}
}

// IsActive returns whether the spinner is currently running
func (s *Spinner) IsActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}
