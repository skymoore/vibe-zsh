---
title: "Configuration Reference"
weight: 6
---


Complete reference for all vibe configuration options.

## Environment Variables

All configuration is done through environment variables in your `~/.zshrc`.

### API Configuration

#### VIBE_API_URL

**Type:** String  
**Default:** `http://localhost:11434/v1`  
**Description:** The base URL for the OpenAI-compatible API endpoint.

**Examples:**

```bash 
# Ollama (local)
export VIBE_API_URL="http://localhost:11434/v1"

# OpenAI
export VIBE_API_URL="https://api.openai.com/v1"

# Anthropic via OpenRouter
export VIBE_API_URL="https://openrouter.ai/api/v1"

# LM Studio
export VIBE_API_URL="http://localhost:1234/v1"

# Groq
export VIBE_API_URL="https://api.groq.com/openai/v1"
```

---

#### VIBE_API_KEY

**Type:** String  
**Default:** `""` (empty, not required for Ollama)  
**Description:** API key for authentication. Required for most cloud providers.

**Examples:**

```bash 
# OpenAI
export VIBE_API_KEY="sk-..."

# OpenRouter
export VIBE_API_KEY="sk-or-..."

# Groq
export VIBE_API_KEY="gsk_..."
```

**Security Note:** Never commit API keys to version control.

---

#### VIBE_MODEL

**Type:** String  
**Default:** `llama3:8b`  
**Description:** The model identifier to use for command generation.

**Examples:**

```bash 
# Ollama
export VIBE_MODEL="llama3:8b"
export VIBE_MODEL="llama3.1:70b"
export VIBE_MODEL="codellama:13b"

# OpenAI
export VIBE_MODEL="gpt-4"
export VIBE_MODEL="gpt-4-turbo"
export VIBE_MODEL="gpt-3.5-turbo"

# Anthropic (via OpenRouter)
export VIBE_MODEL="anthropic/claude-3.5-sonnet"
export VIBE_MODEL="anthropic/claude-3-opus"

# Groq
export VIBE_MODEL="llama-3.1-70b-versatile"
export VIBE_MODEL="mixtral-8x7b-32768"
```

---

#### VIBE_TEMPERATURE

**Type:** Float  
**Default:** `0.7`  
**Range:** `0.0` to `2.0`  
**Description:** Controls randomness in the model's output. Lower values make output more deterministic, higher values more creative.

**Examples:**

```bash 
# More deterministic (recommended for commands)
export VIBE_TEMPERATURE=0.3

# Balanced (default)
export VIBE_TEMPERATURE=0.7

# More creative
export VIBE_TEMPERATURE=1.0
```

**Recommendations:**
- `0.0-0.3`: Most consistent, predictable commands
- `0.4-0.7`: Balanced creativity and consistency
- `0.8-2.0`: More variation (may produce unexpected results)

---

#### VIBE_MAX_TOKENS

**Type:** Integer  
**Default:** `500`  
**Description:** Maximum number of tokens in the API response.

**Examples:**

```bash 
# Shorter responses (faster)
export VIBE_MAX_TOKENS=300

# Default
export VIBE_MAX_TOKENS=500

# Longer responses (for complex commands)
export VIBE_MAX_TOKENS=1000
```

**Note:** Higher values increase API costs and response time.

---

#### VIBE_TIMEOUT

**Type:** Duration  
**Default:** `30s`  
**Description:** Maximum time to wait for API response.

**Examples:**

```bash 
# Shorter timeout
export VIBE_TIMEOUT=15s

# Default
export VIBE_TIMEOUT=30s

# Longer timeout (slow connections)
export VIBE_TIMEOUT=60s
```

**Format:** Number + unit (`s` for seconds, `m` for minutes)

---

### Display Configuration

#### VIBE_SHOW_EXPLANATION

**Type:** Boolean  
**Default:** `true`  
**Description:** Show inline explanations for generated commands.

**Examples:**

```bash 
# Show explanations (default)
export VIBE_SHOW_EXPLANATION=true

# Hide explanations
export VIBE_SHOW_EXPLANATION=false
```

**When enabled:**
```bash
docker ps -a
# docker: Docker command-line tool
# ps: List containers
# -a: Show all containers (not just running)
```

**When disabled:**
```bash
docker ps -a
```

---

#### VIBE_SHOW_WARNINGS

**Type:** Boolean  
**Default:** `true`  
**Description:** Display warnings for potentially dangerous commands.

**Examples:**

```bash 
# Show warnings (default)
export VIBE_SHOW_WARNINGS=true

# Hide warnings
export VIBE_SHOW_WARNINGS=false
```

**When enabled:**
```bash
⚠️  WARNING: This command may delete files
rm -rf *
```

---

### Behavior Configuration

#### VIBE_INTERACTIVE

**Type:** Boolean  
**Default:** `false`  
**Description:** Require confirmation before inserting commands into the prompt.

**Examples:**

```bash 
# Auto-insert (default)
export VIBE_INTERACTIVE=false

# Require confirmation
export VIBE_INTERACTIVE=true
```

**When enabled:**
```
---
docker ps -a
---
Execute this command? [Y/n]
```

---

### Cache Configuration

#### VIBE_ENABLE_CACHE

**Type:** Boolean  
**Default:** `true`  
**Description:** Enable response caching for faster repeated queries.

**Examples:**

