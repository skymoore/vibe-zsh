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

## Behavior Configuration

### Show/Hide Explanations

Control whether command explanations are displayed (default: `true`):

```bash
export VIBE_SHOW_EXPLANATION=false
```

### Interactive Mode

Require confirmation before inserting commands (default: `false`):

```bash
export VIBE_INTERACTIVE=true
```

When enabled, vibe will show the command and ask for confirmation before inserting it into your prompt.

### Show/Hide Warnings

Control whether warnings are shown for dangerous commands (default: `true`):

```bash
export VIBE_SHOW_WARNINGS=false
```

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

Set maximum response tokens (default: `500`):

```bash
export VIBE_MAX_TOKENS=1000
```

## Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `VIBE_API_URL` | `http://localhost:11434/v1` | API endpoint URL |
| `VIBE_API_KEY` | `""` | API key (if required) |
| `VIBE_MODEL` | `llama3:8b` | Model to use |
| `VIBE_TEMPERATURE` | `0.7` | Generation temperature |
| `VIBE_MAX_TOKENS` | `500` | Max response tokens |
| `VIBE_TIMEOUT` | `30s` | Request timeout |
| `VIBE_SHOW_EXPLANATION` | `true` | Show command explanations |
| `VIBE_SHOW_WARNINGS` | `true` | Show warnings for dangerous commands |
| `VIBE_ENABLE_CACHE` | `true` | Enable response caching |
| `VIBE_CACHE_TTL` | `24h` | Cache lifetime |
| `VIBE_INTERACTIVE` | `false` | Confirm before inserting |

## Example Configuration

Here's a complete example configuration in your `~/.zshrc`:

```bash
# vibe configuration
export VIBE_API_URL="https://api.openai.com/v1"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4"
export VIBE_TEMPERATURE=0.5
export VIBE_SHOW_EXPLANATION=true
export VIBE_INTERACTIVE=false
export VIBE_ENABLE_CACHE=true
export VIBE_CACHE_TTL=24h
```
