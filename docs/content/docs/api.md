---
title: "API Reference"
weight: 5
---

This document describes the internal API and architecture of vibe-zsh for developers who want to understand how it works or contribute to the project.

## Architecture Overview

```
┌─────────────────┐
│   Zsh Plugin    │
│ (vibe.plugin)   │
└────────┬────────┘
         │
         │ Captures buffer
         │ on Ctrl+G
         ▼
┌─────────────────────────────────────────────┐
│         Go Binary (Cobra CLI)               │
│  ┌──────────┬──────────┬─────────────────┐ │
│  │ root.go  │history.go│ version/update  │ │
│  └──────────┴──────────┴─────────────────┘ │
└────────┬────────────────────────────────────┘
         │
         ├──────────────────┬──────────────────┐
         │                  │                  │
         ▼                  ▼                  ▼
┌─────────────────┐  ┌──────────────┐  ┌──────────────┐
│  Cache Layer    │  │  LLM Client  │  │   History    │
│ (internal/cache)│  │(internal/    │  │ (internal/   │
│                 │  │  client)     │  │  history)    │
└─────────────────┘  └──────┬───────┘  └──────────────┘
                            │
         ┌──────────────────┼──────────────────┐
         │                  │                  │
         ▼                  ▼                  ▼
┌─────────────────┐  ┌──────────────┐  ┌──────────────┐
│    Parser       │  │   Schema     │  │   Progress   │
│ (JSON/Text)     │  │  (Prompts)   │  │  (Spinner)   │
└─────────────────┘  └──────────────┘  └──────────────┘
         │
         │ OpenAI-compatible API
         ▼
┌──────────────────────────────────────┐
│       LLM Provider                   │
│  (Ollama/OpenAI/Claude/LM Studio)    │
└──────────────────────────────────────┘
```

## Core Components

### 1. Zsh Plugin (`vibe.plugin.zsh`)

The entry point for user interaction.

**Functions:**

```zsh
vibe()          # Main function - generates command from natural language
vh()            # Opens interactive history browser
vibe-regenerate() # Regenerates last command with same query
```

**Key Variables:**

- `VIBE_PLUGIN_DIR`: Plugin installation directory
- `VIBE_BINARY`: Path to the Go binary
- `BUFFER`: Zsh variable containing current command line

**Keybindings:**

```zsh
bindkey '^G' vibe              # Ctrl+G - Generate command
bindkey '^Xh' vibe-history     # Ctrl+X H - Browse history
bindkey '^Xg' vibe-regenerate  # Ctrl+X G - Regenerate last
```

### 2. Go Binary (Cobra CLI)

Main command-line interface built with Cobra.

**Command Structure:**

```bash
vibe-zsh [flags] <query>
vibe-zsh history [list|clear|last]
vibe-zsh version
vibe-zsh update
vibe-zsh check-update
```

**Flags:**

```
API Configuration:
  --api-url string         API endpoint URL (default: http://localhost:11434/v1)
  --api-key string         API authentication key
  --model string           Model to use (default: llama3:8b)
  --temperature float      Generation temperature 0.0-2.0 (default: 0.2)
  --max-tokens int         Maximum response tokens (default: 1000)
  --timeout duration       Request timeout (default: 30s)

Output Control:
  --structured-output      Use JSON schema for responses (default: true)
  --explanation            Show command explanations (default: true)
  --warnings               Show warnings for dangerous commands (default: true)
  --interactive            Confirm before executing commands (default: false)

Caching:
  --cache                  Enable response caching (default: true)
  --cache-ttl duration     Cache lifetime (default: 24h)

Reliability:
  --max-retries int        Max retry attempts (default: 3)
  --json-extraction        Extract JSON from corrupted responses (default: true)
  --strict-validation      Validate response structure (default: true)
  --retry-status           Show retry progress (default: true)

UI/UX:
  --progress               Show progress spinner (default: true)
  --progress-style string  Spinner style: dots, line, circle, bounce, arrow, runes
  --stream                 Stream output with typewriter effect (default: true)
  --stream-delay duration  Delay between streamed words (default: 20ms)
  --debug                  Enable debug logging (default: false)
```

