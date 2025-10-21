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
**Default:** `0.2`  
**Range:** `0.0` to `2.0`  
**Description:** Controls randomness in the model's output. Lower values make output more deterministic, higher values more creative.

**Examples:**

```bash 
# Most deterministic (default, recommended for commands)
export VIBE_TEMPERATURE=0.2

# More creative
export VIBE_TEMPERATURE=0.5

# Higher variation
export VIBE_TEMPERATURE=0.7
```

**Recommendations:**
- `0.0-0.3`: Most consistent, predictable commands (recommended)
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

### Parsing & Reliability Configuration

#### VIBE_USE_STRUCTURED_OUTPUT

**Type:** Boolean  
**Default:** `true`  
**Description:** Use JSON schema for structured responses. Provides better reliability with models that support structured output.

**Examples:**

```bash 
# Use structured output (default, recommended)
export VIBE_USE_STRUCTURED_OUTPUT=true

# Disable structured output (fallback to text parsing)
export VIBE_USE_STRUCTURED_OUTPUT=false
```

**Recommendations:**
- Keep enabled for OpenAI, Claude, and modern models
- Disable if your model doesn't support JSON schema

---

#### VIBE_MAX_RETRIES

**Type:** Integer  
**Default:** `3`  
**Description:** Maximum retry attempts for failed parsing. Higher values increase reliability but may slow down error cases.

**Examples:**

```bash 
# Conservative (faster failures)
export VIBE_MAX_RETRIES=2

# Default (recommended)
export VIBE_MAX_RETRIES=3

# Aggressive (maximum reliability)
export VIBE_MAX_RETRIES=5
```

**When to adjust:**
- Increase if you see occasional parsing failures
- Decrease if you want faster error responses

---

#### VIBE_ENABLE_JSON_EXTRACTION

**Type:** Boolean  
**Default:** `true`  
**Description:** Automatically extract valid JSON from corrupted LLM responses containing garbage characters, escape sequences, or Unicode pollution.

**Examples:**

```bash 
# Enable JSON extraction (default, recommended)
export VIBE_ENABLE_JSON_EXTRACTION=true

# Disable extraction (strict parsing only)
export VIBE_ENABLE_JSON_EXTRACTION=false
```

**What it handles:**
- Terminal escape sequences (`\x1b[200~`, `\x1b[201~`)
- Unicode separators (`\u2028`, `\u2029`)
- Control characters and garbage text
- JSON buried in explanatory text

---

#### VIBE_STRICT_VALIDATION

**Type:** Boolean  
**Default:** `true`  
**Description:** Validate response structure for required fields and proper formatting.

**Examples:**

```bash 
# Enable validation (default, recommended)
export VIBE_STRICT_VALIDATION=true

# Disable validation (permissive)
export VIBE_STRICT_VALIDATION=false
```

**What it checks:**
- Command field is non-empty
- Explanation array has entries
- No whitespace-only fields

**When to disable:**
- Debugging parsing issues
- Testing with experimental models

---

#### VIBE_DEBUG_LOGS

**Type:** Boolean  
**Default:** `false`  
**Description:** Enable detailed debug logging for troubleshooting parsing and API issues.

**Examples:**

```bash 
# Disable debug logs (default)
export VIBE_DEBUG_LOGS=false

# Enable debug logs
export VIBE_DEBUG_LOGS=true
```

**What it logs:**
- Raw LLM responses (first 500 chars)
- Parsing attempts and failures
- Which fallback layer succeeded
- JSON extraction details
- Garbage removed from responses

**Output example:**
```
[VIBE] [ATTEMPT 1][structured_output] Parsing failed: invalid JSON
[VIBE] [ATTEMPT 2][enhanced_parsing] Parsing failed: ...
[VIBE] [JSON_EXTRACT] Success
Trimmed prefix: "AWS\x1b[200~..."
[VIBE] [SUCCESS] Layer 'enhanced_parsing' succeeded on attempt 2
```

---

#### VIBE_SHOW_RETRY_STATUS

**Type:** Boolean  
**Default:** `true`  
**Description:** Show progress indicator during retry attempts.

