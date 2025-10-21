package progress

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewSpinner(t *testing.T) {
	tests := []struct {
		name  string
		style SpinnerStyle
		want  int // expected number of frames
	}{
		{"dots style", StyleDots, 10},
		{"line style", StyleLine, 4},
		{"circle style", StyleCircle, 4},
		{"bounce style", StyleBounce, 4},
		{"arrow style", StyleArrow, 8},
		{"empty style defaults to dots", "", 10},
		{"invalid style defaults to dots", "invalid", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSpinner(tt.style)
			if len(s.frames) != tt.want {
				t.Errorf("NewSpinner() frames = %v, want %v", len(s.frames), tt.want)
			}
			if s.active {
				t.Error("NewSpinner() should not be active initially")
			}
		})
	}
}

func TestSpinnerStartStop(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWithWriter(StyleDots, buf)
	ctx := context.Background()

	if s.IsActive() {
		t.Error("Spinner should not be active before Start()")
	}

	s.Start(ctx, "Testing...")
	time.Sleep(50 * time.Millisecond) // Let it start

	if !s.IsActive() {
		t.Error("Spinner should be active after Start()")
	}

	s.Stop()
	time.Sleep(50 * time.Millisecond) // Let it stop

	if s.IsActive() {
		t.Error("Spinner should not be active after Stop()")
	}

	output := buf.String()
	if !strings.Contains(output, "Testing...") {
		t.Errorf("Output should contain initial message, got: %s", output)
	}

	// Check for cursor hide/show sequences
	if !strings.Contains(output, "\033[?25l") {
		t.Error("Output should contain cursor hide sequence")
	}
	if !strings.Contains(output, "\033[?25h") {
		t.Error("Output should contain cursor show sequence")
	}
}

func TestSpinnerUpdate(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWithWriter(StyleLine, buf)
	ctx := context.Background()

	s.Start(ctx, "Initial")
	time.Sleep(50 * time.Millisecond)

	s.Update("Updated")
	time.Sleep(150 * time.Millisecond) // Wait for at least one frame with new message

	s.Stop()
	time.Sleep(50 * time.Millisecond)

	output := buf.String()
	if !strings.Contains(output, "Initial") {
		t.Error("Output should contain initial message")
	}
	if !strings.Contains(output, "Updated") {
		t.Error("Output should contain updated message")
	}
}

func TestSpinnerContextCancellation(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWithWriter(StyleDots, buf)
	ctx, cancel := context.WithCancel(context.Background())

	s.Start(ctx, "Testing context...")
	time.Sleep(50 * time.Millisecond)

	if !s.IsActive() {
		t.Error("Spinner should be active")
	}

	// Cancel context
	cancel()
	time.Sleep(100 * time.Millisecond) // Wait for goroutine to exit

	// Spinner should still report active until Stop() is called
	// but the goroutine should have exited
	output := buf.String()

	// Should have cleanup sequences
	if !strings.Contains(output, "\033[?25h") {
		t.Error("Output should contain cursor show sequence after context cancellation")
	}
}

func TestSpinnerMultipleStops(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWithWriter(StyleDots, buf)
	ctx := context.Background()

	s.Start(ctx, "Testing...")
	time.Sleep(50 * time.Millisecond)

	// Multiple stops should not panic
	s.Stop()
	s.Stop()
	s.Stop()

	time.Sleep(50 * time.Millisecond)
}

func TestSpinnerUpdateBeforeStart(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWithWriter(StyleDots, buf)

	// Update before start should not panic
	s.Update("Should not crash")

	if s.IsActive() {
		t.Error("Spinner should not be active")
	}
}

func TestSpinnerAnimation(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWithWriter(StyleLine, buf)
	ctx := context.Background()

	s.Start(ctx, "Animating")
	time.Sleep(400 * time.Millisecond) // Wait for multiple frames (80ms per frame)

	s.Stop()
	time.Sleep(50 * time.Millisecond)

	output := buf.String()

	// Should contain multiple different frames from StyleLine: -, \, |, /
	frames := []string{"-", "\\", "|", "/"}
	foundFrames := 0
	for _, frame := range frames {
		if strings.Contains(output, frame+" Animating") {
			foundFrames++
		}
	}

	if foundFrames < 2 {
		t.Errorf("Expected to see at least 2 different frames, found %d", foundFrames)
	}
}

func TestSpinnerClearLine(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWithWriter(StyleDots, buf)
	ctx := context.Background()

	s.Start(ctx, "Testing clear")
	time.Sleep(100 * time.Millisecond)
	s.Stop()
	time.Sleep(50 * time.Millisecond)

	output := buf.String()

	// Should contain line clear sequence
	if !strings.Contains(output, "\033[2K") {
		t.Error("Output should contain line clear sequence")
	}

	// Should contain carriage return
	if !strings.Contains(output, "\r") {
		t.Error("Output should contain carriage return")
	}
}

func TestSpinnerDoubleStart(t *testing.T) {
	buf := &bytes.Buffer{}
	s := NewSpinnerWithWriter(StyleDots, buf)
	ctx := context.Background()

	s.Start(ctx, "First start")
	time.Sleep(50 * time.Millisecond)

	// Second start should be ignored
	s.Start(ctx, "Second start")
	time.Sleep(50 * time.Millisecond)

	s.Stop()
	time.Sleep(50 * time.Millisecond)

	output := buf.String()

	// Should only contain first message
	if !strings.Contains(output, "First start") {
		t.Error("Output should contain first start message")
	}
}