**Exit Codes:**

- `0`: Success
- `1`: Error (details on stderr)
- `130`: Interrupted (Ctrl+C)

### 3. Configuration (`internal/config/config.go`)

Manages all configuration from environment variables and command-line flags.

**Type:**

```go
type Config struct {
    // API Configuration
    APIURL               string
    APIKey               string
    Model                string
    Temperature          float64
    MaxTokens            int
    Timeout              time.Duration
    
    // Output Control
    UseStructuredOutput  bool
    ShowExplanation      bool
    ShowWarnings         bool
    InteractiveMode      bool
    
    // Caching
    EnableCache          bool
    CacheDir             string
    CacheTTL             time.Duration
    
    // Reliability
    MaxRetries           int
    EnableJSONExtraction bool
    StrictValidation     bool
    ShowRetryStatus      bool
    
    // UI/UX
    ShowProgress         bool
    ProgressStyle        progress.SpinnerStyle
    StreamOutput         bool
    StreamDelay          time.Duration
    EnableDebugLogs      bool
    
    // System Context
    OSName               string  // Detected OS (macOS/Linux/Windows)
    Shell                string  // Detected shell (zsh/bash/fish)
    
    // History
    EnableHistory        bool
    HistorySize          int
    HistoryKey           string
    RegenerateKey        string
}
```

**Loading:**

```go
cfg := config.Load()  // Loads from environment variables
```

**Environment Variables:**

All flags can be set via environment variables with `VIBE_` prefix:
- `VIBE_API_URL`
- `VIBE_MODEL`
- `VIBE_TEMPERATURE`
- `VIBE_SHOW_PROGRESS`
- etc.

### 4. Client (`internal/client/client.go`)

Handles communication with OpenAI-compatible APIs with multi-layer fallback strategy.

**Main Type:**

```go
type Client struct {
    config     *config.Config
    httpClient *http.Client
    cache      *cache.Cache
}
```

**Primary Method:**

```go
func (c *Client) GenerateCommand(ctx context.Context, query string) (*schema.CommandResponse, error)
```

**Multi-Layer Fallback Strategy:**

1. **Layer 1: Structured Output** - Uses JSON schema with strict validation
2. **Layer 2: Enhanced Parsing** - Extracts JSON from text with retries
3. **Layer 3: Explicit JSON Prompt** - Adds explicit JSON formatting instructions
4. **Layer 4: Emergency Fallback** - Returns helpful error message

**Request Structure:**

```go
type ChatCompletionRequest struct {
    Model          string          `json:"model"`
    Messages       []Message       `json:"messages"`
    Temperature    float64         `json:"temperature,omitempty"`
    MaxTokens      int             `json:"max_tokens,omitempty"`
    Stream         bool            `json:"stream"`
    ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type Message struct {
    Role    string `json:"role"`    // "system" or "user"
    Content string `json:"content"`
}

type ResponseFormat struct {
    Type       string      `json:"type"`              // "json_schema"
    JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}
```

**Response Structure:**

```go
type ChatCompletionResponse struct {
    ID      string   `json:"id"`
    Object  string   `json:"object"`
    Created int64    `json:"created"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   Usage    `json:"usage,omitempty"`
}

type Choice struct {
    Index        int     `json:"index"`
    Message      Message `json:"message"`
    FinishReason string  `json:"finish_reason"`
}
```

**Retry Logic (`internal/client/retry.go`):**

```go
func (c *Client) withRetry(ctx context.Context, fn func() error) error
```

- Exponential backoff: 1s → 2s → 4s (max 10s)
- Retries on network errors and 5xx status codes
- Respects context cancellation

### 5. Parser (`internal/parser/text_parser.go`)

Parses LLM responses with robust JSON extraction and text fallback.

**Primary Functions:**

```go
// Extract valid JSON from corrupted/wrapped responses
func ExtractJSON(corruptedText string) (string, error)

// Parse text-based responses (fallback)
func ParseTextResponse(text string) (*schema.CommandResponse, error)

// Remove ANSI codes, Unicode garbage, etc.
func RemoveGarbagePatterns(text string) string