```bash 
# Enable cache (default)
export VIBE_ENABLE_CACHE=true

# Disable cache
export VIBE_ENABLE_CACHE=false
```

**Impact:**
- **Enabled:** Cached queries return in ~5-10ms (100-400x faster)
- **Disabled:** Every query hits the API

---

#### VIBE_CACHE_TTL

**Type:** Duration  
**Default:** `24h`  
**Description:** How long cached responses remain valid.

**Examples:**

```bash 
# 1 hour
export VIBE_CACHE_TTL=1h

# 12 hours
export VIBE_CACHE_TTL=12h

# 24 hours (default)
export VIBE_CACHE_TTL=24h

# 7 days
export VIBE_CACHE_TTL=168h
```

**Format:** Number + unit (`h` for hours, `m` for minutes)

**Cache location:** `~/.cache/vibe/`

---

## Configuration Templates

### Local LLM (Ollama)

```bash 
# Minimal configuration (uses defaults)
export VIBE_API_URL="http://localhost:11434/v1"
export VIBE_MODEL="llama3:8b"
```

### OpenAI

```bash 
export VIBE_API_URL="https://api.openai.com/v1"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4"
export VIBE_TEMPERATURE=0.5
export VIBE_TIMEOUT=30s
```

### Anthropic (via OpenRouter)

```bash 
export VIBE_API_URL="https://openrouter.ai/api/v1"
export VIBE_API_KEY="sk-or-..."
export VIBE_MODEL="anthropic/claude-3.5-sonnet"
export VIBE_TEMPERATURE=0.7
export VIBE_MAX_TOKENS=500
```

### Groq

```bash 
export VIBE_API_URL="https://api.groq.com/openai/v1"
export VIBE_API_KEY="gsk_..."
export VIBE_MODEL="llama-3.1-70b-versatile"
export VIBE_TEMPERATURE=0.4
```

### Complete Configuration Example

```bash 
# API Configuration
export VIBE_API_URL="https://api.openai.com/v1"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4"
export VIBE_TEMPERATURE=0.5
export VIBE_MAX_TOKENS=500
export VIBE_TIMEOUT=30s

# Display Configuration
export VIBE_SHOW_EXPLANATION=true
export VIBE_SHOW_WARNINGS=true

# Behavior Configuration
export VIBE_INTERACTIVE=false

# Cache Configuration
export VIBE_ENABLE_CACHE=true
export VIBE_CACHE_TTL=24h
```

---

## Plugin Configuration

### VIBE_PLUGIN_DIR

**Type:** String (read-only)  
**Default:** Auto-detected  
**Description:** Directory where the vibe plugin is installed. Automatically set by the plugin.

```bash
# Typically:
~/.oh-my-zsh/custom/plugins/vibe
```

---

### VIBE_BINARY

**Type:** String  
**Default:** `${VIBE_PLUGIN_DIR}/vibe`  
**Description:** Path to the vibe binary. Override if you want to use a custom binary location.

**Example:**

```bash 
export VIBE_BINARY="/usr/local/bin/vibe"
```

---

## Quick Reference Table

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `VIBE_API_URL` | String | `http://localhost:11434/v1` | API endpoint URL |
| `VIBE_API_KEY` | String | `""` | API authentication key |
| `VIBE_MODEL` | String | `llama3:8b` | Model identifier |
| `VIBE_TEMPERATURE` | Float | `0.7` | Randomness (0.0-2.0) |
| `VIBE_MAX_TOKENS` | Integer | `500` | Max response tokens |
| `VIBE_TIMEOUT` | Duration | `30s` | Request timeout |
| `VIBE_SHOW_EXPLANATION` | Boolean | `true` | Show explanations |
| `VIBE_SHOW_WARNINGS` | Boolean | `true` | Show warnings |
| `VIBE_INTERACTIVE` | Boolean | `false` | Confirm before insert |
| `VIBE_ENABLE_CACHE` | Boolean | `true` | Enable caching |
| `VIBE_CACHE_TTL` | Duration | `24h` | Cache lifetime |
| `VIBE_BINARY` | String | Auto | Binary path |

---

## Applying Configuration Changes

After modifying environment variables in `~/.zshrc`:

```bash 
source ~/.zshrc
```

Or restart your terminal.

---

## Configuration Files

vibe uses **environment variables only** - no configuration files are required or used.

**Cache location:**
```
~/.cache/vibe/
```

**Plugin location:**
```
~/.oh-my-zsh/custom/plugins/vibe/
```

To reset cache:

```bash 
rm -rf ~/.cache/vibe/*
```

---

## Troubleshooting Configuration

### View Current Configuration

Check what vibe sees:

```bash 
env | grep VIBE_
```

### Test Configuration

Test with a simple query:

```bash 
~/.oh-my-zsh/custom/plugins/vibe/vibe "list files"
```

### Common Issues

**"Connection refused"**
- Check `VIBE_API_URL` is correct
- Verify service is running (e.g., `ollama serve`)

**"Unauthorized" or "Invalid API key"**
- Verify `VIBE_API_KEY` is set and correct
- Check for extra quotes or spaces

**"Model not found"**
- Verify `VIBE_MODEL` exists for your provider
- For Ollama: `ollama list`

**Slow responses**
- Enable cache: `VIBE_ENABLE_CACHE=true`
- Reduce `VIBE_MAX_TOKENS`
- Use faster model or local LLM