**Examples:**

```bash 
# Show retry status (default)
export VIBE_SHOW_RETRY_STATUS=true

# Hide retry status
export VIBE_SHOW_RETRY_STATUS=false
```

**When enabled:**
```
⏳ Generating command (attempt 2/3)...
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
export VIBE_TEMPERATURE=0.2
export VIBE_TIMEOUT=30s
```

### Anthropic (via OpenRouter)

```bash 
export VIBE_API_URL="https://openrouter.ai/api/v1"
export VIBE_API_KEY="sk-or-..."
export VIBE_MODEL="anthropic/claude-3.5-sonnet"
export VIBE_TEMPERATURE=0.2
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
export VIBE_TEMPERATURE=0.2
export VIBE_MAX_TOKENS=500
export VIBE_TIMEOUT=30s

# Parsing & Reliability (defaults shown, usually don't need to change)
export VIBE_USE_STRUCTURED_OUTPUT=true
export VIBE_MAX_RETRIES=3
export VIBE_ENABLE_JSON_EXTRACTION=true
export VIBE_STRICT_VALIDATION=true
export VIBE_DEBUG_LOGS=false
export VIBE_SHOW_RETRY_STATUS=true

# Display Configuration
export VIBE_SHOW_EXPLANATION=true
export VIBE_SHOW_WARNINGS=true

# Behavior Configuration
export VIBE_INTERACTIVE=false

# Cache Configuration
export VIBE_ENABLE_CACHE=true
export VIBE_CACHE_TTL=24h
```

### Debugging Configuration

When troubleshooting issues, use this minimal debug config:

```bash 
# Enable debug logging
export VIBE_DEBUG_LOGS=true

# Increase retries to see all fallback attempts
export VIBE_MAX_RETRIES=5

# Keep other settings at defaults
export VIBE_ENABLE_JSON_EXTRACTION=true
export VIBE_STRICT_VALIDATION=true
export VIBE_SHOW_RETRY_STATUS=true
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

### Core Settings

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `VIBE_API_URL` | String | `http://localhost:11434/v1` | API endpoint URL |
| `VIBE_API_KEY` | String | `""` | API authentication key |
| `VIBE_MODEL` | String | `llama3:8b` | Model identifier |
| `VIBE_TEMPERATURE` | Float | `0.2` | Randomness (0.0-2.0) |
| `VIBE_MAX_TOKENS` | Integer | `500` | Max response tokens |
| `VIBE_TIMEOUT` | Duration | `30s` | Request timeout |

### Parsing & Reliability

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `VIBE_USE_STRUCTURED_OUTPUT` | Boolean | `true` | Use JSON schema for responses |
| `VIBE_MAX_RETRIES` | Integer | `3` | Max retry attempts |
| `VIBE_ENABLE_JSON_EXTRACTION` | Boolean | `true` | Extract JSON from corrupted responses |
| `VIBE_STRICT_VALIDATION` | Boolean | `true` | Validate response structure |
| `VIBE_DEBUG_LOGS` | Boolean | `false` | Enable debug logging |
| `VIBE_SHOW_RETRY_STATUS` | Boolean | `true` | Show retry progress |

### Display & Behavior

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `VIBE_SHOW_EXPLANATION` | Boolean | `true` | Show explanations |
| `VIBE_SHOW_WARNINGS` | Boolean | `true` | Show warnings |
| `VIBE_INTERACTIVE` | Boolean | `false` | Confirm before insert |

### Cache Settings

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `VIBE_ENABLE_CACHE` | Boolean | `true` | Enable caching |
| `VIBE_CACHE_TTL` | Duration | `24h` | Cache lifetime |

### System

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
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

**Parsing failures or corrupted output**
- Enable debug logs: `VIBE_DEBUG_LOGS=true`
- Check logs for which layer is failing
- Ensure JSON extraction is enabled: `VIBE_ENABLE_JSON_EXTRACTION=true`
- Try lowering temperature: `VIBE_TEMPERATURE=0.3`
- Increase retries: `VIBE_MAX_RETRIES=5`