// Attempt to repair malformed JSON
func AttemptJSONRepair(text string) string
```

**JSON Extraction Strategy:**

1. Find first `{` and last `}` - validate as JSON
2. Apply `RemoveGarbagePatterns` to clean text
3. Attempt `AttemptJSONRepair` for malformed JSON
4. Return error if all strategies fail

**Text Parsing (Fallback):**

Parses responses in format:
```
command
# explanation line 1
# explanation line 2
WARNING: dangerous operation
```

Handles:
- Code blocks (```bash)
- Markdown formatting
- Comment prefixes (#, -, *)
- Warning detection

### 6. Cache (`internal/cache/cache.go`)

File-based caching system for responses with TTL support.

**Type:**

```go
type Cache struct {
    dir string
    ttl time.Duration
}
```

**Methods:**

```go
func New(dir string, ttl time.Duration) (*Cache, error)
func (c *Cache) Get(query string) (*schema.CommandResponse, bool)
func (c *Cache) Set(query string, resp *schema.CommandResponse) error
func (c *Cache) Clear() error
```

**Cache Location:**

```
~/.cache/vibe/
```

**Key Generation:**

```go
key := sha256(query)  // Hashed query as filename
```

**Cache Entry Format:**

```json
{
  "query": "original query",
  "response": {
    "command": "...",
    "explanation": ["..."],
    "warning": "..."
  },
  "timestamp": "2024-01-01T12:00:00Z",
  "model": "llama3:8b",
  "temperature": 0.2
}
```

**TTL Handling:**

- Default: 24 hours
- Configurable via `VIBE_CACHE_TTL`
- Expired entries automatically ignored on read

### 7. Schema (`internal/schema/schema.go`)

Defines response structure, JSON schema, and system prompts.

**Response Type:**

```go
type CommandResponse struct {
    Command      string   `json:"command"`
    Explanation  []string `json:"explanation"`
    Warning      string   `json:"warning,omitempty"`
    Alternatives []string `json:"alternatives,omitempty"`
    SafetyLevel  string   `json:"safety_level,omitempty"`
}

func (c *CommandResponse) Validate() error
```

**JSON Schema:**

```go
func GetJSONSchema() map[string]interface{}
```

Returns OpenAI-compatible JSON schema for structured output with:
- Required fields: `command`, `explanation`
- Optional fields: `warning`, `alternatives`, `safety_level`
- Strict validation enabled

**System Prompt:**

```go
func GetSystemPrompt(osName, shell string) string
```

Generates OS and shell-aware system prompt that:
- Includes OS context (macOS/Linux/Windows)
- Includes shell context (zsh/bash/fish)
- Warns about OS-specific utilities (BSD vs GNU)
- Enforces strict JSON formatting rules
- Prevents common LLM mistakes (markdown, Unicode, truncation)

**Key Prompt Features:**

- OS-specific command generation (e.g., BSD `find` on macOS)
- Strict JSON formatting requirements
- Complete explanation enforcement (no "...", "???")
- Warning detection for dangerous commands
- Portable POSIX command preference

### 8. Additional Components

**Progress Spinner (`internal/progress/spinner.go`):**

```go
type Spinner struct {
    frames   []string
    message  string
    // ...
}

func NewSpinner(style SpinnerStyle) *Spinner
func (s *Spinner) Start(ctx context.Context, initialMsg string)
func (s *Spinner) Update(msg string)
func (s *Spinner) Stop()
```

Styles: `dots`, `line`, `circle`, `bounce`, `arrow`, `runes`

**Streamer (`internal/streamer/streamer.go`):**

```go
type Streamer struct {
    config *Config
    writer io.Writer
}

func StreamWord(w io.Writer, text string, delay time.Duration) error
func StreamCharacter(w io.Writer, text string, delay time.Duration) error
func StreamLine(w io.Writer, text string, delay time.Duration) error
```

Provides typewriter effect for output.

**History (`internal/history/history.go`):**

```go
type History struct {
    dir      string
    maxSize  int
    filePath string
}

type Entry struct {
    Query     string    `json:"query"`
    Command   string    `json:"command"`
    Timestamp time.Time `json:"timestamp"`
    Count     int       `json:"count"`
}

func New(cacheDir string, maxSize int) (*History, error)
func (h *History) Add(query, command string) error
func (h *History) List() ([]Entry, error)
func (h *History) Clear() error
```

**Interactive History UI (`internal/history/ui.go`):**

```go
func ShowInteractive(entries []Entry) (*SelectionResult, error)
func FormatPlainList(entries []Entry) string
```

Uses Bubble Tea for TUI.

**Confirmation Dialog (`internal/confirm/confirm.go`):**

```go
func ShowConfirmation(command string) (bool, error)
```

Interactive confirmation prompt for dangerous commands.

**Updater (`internal/updater/updater.go`):**

```go
func CheckForUpdatesBackground(currentVersion string)
func ShowUpdateNotification(currentVersion string)
func PerformUpdate(currentVersion string) error
```

Checks GitHub releases for updates.

**Formatter (`internal/formatter/formatter.go`):**

```go
func Format(resp *schema.CommandResponse, showExplanation bool, showWarnings bool) string
```

Formats response for display.

**Logger (`internal/logger/logger.go`):**

```go
func Init(enableDebug bool)
func Debug(format string, args ...interface{})
func LogParsingFailure(attempt int, layer, content string, err error)
func LogLayerSuccess(layer string, layerNum int)
```

Debug logging system.

**Errors (`internal/errors/errors.go`):**

```go
var (
    ErrInvalidJSON   = errors.New("invalid JSON")
    ErrNoResponse    = errors.New("no response from API")
    ErrEmptyResponse = errors.New("empty response")
)

type APIError struct {
    StatusCode int
    Body       string
}

func IsRetryable(err error) bool
```

## Request Flow

1. **User Input**
   - User types natural language
   - Presses `Ctrl+G`

2. **Zsh Plugin**
   - Captures `$BUFFER`
   - Calls Go binary: `vibe-zsh "$BUFFER"`

3. **Cobra CLI**
   - Parses flags and arguments
   - Initializes configuration

4. **Configuration**
   - Loads environment variables
   - Merges with command-line flags
   - Detects OS and shell

5. **Progress Spinner**
   - Starts spinner if TTY detected
   - Shows "Checking cache..." message

6. **Cache Check**
   - Generates SHA256 hash of query
   - Checks `~/.cache/vibe/` for cached response
   - Returns cached response if valid and not expired

7. **LLM Request (Multi-Layer)**
   - **Layer 1**: Structured output with JSON schema
   - **Layer 2**: Enhanced parsing with retries (up to 3)
   - **Layer 3**: Explicit JSON prompt with lower temperature
   - **Layer 4**: Emergency fallback with error message

8. **HTTP Request**
   - Constructs OpenAI-compatible request
   - Adds authentication header if API key present
   - Applies retry logic with exponential backoff
   - Respects context cancellation (Ctrl+C)

9. **Response Parsing**
   - Extracts JSON from response
   - Validates structure
   - Cleans explanations (removes ANSI codes, Unicode)
   - Detects garbage/truncated text

10. **Output Formatting**
    - Streams explanations to stderr with typewriter effect
    - Detects and warns about incomplete explanations
    - Shows warnings for dangerous commands

11. **Interactive Confirmation** (if enabled)
    - Shows Bubble Tea confirmation dialog
    - User can accept/reject command

12. **Command Output**
    - Writes command to stdout (captured by Zsh)
    - Saves to history if enabled
    - Shows update notification if available

13. **Zsh Plugin**
    - Reads stdout
    - Replaces `$BUFFER` with command
    - User can edit or execute

## Error Handling

### Error Types

**Network Errors:**
- Connection refused → Check if API service is running
- Host not found → Check `VIBE_API_URL`
- Timeout → Check network connection

**API Errors:**
- 401 Unauthorized → Check `VIBE_API_KEY`
- 429 Rate Limited → Automatic retry with backoff
- 500 Server Error → Automatic retry with backoff

**Parsing Errors:**
- Invalid JSON → Falls back to text parsing
- Missing fields → Falls back to next layer
- Validation failure → Falls back to next layer

**Cache Errors:**
- Non-fatal, logged but don't stop execution

**History Errors:**
- Non-fatal, logged but don't stop execution

### Error Output

All errors written to stderr:

```go
fmt.Fprintf(os.Stderr, "Error: %v\n", err)
os.Exit(1)
```

### Retry Logic

Retryable errors:
- Network timeouts
- Connection errors
- 5xx status codes
- Rate limits (429)

Non-retryable errors:
- 4xx client errors (except 429)
- Invalid configuration
- Context cancellation

Backoff strategy:
- Initial: 1 second
- Multiplier: 2x
- Maximum: 10 seconds
- Max attempts: 3

## Testing

### Unit Tests

Run all tests:

```bash
make test
```

Test specific package:

```bash
go test ./internal/parser
go test ./internal/config
go test ./internal/history
```

Test with coverage:

```bash
go test -cover ./...
```

Test with verbose output:

```bash
go test -v ./internal/parser
```

### Integration Tests

Test with live API:

```bash
export VIBE_API_URL="http://localhost:11434/v1"
export VIBE_MODEL="llama3:8b"
vibe-zsh "list all files"
```

Test with debug logging:

```bash
vibe-zsh --debug "find python files"
```

Test caching:

```bash
# First run - should hit API
vibe-zsh "list docker containers"

# Second run - should use cache
vibe-zsh "list docker containers"
```

Test history:

```bash
vibe-zsh history list
vibe-zsh history last
vibe-zsh history clear
```

Test interactive mode:

```bash
vibe-zsh --interactive "remove all logs"
```

### Testing Different Providers

**Ollama:**
```bash
export VIBE_API_URL="http://localhost:11434/v1"
export VIBE_MODEL="llama3:8b"
```

**OpenAI:**
```bash
export VIBE_API_URL="https://api.openai.com/v1"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4"
```

**LM Studio:**
```bash
export VIBE_API_URL="http://localhost:1234/v1"
export VIBE_MODEL="local-model"
```

## Development

### Building

```bash
make build        # Current platform
make build-all    # All platforms (Linux, macOS, Windows)
make install      # Build and install to /usr/local/bin
```

Build manually:

```bash
go build -o vibe-zsh -ldflags "-X main.version=dev" .
```

### Project Structure

```
vibe-zsh/
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Main command + generation logic
│   └── history.go         # History subcommands
├── internal/
│   ├── cache/             # Response caching
│   ├── client/            # HTTP client + retry logic
│   │   ├── client.go      # Main client with fallback layers
│   │   ├── http.go        # HTTP request handling
│   │   └── retry.go       # Exponential backoff retry
│   ├── config/            # Configuration management
│   ├── confirm/           # Interactive confirmation dialog
│   ├── errors/            # Error types and handling
│   ├── formatter/         # Response formatting
│   ├── history/           # Query history management
│   │   ├── history.go     # History storage
│   │   └── ui.go          # Interactive TUI
│   ├── logger/            # Debug logging
│   ├── parser/            # JSON extraction + text parsing
│   ├── progress/          # Spinner animations
│   ├── schema/            # Response schema + prompts
│   ├── streamer/          # Typewriter effect output
│   └── updater/           # Auto-update functionality
├── main.go                # Entry point
└── vibe.plugin.zsh        # Zsh integration
```

### Adding a New Feature

**Example: Adding a new spinner style**

1. Add style to `internal/progress/spinner.go`:
```go
const StyleCustom SpinnerStyle = "custom"

var spinnerFrames = map[SpinnerStyle][]string{
    // ...
    StyleCustom: {"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
}
```

2. Update config parser in `cmd/root.go`:
```go
case "custom":
    return progress.StyleCustom
```

3. Update documentation

**Example: Adding a new command**

1. Create new file in `cmd/`:
```go
package cmd

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

2. Add tests
3. Update documentation

### Extending Functionality

**Adding a new configuration option:**

1. Add to `Config` struct in `internal/config/config.go`
2. Add environment variable parsing with `getEnv*` helper
3. Add command-line flag in `cmd/root.go`
4. Add flag handling in `initConfig()`
5. Update documentation
6. Add tests

**Modifying the system prompt:**

1. Edit `GetSystemPrompt()` in `internal/schema/schema.go`
2. Test with multiple LLM providers (Ollama, OpenAI, Claude)
3. Ensure OS-specific instructions work correctly
4. Test with different shells (zsh, bash)
5. Verify JSON output format compliance

**Adding a new parsing layer:**

1. Add method to `Client` in `internal/client/client.go`
2. Call in `GenerateCommand()` fallback chain
3. Add logging with `logger.LogLayerSuccess()`
4. Add error handling with `logger.LogParsingFailure()`
5. Test with various LLM outputs

## Contributing

See the main README for contribution guidelines.

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` before committing
- Add tests for new functionality
- Update documentation
- Keep functions focused and small
- Use meaningful variable names
- Add comments for complex logic

### Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests (`go test ./...`)
5. Update documentation (README, API docs, usage docs)
6. Commit with clear messages
7. Push to your fork
8. Submit PR with description of changes

### Testing Checklist

Before submitting PR:

- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`gofmt -w .`)
- [ ] No lint errors (`go vet ./...`)
- [ ] Documentation updated
- [ ] Tested with Ollama
- [ ] Tested with OpenAI (if applicable)
- [ ] Tested on macOS (if applicable)
- [ ] Tested on Linux (if applicable)

