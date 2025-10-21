# Streamer Package

Provides streaming output capabilities for terminal applications, creating a typewriter effect for better user experience.

## Features

- **Multiple streaming modes**: character, word, line, instant
- **Configurable delays**: Control speed of streaming
- **TTY detection**: Auto-detects if output should be streamed
- **Format preservation**: Maintains whitespace and special characters
- **Word-aware**: Doesn't break words mid-stream in word mode
- **Zero-allocation tokenization**: Efficient text processing

## Usage

### Quick Start

```go
import (
    "os"
    "time"
    "github.com/skymoore/vibe-zsh/internal/streamer"
)

func main() {
    text := "docker ps -a --format json"
    
    // Simple word-by-word streaming
    streamer.StreamOutput(os.Stdout, text, 10*time.Millisecond)
}
```

### Convenience Functions

```go
// Instant output (no streaming)
streamer.StreamInstant(os.Stdout, "Hello World")

// Character-by-character
streamer.StreamCharacter(os.Stdout, "Hello", 10*time.Millisecond)

// Word-by-word (recommended for commands)
streamer.StreamWord(os.Stdout, "docker ps -a", 20*time.Millisecond)

// Line-by-line
streamer.StreamLine(os.Stdout, "line1\nline2\nline3", 50*time.Millisecond)
```

### Advanced Configuration

```go
config := &streamer.Config{
    Mode:      streamer.ModeWord,
    CharDelay: 10 * time.Millisecond,
    WordDelay: 30 * time.Millisecond,
    LineDelay: 50 * time.Millisecond,
}

s := streamer.New(config)
s.Stream(os.Stdout, "Your text here")
```

### TTY-Aware Streaming

```go
import "github.com/skymoore/vibe-zsh/internal/streamer"

if streamer.ShouldStream() {
    // Terminal output - use streaming
    streamer.StreamWord(os.Stdout, command, 20*time.Millisecond)
} else {
    // Piped/redirected - instant output
    fmt.Println(command)
}
```

### Streaming Modes

#### ModeInstant
Outputs everything immediately. Use when streaming is disabled or not needed.

```go
streamer.StreamInstant(os.Stdout, text)
```

#### ModeCharacter
Streams one character at a time. Creates classic typewriter effect.

```go
// Slower, more dramatic
streamer.StreamCharacter(os.Stdout, text, 20*time.Millisecond)
```

#### ModeWord (Recommended)
Streams one word at a time. Best for command output - maintains readability while providing visual feedback.

```go
// Balanced speed and readability
streamer.StreamWord(os.Stdout, "docker ps -a", 30*time.Millisecond)
```

#### ModeLine
Streams one line at a time. Good for multi-line output like explanations.

```go
explanation := "# docker: Docker CLI\n# ps: List containers\n# -a: Show all"
streamer.StreamLine(os.Stdout, explanation, 50*time.Millisecond)
```

## Integration Example

### With Progress Spinner

```go
import (
    "context"
    "os"
    "github.com/skymoore/vibe-zsh/internal/progress"
    "github.com/skymoore/vibe-zsh/internal/streamer"
)

func generateCommand(query string) {
    ctx := context.Background()
    
    // Show spinner during API call
    spinner := progress.NewSpinner(progress.StyleDots)
    defer spinner.Stop()
    spinner.Start(ctx, "Generating command...")
    
    // ... API call to generate command ...
    command := "docker ps -a"
    
    // Stop spinner before streaming output
    spinner.Stop()
    
    // Stream the command
    if streamer.ShouldStream() {
        streamer.StreamWord(os.Stdout, command, 20*time.Millisecond)
    } else {
        fmt.Print(command)
    }
}
```

## Configuration Guidelines

### Recommended Delays

- **Fast**: 5-10ms per character, 15-20ms per word
- **Medium**: 10-15ms per character, 30-40ms per word
- **Slow**: 20-30ms per character, 50-80ms per word

### Mode Selection

- **Commands**: Use `ModeWord` - maintains readability
- **Explanations**: Use `ModeLine` - easier to read line-by-line
- **Dramatic effect**: Use `ModeCharacter` - classic typewriter
- **Piped output**: Use `ModeInstant` - no delays

### Performance Considerations

```go
// For short text (< 50 chars), instant is fine
if len(text) < 50 {
    streamer.StreamInstant(os.Stdout, text)
} else {
    streamer.StreamWord(os.Stdout, text, 20*time.Millisecond)
}
```

## Implementation Details

### Word-Aware Tokenization

The streamer intelligently splits text into words and whitespace, preserving formatting:

```go
Input:  "hello  world\n"
Tokens: ["hello", "  ", "world", "\n"]
```

This ensures:
- Multiple spaces are preserved
- Tabs and newlines are maintained
- Words are never broken mid-stream

### Format Preservation

All streaming modes preserve:
- Multiple consecutive spaces
- Tabs and newlines
- Special characters
- Command flags and arguments

```go
// Input and output are identical
input := "docker ps -a --format json"
streamer.StreamWord(os.Stdout, input, 20*time.Millisecond)
// Output: "docker ps -a --format json" (streamed word-by-word)
```

## Testing

```go
func TestMyStreaming(t *testing.T) {
    buf := &bytes.Buffer{}
    text := "test output"
    
    err := streamer.StreamWord(buf, text, 0) // Zero delay for tests
    if err != nil {
        t.Fatal(err)
    }
    
    if buf.String() != text {
        t.Errorf("got %q, want %q", buf.String(), text)
    }
}
```

## Best Practices

1. **Check TTY**: Use `ShouldStream()` to detect if streaming is appropriate
2. **Stop spinner first**: Always stop progress indicators before streaming
3. **Use word mode**: Best balance of speed and readability for commands
4. **Zero delay in tests**: Set delays to 0 for fast test execution
5. **Handle errors**: Check return values from Stream functions
6. **Consider length**: Skip streaming for very short text

## Performance

- **Memory**: Minimal allocation, efficient tokenization
- **CPU**: Negligible (sleeps between outputs)
- **Latency**: Configurable via delay settings
- **Throughput**: ~100 words/second at 10ms delay

## Common Patterns

### Command + Explanation

```go
// Stream command
streamer.StreamWord(os.Stdout, command, 20*time.Millisecond)
fmt.Println() // Newline

// Stream explanation line-by-line
streamer.StreamLine(os.Stdout, explanation, 50*time.Millisecond)
```

### Conditional Streaming

```go
delay := 20 * time.Millisecond
if !streamer.ShouldStream() {
    delay = 0 // Instant if not a terminal
}
streamer.StreamWord(os.Stdout, text, delay)
```

### Error Handling

```go
if err := streamer.StreamWord(os.Stdout, text, 20*time.Millisecond); err != nil {
    // Fallback to instant output
    fmt.Print(text)
}
```
