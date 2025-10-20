# vibe-zsh

**vibe-zsh** is an Oh-My-Zsh plugin that transforms natural language into shell commands using any OpenAI-compatible LLM API. Written in Go for speed and zero runtime dependencies.

> "list all docker containers running"  
> ⤷ becomes: `docker ps`  
> ⤷ plus inline explanation of the command and flags.

---

## 🚀 Features

- 🧠 Natural language → shell command translation
- 💬 Inline explanation of each command and flag
- ✍️ Editable buffer before executing anything
- 🔌 OpenAI API compatible (works with Ollama, OpenAI, LM Studio, Groq, OpenRouter, etc.)
- ⚡ Fast compiled Go binary (no Python/Node.js required)
- 📋 Structured JSON output for reliability
- 🛡️ Built-in safety warnings for dangerous commands

---

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development Plan](#development-plan)
- [Technical Details](#technical-details)
- [Contributing](#contributing)
- [License](#license)

---

## 🏗️ Architecture Overview

### Component Stack

1. **Go Binary** (`vibe`) - Calls OpenAI-compatible API with JSON schema
2. **Zsh Plugin** (`vibe.plugin.zsh`) - ZLE widget that captures buffer and calls binary
3. **LLM Provider** - Any OpenAI-compatible API endpoint

### Data Flow

```
User types natural language
         ↓
    Press Ctrl+G
         ↓
Zsh captures $BUFFER
         ↓
Calls Go binary with query
         ↓
Go sends to LLM with JSON schema
         ↓
LLM returns structured JSON
         ↓
Go formats as command + comments
         ↓
Zsh replaces buffer with output
         ↓
User reviews and executes
```

---

## 📦 Installation

### Prerequisites

- Zsh with Oh-My-Zsh installed
- Go 1.21+ (for building from source)
- An OpenAI-compatible LLM provider running (e.g., Ollama)

### Quick Install

```bash
# Clone into Oh-My-Zsh custom plugins directory
git clone https://github.com/skymoore/vibe-zsh.git ~/.oh-my-zsh/custom/plugins/vibe

# Build the binary
cd ~/.oh-my-zsh/custom/plugins/vibe
make build

# Add to your .zshrc plugins list
plugins=(git vibe)

# Reload your shell
source ~/.zshrc
```

### Manual Build

```bash
# Clone the repo
git clone https://github.com/skymoore/vibe-zsh.git
cd vibe-zsh

# Build for your platform
go build -o vibe main.go

# Move binary to plugin directory or somewhere in PATH
mv vibe ~/.oh-my-zsh/custom/plugins/vibe/
```

---

## ⚙️ Configuration

Configure vibe using environment variables in your `~/.zshrc`:

### Default (Ollama Local)

```bash
export VIBE_API_URL="http://localhost:11434/v1"
export VIBE_MODEL="llama3:8b"
```

### OpenAI

```bash
export VIBE_API_URL="https://api.openai.com/v1"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4"
```

### LM Studio

```bash
export VIBE_API_URL="http://localhost:1234/v1"
export VIBE_MODEL="local-model"
```

### Groq

```bash
export VIBE_API_URL="https://api.groq.com/openai/v1"
export VIBE_API_KEY="gsk_..."
export VIBE_MODEL="llama-3.1-70b-versatile"
```

### OpenRouter

```bash
export VIBE_API_URL="https://openrouter.ai/api/v1"
export VIBE_API_KEY="sk-or-..."
export VIBE_MODEL="anthropic/claude-3.5-sonnet"
```

### Azure OpenAI

```bash
export VIBE_API_URL="https://<resource>.openai.azure.com/openai/deployments/<deployment>"
export VIBE_API_KEY="<azure-key>"
export VIBE_MODEL="gpt-4"
```

### All Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `VIBE_API_URL` | `http://localhost:11434/v1` | Base URL for OpenAI-compatible API |
| `VIBE_API_KEY` | `""` | API key (optional for local, required for cloud) |
| `VIBE_MODEL` | `llama3:8b` | Model name to use |
| `VIBE_TEMPERATURE` | `0.7` | Temperature for generation (0.0-2.0) |
| `VIBE_MAX_TOKENS` | `500` | Maximum tokens in response |
| `VIBE_TIMEOUT` | `30s` | Request timeout duration |
| `VIBE_USE_STRUCTURED_OUTPUT` | `true` | Enable JSON schema structured output |
| `VIBE_SHOW_EXPLANATION` | `true` | Show explanation comments below command |

---

## 🎯 Usage

### Basic Usage

1. Type a natural language command in your terminal
2. Press `Ctrl+G` (or your configured keybinding)
3. vibe replaces your input with the shell command + explanation
4. Review the command (edit if needed)
5. Press Enter to execute

### Examples

**Example 1: List files**
```bash
# Type:
show me all files including hidden ones

# Press Ctrl+G, get:
ls -la
# ls: list directory contents
# -l: use a long listing format
# -a: include hidden files
```

**Example with explanations disabled:**
```bash
# In ~/.zshrc:
export VIBE_SHOW_EXPLANATION=false

# Type:
list all files

# Press Ctrl+G, get:
ls -la
```

**Example 2: Docker**
```bash
# Type:
show logs of nginx container

# Press Ctrl+G, get:
docker logs -f nginx
# docker: Docker command-line interface
# logs: Show logs for a container
# -f: Follow log output in real-time
# nginx: Name of the container
```

**Example 3: Find files**
```bash
# Type:
find all python files modified in the last week

# Press Ctrl+G, get:
find . -name "*.py" -mtime -7
# find: Search for files and directories
# .: Start search in current directory
# -name "*.py": Match files ending in .py
# -mtime -7: Modified within last 7 days
```

**Example 4: System commands (with warning)**
```bash
# Type:
delete all files in /tmp

# Press Ctrl+G, get:
rm -rf /tmp/*
# rm: Remove files or directories
# -r: Remove directories recursively
# -f: Force removal without prompting
# /tmp/*: All files in /tmp directory
# WARNING: This command will permanently delete all files in /tmp
```

---

## 🗺️ Development Plan

### Phase 0: De-risk Integration ✅ (CRITICAL FIRST)

**Goal:** Verify zsh can handle multi-line buffer replacement

- [ ] Create minimal Go program that outputs multi-line string
- [ ] Test zsh `$BUFFER` replacement with multi-line content
- [ ] Verify comment lines work correctly in zsh buffer
- [ ] Document any quirks or limitations

**Success Criteria:** Simple Go program can successfully replace zsh buffer with command + comments

---

### Phase 1: MVP (Core Functionality)

**Goal:** Working plugin with basic features

#### 1.1 Project Structure
```
vibe-zsh/
├── main.go
├── cmd/
│   └── vibe.go              # CLI logic
├── internal/
│   ├── client/
│   │   └── openai.go        # OpenAI-compatible API client
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── schema/
│   │   └── command_schema.go # JSON schema definition
│   └── formatter/
│       └── formatter.go     # Output formatting
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── Makefile
└── vibe.plugin.zsh
```

#### 1.2 Core Components

**JSON Schema for Structured Output:**
```json
{
  "type": "object",
  "properties": {
    "command": {
      "type": "string",
      "description": "The shell command"
    },
    "explanation": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Explanation lines"
    },
    "warning": {
      "type": "string",
      "description": "Optional warning message"
    }
  },
  "required": ["command", "explanation"],
  "additionalProperties": false
}
```

**OpenAI API Request Format:**
```json
{
  "model": "llama3:8b",
  "messages": [
    {
      "role": "system",
      "content": "You are VibeCLI, a helpful assistant..."
    },
    {
      "role": "user",
      "content": "list all files"
    }
  ],
  "response_format": {
    "type": "json_schema",
    "json_schema": {
      "name": "shell_command_response",
      "strict": true,
      "schema": { /* schema above */ }
    }
  },
  "temperature": 0.7,
  "max_tokens": 500,
  "stream": false
}
```

#### 1.3 Implementation Tasks

- [ ] Initialize Go module
- [ ] Implement config package (environment variables)
- [ ] Implement OpenAI client with HTTP requests
- [ ] Define JSON schema in code
- [ ] Implement request builder with structured output
- [ ] Implement response parser
- [ ] Implement output formatter (command + comments)
- [ ] Create CLI entry point
- [ ] Update vibe.plugin.zsh to call Go binary
- [ ] Add error handling for common cases
- [ ] Create Makefile for building

#### 1.4 System Prompt

```
You are VibeCLI, a helpful assistant for command-line users.

Generate a shell command that accomplishes the user's request.

Guidelines:
- Provide the exact command in the "command" field
- Break down the explanation into clear, concise lines
- Each explanation line should describe one part of the command
- If the command uses sudo or is potentially dangerous, add a warning
- Be precise and accurate - users will run these commands

User request: {user_query}
```

---

### Phase 2: Robustness

**Goal:** Production-ready reliability and compatibility

#### 2.1 Fallback Support

- [ ] Implement text-based response parser (for non-structured models)
- [ ] Auto-detect structured output support
- [ ] Graceful fallback when JSON schema not supported
- [ ] Add `VIBE_USE_STRUCTURED_OUTPUT` flag

#### 2.2 Error Handling

- [ ] Handle network errors gracefully
- [ ] Handle timeout errors
- [ ] Handle JSON parse errors
- [ ] Handle empty responses
- [ ] Handle rate limit errors (429)
- [ ] Handle authentication errors (401)
- [ ] User-friendly error messages

#### 2.3 Retry Logic

- [ ] Implement exponential backoff for transient errors
- [ ] Configurable retry attempts
- [ ] Rate limit handling with backoff

#### 2.4 Testing

- [ ] Unit tests for config loading
- [ ] Unit tests for JSON parsing
- [ ] Unit tests for text parsing (fallback)
- [ ] Unit tests for formatting
- [ ] Integration tests with mocked API
- [ ] Integration tests for error scenarios
- [ ] E2E tests with real Ollama instance

#### 2.5 Build & Distribution

- [ ] Cross-compilation for macOS (amd64, arm64)
- [ ] Cross-compilation for Linux (amd64, arm64)
- [ ] CI/CD pipeline for releases
- [ ] Installation script
- [ ] Release binaries on GitHub

---

### Phase 3: V2 Features (Future)

**Goal:** Enhanced user experience and advanced features

#### 3.1 Tab Completion

- [ ] Research zsh completion system (`compinit`, `compdef`)
- [ ] Implement completion function `_vibe`
- [ ] Double-tab to show command suggestions
- [ ] Show recent commands on completion
- [ ] Show command categories/templates

#### 3.2 Command History

- [ ] Store generated commands in history file
- [ ] Search previous commands
- [ ] Reuse/modify previous queries
- [ ] Analytics on most common commands

#### 3.3 Interactive Mode

- [ ] Confirmation prompt before executing dangerous commands
- [ ] Show command preview before replacement
- [ ] Allow editing before insertion
- [ ] Multi-step command generation

#### 3.4 Extended Schema

- [ ] Add `alternatives` field (2-3 command options)
- [ ] Add `safety_level` field (safe, caution, dangerous)
- [ ] Add `explanation_detail` option (brief, detailed)
- [ ] Support for multi-line commands (scripts)

#### 3.5 Configuration File

- [ ] Support `.viberc` config file
- [ ] Multiple provider profiles
- [ ] Switch profiles easily
- [ ] Per-directory configuration

---

## 🔧 Technical Details

### Why Go?

- ✅ **Single compiled binary** - No runtime dependencies
- ✅ **Fast startup time** - Critical for shell responsiveness
- ✅ **Excellent HTTP client** - stdlib `net/http` is battle-tested
- ✅ **Easy cross-compilation** - Build for any platform
- ✅ **Small binary size** - ~5-10MB statically compiled

### Why OpenAI API Standard?

- ✅ **Universal compatibility** - Works with 10+ providers
- ✅ **Future-proof** - Industry standard for LLM APIs
- ✅ **Local-first** - Ollama supports it out of the box
- ✅ **Cloud-ready** - Easy to switch to OpenAI/Anthropic
- ✅ **Well-documented** - Extensive examples and libraries

### Why JSON Schema Structured Output?

- ✅ **Guaranteed format** - No string parsing needed
- ✅ **Type safety** - Invalid JSON = clear failure
- ✅ **Easier parsing** - Direct JSON decode
- ✅ **Extensible** - Easy to add new fields
- ✅ **Better reliability** - No regex/text manipulation

### Fallback Strategy

For providers/models that don't support structured output:

1. Try structured output first
2. If it fails, automatically fall back to text parsing
3. Parse response looking for command and explanation patterns
4. Clean up markdown code blocks, backticks, etc.
5. Format consistently with structured output

### Security Considerations

- ✅ **No auto-execution** - Commands appear in buffer for review
- ✅ **Warning system** - Dangerous commands flagged
- ✅ **Local-first default** - Data stays on your machine with Ollama
- ✅ **No telemetry** - Zero data collection
- ✅ **Open source** - Audit the code yourself

---

## 🤝 Contributing

Contributions are welcome! Whether it's:

- Adding new provider configurations
- Improving the system prompt
- Adding tests
- Fixing bugs
- Improving documentation

### Development Setup

```bash
# Clone the repo
git clone https://github.com/skymoore/vibe-zsh.git
cd vibe-zsh

# Install dependencies (none! stdlib only)
go mod download

# Run tests
go test ./...

# Build
go build -o vibe main.go

# Run
./vibe "your natural language query"
```

### Testing with Different Providers

```bash
# Test with Ollama
export VIBE_API_URL="http://localhost:11434/v1"
./vibe "list files"

# Test with OpenAI
export VIBE_API_URL="https://api.openai.com/v1"
export VIBE_API_KEY="sk-..."
./vibe "list files"
```

---

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Original Python implementation: [LoganPederson/vibe](https://github.com/LoganPederson/vibe)
- Inspired by natural language shell tools
- Built with ❤️ using Go and Zsh

---

## 📚 Resources

- [OpenAI API Documentation](https://platform.openai.com/docs/api-reference)
- [Ollama Documentation](https://ollama.ai/docs)
- [LM Studio Structured Output Docs](https://lmstudio.ai/docs/advanced/structured-output)
- [Oh-My-Zsh Plugin Development](https://github.com/ohmyzsh/ohmyzsh/wiki/Customization)
- [Zsh Line Editor (ZLE) Documentation](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html)

---

## 🚀 Project Status

**Current Phase:** Planning Complete ✅  
**Next Milestone:** Phase 0 - De-risk Integration  
**Target MVP Date:** TBD  

### Roadmap

- [x] Research and planning
- [x] OpenAI API compatibility design
- [x] JSON schema structured output design
- [ ] Phase 0: Zsh integration proof-of-concept
- [ ] Phase 1: MVP implementation
- [ ] Phase 2: Robustness and testing
- [ ] Phase 3: V2 features

---

**Made with 🤖 and ☕ for the command line.**