## Dependencies

vibe-zsh uses minimal external dependencies:

**Core Dependencies:**
- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - TUI components
- `github.com/charmbracelet/lipgloss` - Terminal styling

**Standard Library:**
- `net/http` - HTTP client
- `encoding/json` - JSON parsing
- `crypto/sha256` - Cache key generation
- `context` - Cancellation and timeouts
- `os`, `io`, `time`, `strings`, etc.

**Why These Dependencies?**

- **Cobra**: Industry-standard CLI framework with excellent flag parsing
- **Bubble Tea**: Modern, composable TUI framework for interactive features
- **Lipgloss**: Beautiful terminal styling without complexity

All dependencies are well-maintained, widely-used, and have minimal transitive dependencies.

## API Compatibility

vibe-zsh is compatible with any OpenAI-compatible API:

**Required Endpoints:**
- `POST /chat/completions`

**Required Request Format:**
```json
{
  "model": "string",
  "messages": [
    {"role": "system", "content": "string"},
    {"role": "user", "content": "string"}
  ],
  "temperature": 0.2,
  "max_tokens": 1000,
  "stream": false,
  "response_format": {
    "type": "json_schema",
    "json_schema": {...}
  }
}
```

**Required Response Format:**
```json
{
  "id": "string",
  "object": "chat.completion",
  "created": 1234567890,
  "model": "string",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "string"
      },
      "finish_reason": "stop"
    }
  ]
}
```

