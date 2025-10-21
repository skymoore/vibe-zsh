---
title: "Configuration"
weight: 2
---


All configuration is done through environment variables in your `~/.zshrc`. All settings are optional and have sensible defaults.

## Provider Configuration

### Local LLM (Ollama) - Default

```bash
export VIBE_API_URL="http://localhost:11434/v1"
export VIBE_MODEL="llama3:8b"
```

This is the default configuration. If you have Ollama running locally, vibe will work out of the box.

### OpenAI

```bash
export VIBE_API_URL="https://api.openai.com/v1"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4"
```

### Anthropic (via OpenRouter)

```bash
export VIBE_API_URL="https://openrouter.ai/api/v1"
export VIBE_API_KEY="sk-or-..."
export VIBE_MODEL="anthropic/claude-3.5-sonnet"
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

## Display Configuration

### Progress Indicators

Control the progress spinner shown during command generation (default: `true`):

```bash
export VIBE_SHOW_PROGRESS=false
```

Choose spinner style (default: `dots`):

```bash
export VIBE_PROGRESS_STYLE=arrow  # Options: dots, line, circle, bounce, arrow
```

### Streaming Output

Enable typewriter effect for command output (default: `true`):

```bash
export VIBE_STREAM_OUTPUT=false
```

Control streaming speed (default: `20ms`):

```bash
export VIBE_STREAM_DELAY=10ms  # Faster
export VIBE_STREAM_DELAY=50ms  # Slower
```

### Show/Hide Explanations

Control whether command explanations are displayed (default: `true`):

```bash
export VIBE_SHOW_EXPLANATION=false
```

### Show/Hide Warnings

Control whether warnings are shown for dangerous commands (default: `true`):

```bash
export VIBE_SHOW_WARNINGS=false
```

## Behavior Configuration

### Interactive Mode

Require confirmation before inserting commands (default: `false`):

```bash
export VIBE_INTERACTIVE=true
```

When enabled, vibe will show the command and ask for confirmation before inserting it into your prompt.

## Performance Configuration

### Cache Settings

Enable/disable response caching (default: `true`):

```bash
export VIBE_ENABLE_CACHE=false
```

Set cache lifetime (default: `24h`):

```bash
export VIBE_CACHE_TTL=12h
```

Caching dramatically improves performance for repeated queries. Cached responses are 100-400x faster.

### API Settings

Set request timeout (default: `30s`):

```bash
export VIBE_TIMEOUT=60s
```

Set generation temperature (default: `0.7`):

```bash
export VIBE_TEMPERATURE=0.3
```

Lower temperature = more deterministic, higher = more creative.

Set maximum response tokens (default: `1000`):

```bash
export VIBE_MAX_TOKENS=1000
```

## Advanced Configuration

### Parsing & Reliability

Control how vibe handles LLM responses (defaults are recommended):

```bash
# Max retry attempts for failed parsing (default: 3)
export VIBE_MAX_RETRIES=5

# Extract JSON from corrupted responses (default: true)
export VIBE_ENABLE_JSON_EXTRACTION=true

# Validate response structure (default: true)
export VIBE_STRICT_VALIDATION=true

# Use JSON schema for structured responses (default: true)
export VIBE_USE_STRUCTURED_OUTPUT=true
```

### Debug & Troubleshooting

Enable detailed logging for debugging parsing issues:

```bash
# Enable debug logging (default: false)
export VIBE_DEBUG_LOGS=true

# Show retry progress during generation (default: true)
export VIBE_SHOW_RETRY_STATUS=true
```

When `VIBE_DEBUG_LOGS=true`, vibe will log:
- Raw LLM responses
- Parsing attempts and failures
- Which fallback layer succeeded
- Extracted/cleaned JSON content

## Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| **API Configuration** | | |
| `VIBE_API_URL` | `http://localhost:11434/v1` | API endpoint URL |
| `VIBE_API_KEY` | `""` | API key (if required) |
| `VIBE_MODEL` | `llama3:8b` | Model to use |
| `VIBE_TEMPERATURE` | `0.2` | Generation temperature (0.0-2.0) |
| `VIBE_MAX_TOKENS` | `1000` | Max response tokens |
| `VIBE_TIMEOUT` | `30s` | Request timeout |
| **Display Options** | | |
| `VIBE_SHOW_EXPLANATION` | `true` | Show command explanations |
| `VIBE_SHOW_WARNINGS` | `true` | Show warnings for dangerous commands |
| `VIBE_SHOW_PROGRESS` | `true` | Show progress spinner during generation |
| `VIBE_PROGRESS_STYLE` | `dots` | Spinner style: dots, line, circle, bounce, arrow |
| `VIBE_STREAM_OUTPUT` | `true` | Stream output with typewriter effect |
| `VIBE_STREAM_DELAY` | `20ms` | Delay between streamed words |
| **Behavior** | | |
| `VIBE_INTERACTIVE` | `false` | Confirm before inserting command |
| `VIBE_USE_STRUCTURED_OUTPUT` | `true` | Use JSON schema for structured responses |
| `VIBE_ENABLE_CACHE` | `true` | Enable response caching |
| `VIBE_CACHE_TTL` | `24h` | Cache lifetime |
| **Parsing & Retry** | | |
| `VIBE_MAX_RETRIES` | `3` | Max retry attempts for failed parsing |
| `VIBE_ENABLE_JSON_EXTRACTION` | `true` | Extract JSON from corrupted responses |
| `VIBE_STRICT_VALIDATION` | `true` | Validate response structure |
| `VIBE_SHOW_RETRY_STATUS` | `true` | Show retry progress during generation |
| **Debugging** | | |
| `VIBE_DEBUG_LOGS` | `false` | Enable debug logging for troubleshooting |

## OS-Aware Command Generation

vibe automatically detects your operating system and shell to generate commands that work on YOUR system:

- **macOS**: Uses BSD utilities (e.g., `find` without `-printf`, `sed -i ''`)
- **Linux**: Uses GNU utilities (e.g., `find -printf`, `sed -i`)
- **Shell-specific**: Generates syntax appropriate for zsh, bash, etc.

This detection happens automatically - no configuration needed! Commands are tailored to your environment.

## Example Configuration

Here's a complete example configuration in your `~/.zshrc`:

```bash
# vibe configuration
export VIBE_API_URL="https://api.openai.com/v1"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4"
export VIBE_TEMPERATURE=0.2
export VIBE_SHOW_EXPLANATION=true
export VIBE_SHOW_PROGRESS=true
export VIBE_PROGRESS_STYLE=dots
export VIBE_STREAM_OUTPUT=true
export VIBE_STREAM_DELAY=20ms
export VIBE_INTERACTIVE=false
export VIBE_ENABLE_CACHE=true
export VIBE_CACHE_TTL=24h
```
