---
title: "API Reference"
weight: 5
---


This document describes the internal API and architecture of vibe for developers who want to understand how it works or contribute to the project.

## Architecture Overview

```
┌─────────────────┐
│   Zsh Plugin    │
│  (vibe.plugin)  │
└────────┬────────┘
         │
         │ Captures buffer
         │ on Ctrl+G
         ▼
┌─────────────────┐
│  Go Binary      │
│  (cmd/vibe.go)  │
└────────┬────────┘
         │
         ├─────────────────┐
         │                 │
         ▼                 ▼
┌─────────────────┐  ┌──────────────┐
│  Cache Layer    │  │  LLM Client  │
│ (internal/cache)│  │(internal/    │
│                 │  │  client)     │
└─────────────────┘  └──────┬───────┘
                            │
                            │ OpenAI API
                            ▼
                     ┌──────────────┐
                     │  LLM Provider│
                     │ (Ollama/GPT) │
                     └──────────────┘
```

## Core Components

### 1. Zsh Plugin (`vibe.plugin.zsh`)

The entry point for user interaction.

**Functions:**

```zsh
vibe()
```

Captures the current command buffer, sends it to the Go binary, and replaces the buffer with the generated command.

**Key Variables:**

- `VIBE_PLUGIN_DIR`: Plugin installation directory
- `VIBE_BINARY`: Path to the Go binary
- `BUFFER`: Zsh variable containing current command line

**Keybinding:**

```zsh
bindkey '^G' vibe
```

### 2. Go Binary (`cmd/vibe.go`)

Main command-line interface.

**Command Structure:**

```bash
vibe [flags] <query>
```

**Flags:**

None currently. Configuration is done via environment variables.

**Exit Codes:**

- `0`: Success
- `1`: Error (details on stderr)

### 3. Configuration (`internal/config/config.go`)

Manages all configuration from environment variables.

**Type:**

```go
type Config struct {
    APIURL          string
    APIKey          string
    Model           string
    Temperature     float64
    MaxTokens       int
    Timeout         time.Duration
    ShowExplanation bool
    ShowWarnings    bool
    EnableCache     bool
    CacheTTL        time.Duration
    Interactive     bool
}
```

**Loading:**

```go
cfg := config.Load()
```

### 4. Client (`internal/client/client.go`)

Handles communication with OpenAI-compatible APIs.

**Interface:**

```go
type Client interface {
    Complete(ctx context.Context, prompt string) (*Response, error)
}
```

**Request Structure:**

```go
type CompletionRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    Temperature float64   `json:"temperature"`
    MaxTokens   int       `json:"max_tokens"`
}
```

**Response Structure:**

```go
type Response struct {
    Command     string
    Explanation []string
    Warning     string
}
```

### 5. Parser (`internal/parser/text_parser.go`)

Parses LLM responses into structured command + explanations.

**Interface:**

```go
type Parser interface {
    Parse(response string) (*ParsedCommand, error)
}
```

**Output:**

```go
type ParsedCommand struct {
    Command     string
    Explanation []string
    IsValid     bool
}
```

**Format:**

The parser expects responses in this format:

```
command
# explanation line 1
# explanation line 2
```

### 6. Cache (`internal/cache/cache.go`)

File-based caching system for responses.

**Interface:**

```go
type Cache interface {
    Get(key string) (string, bool)
    Set(key string, value string, ttl time.Duration) error
    Clear() error
}
```

**Cache Location:**

```
~/.cache/vibe/
```

**Key Generation:**

```go
key := sha256(query + model + temperature)
```

### 7. Schema (`internal/schema/schema.go`)

Defines the system prompt sent to the LLM.

**System Prompt:**

```go
const SystemPrompt = `You are a shell command generator...`
```

This prompt instructs the LLM to:
- Generate only valid shell commands
- Provide explanations as comments
- Warn about dangerous operations
- Be concise and accurate

## Request Flow

1. **User Input**
   - User types natural language
   - Presses `Ctrl+G`

2. **Zsh Plugin**
   - Captures `$BUFFER`
   - Calls Go binary: `vibe "$BUFFER"`

3. **Configuration**
   - Loads environment variables
   - Validates settings

4. **Cache Check**
   - Generates cache key from query + config
   - Returns cached response if available

5. **LLM Request**
   - Constructs completion request
   - Sends to configured API
   - Applies retry logic with backoff

6. **Response Parsing**
   - Extracts command and explanations
   - Validates format

7. **Output**
   - Writes to stdout
   - Go binary exits

8. **Zsh Plugin**
   - Reads stdout
   - Replaces `$BUFFER` with command
   - Shows explanations as comments

## Error Handling

All errors are returned on stderr with descriptive messages:

```go
if err != nil {
    fmt.Fprintf(os.Stderr, "vibe: %v\n", err)
    os.Exit(1)
}
```

## Testing

### Unit Tests

Run all tests:

```bash
make test
```

Test specific package:

```bash
go test ./internal/parser
```

### Integration Tests

Test with live API:

```bash
export VIBE_API_URL="http://localhost:11434/v1"
export VIBE_MODEL="llama3:8b"
./vibe "list all files"
```

## Development

### Building

```bash
make build        # Current platform
make build-all    # All platforms
```

### Adding a New Provider

1. Implement the `Client` interface
2. Add provider-specific authentication
3. Map provider response format to `Response` struct
4. Add configuration options
5. Update documentation

### Extending Functionality

**Adding a new configuration option:**

1. Add to `Config` struct in `internal/config/config.go`
2. Add environment variable parsing
3. Update documentation
4. Add tests

**Modifying the system prompt:**

1. Edit `internal/schema/schema.go`
2. Test with multiple LLM providers
3. Ensure backward compatibility

## Contributing

See the main README for contribution guidelines.

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Add tests for new functionality
- Update documentation

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Update documentation
6. Submit PR

## Dependencies

vibe uses **only** the Go standard library:

- `net/http`: HTTP client
- `encoding/json`: JSON parsing
- `crypto/sha256`: Cache key generation
- `os`, `io`, `time`, `context`: Standard utilities

No external dependencies = simple, secure, maintainable code.