**Tested Providers:**
- ✅ Ollama
- ✅ OpenAI
- ✅ LM Studio
- ✅ Claude (via OpenAI-compatible proxy)
- ✅ Anthropic API
- ✅ Azure OpenAI

## Performance Considerations

**Caching:**
- Cache hits return in <1ms
- Cache misses require API call (1-5s depending on model)
- Cache stored in `~/.cache/vibe/`
- Default TTL: 24 hours

**Network:**
- Timeout: 30s default
- Retry: 3 attempts with exponential backoff
- Connection pooling via `http.Client`

**Memory:**
- Minimal memory footprint (~10MB)
- No persistent connections
- Cache files are small (<1KB each)

**Startup:**
- Cold start: ~50ms
- Warm start (cached): <10ms

## Security Considerations

**API Keys:**
- Never logged or printed
- Stored in environment variables only
- Transmitted via HTTPS only
- Not included in cache files

**Command Execution:**
- Commands never auto-executed
- Interactive mode available for confirmation
- Warnings shown for dangerous commands
- User always has final control

**Cache:**
- Stored in user's home directory
- Readable only by user (0644 permissions)
- No sensitive data stored
- Can be cleared anytime

**Network:**
- HTTPS enforced for remote APIs
- Certificate validation enabled
- No telemetry or tracking
- All requests user-initiated
