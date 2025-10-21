# Progress Package

Provides animated progress indicators (spinners) for terminal applications.

## Features

- **Multiple spinner styles**: dots, line, circle, bounce, arrow
- **Thread-safe**: Safe to update from multiple goroutines
- **Context-aware**: Respects context cancellation (Ctrl+C)
- **TTY detection**: Auto-detects if output is a terminal
- **Clean cleanup**: Always restores cursor state, even on panic
- **Non-blocking updates**: Buffered channels prevent blocking

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/skymoore/vibe-zsh/internal/progress"
)

func main() {
    ctx := context.Background()
    
    // Create spinner with dots style
    spinner := progress.NewSpinner(progress.StyleDots)
    
    // Always ensure cleanup
    defer spinner.Stop()
    
    // Start animation
    spinner.Start(ctx, "Loading...")
    
    // Do work...
    time.Sleep(2 * time.Second)
    
    // Update message
    spinner.Update("Processing...")
    
    // Do more work...
    time.Sleep(2 * time.Second)
    
    // Stop is called by defer
}
```

### With Context Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Handle Ctrl+C
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, os.Interrupt)
go func() {
    <-sigCh
    cancel() // Spinner will stop automatically
}()

spinner := progress.NewSpinner(progress.StyleDots)
defer spinner.Stop()

spinner.Start(ctx, "Working...")
// Spinner stops when context is cancelled
```

### TTY Detection

```go
import "github.com/skymoore/vibe-zsh/internal/progress"

if progress.IsStderrTerminal() {
    spinner := progress.NewSpinner(progress.StyleDots)
    defer spinner.Stop()
    spinner.Start(ctx, "Processing...")
    // ... do work
} else {
    // Not a terminal, skip spinner
    fmt.Fprintln(os.Stderr, "Processing...")
    // ... do work
}
```

### Available Styles

```go
progress.StyleDots   // ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ (default)
progress.StyleLine   // -\|/
progress.StyleCircle // ◐◓◑◒
progress.StyleBounce // ⠁⠂⠄⠂
progress.StyleArrow  // ←↖↑↗→↘↓↙
```

### Testing

For testing, use `NewSpinnerWithWriter` to capture output:

```go
func TestMyFunction(t *testing.T) {
    buf := &bytes.Buffer{}
    spinner := progress.NewSpinnerWithWriter(progress.StyleDots, buf)
    
    ctx := context.Background()
    spinner.Start(ctx, "Testing...")
    time.Sleep(100 * time.Millisecond)
    spinner.Stop()
    
    output := buf.String()
    if !strings.Contains(output, "Testing...") {
        t.Error("Expected spinner message")
    }
}
```

## Implementation Details

### ANSI Escape Sequences

- `\033[?25l` - Hide cursor
- `\033[?25h` - Show cursor
- `\033[2K` - Clear entire line
- `\r` - Return cursor to start of line

### Thread Safety

- All public methods are thread-safe
- Uses `sync.Mutex` for state protection
- Buffered channels for non-blocking updates

### Cleanup Guarantees

- `defer` in `run()` ensures cursor restoration
- Panic recovery prevents terminal corruption
- Context cancellation triggers cleanup
- Multiple `Stop()` calls are safe

## Performance

- Frame rate: ~12.5 FPS (80ms per frame)
- Memory: Minimal (single goroutine, small buffers)
- CPU: Negligible (sleeps between frames)

## Best Practices

1. **Always use defer**: `defer spinner.Stop()`
2. **Check TTY**: Only show spinner if stderr is a terminal
3. **Use context**: Pass context for cancellation support
4. **Update sparingly**: Don't update on every iteration of tight loops
5. **Stop before output**: Ensure spinner is stopped before printing results
