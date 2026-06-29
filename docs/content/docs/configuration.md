---
title: "Configuration"
weight: 2
---


All configuration is done through environment variables in your `~/.zshrc`. All settings are optional and have sensible defaults.

## Provider Configuration

vibe-zsh uses [gollm](https://github.com/teilomillet/gollm) to talk to each
provider natively. There are two kinds of provider:

- **Hosted providers** (`openai`, `anthropic`, `groq`, `openrouter`, `deepseek`,
  `google-openai`, `mistral`, `cohere`) have a fixed endpoint built in. Set
  `VIBE_PROVIDER`, `VIBE_API_KEY`, and `VIBE_MODEL` — do **not** set `VIBE_API_URL`.
- **Local providers** (`ollama`, `lmstudio`, `vllm`) run on your machine. Set
  `VIBE_PROVIDER` and point `VIBE_API_URL` at the local server; no API key needed.

Select a provider with `VIBE_PROVIDER`. If you don't set it, vibe infers the
provider from `VIBE_API_URL` (mainly to keep older configs working). **Setting
`VIBE_PROVIDER` explicitly is the recommended approach.**

{{< callout type="warning" >}}
The default `VIBE_MODEL` is `llama3:8b` (chosen for the default Ollama setup).
Always set `VIBE_MODEL` when you use a hosted provider, or requests will ask for
a model that doesn't exist there.
{{< /callout >}}

### Ollama (local) - Default

```bash
export VIBE_PROVIDER="ollama"
export VIBE_API_URL="http://localhost:11434/v1"
export VIBE_MODEL="llama3:8b"
```

This is the default configuration. If you have Ollama running locally, vibe will work out of the box.

### OpenAI

```bash
export VIBE_PROVIDER="openai"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4o"
# No VIBE_API_URL needed — the openai provider always targets api.openai.com.
```

### Anthropic

```bash
export VIBE_PROVIDER="anthropic"
export VIBE_API_KEY="sk-ant-..."
export VIBE_MODEL="claude-3-5-sonnet-20241022"
```

### Groq

```bash
export VIBE_PROVIDER="groq"
export VIBE_API_KEY="gsk_..."
export VIBE_MODEL="llama-3.1-70b-versatile"
```

### OpenRouter

```bash
export VIBE_PROVIDER="openrouter"
export VIBE_API_KEY="sk-or-..."
export VIBE_MODEL="anthropic/claude-3.5-sonnet"
```

### LM Studio (local)

```bash
export VIBE_PROVIDER="lmstudio"
export VIBE_API_URL="http://localhost:1234/v1"
export VIBE_MODEL="local-model"
```

{{< callout type="info" >}}
**Validation:** gollm checks your configuration when vibe starts. Hosted
providers validate the API key format up front (e.g. Anthropic keys must start
with `sk-ant-`), and local providers must already be running and reachable.

**Custom OpenAI-compatible gateways:** the `openai` provider always points at
`api.openai.com` and cannot be redirected via `VIBE_API_URL`. To use an
OpenAI-compatible endpoint that isn't OpenAI itself, set `VIBE_PROVIDER` to
`lmstudio` or `vllm` and point `VIBE_API_URL` at your gateway.
{{< /callout >}}

## Display Configuration

### Progress Indicators

Control the progress spinner shown during command generation (default: `true`):

```bash
export VIBE_SHOW_PROGRESS=false
```

Choose spinner style (default: `dots`):

```bash
export VIBE_PROGRESS_STYLE=arrow  # Options: dots, line, circle, bounce, arrow, runes
```

**Available Styles:**
- `dots` - ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
- `line` - - \ | /
- `circle` - ◐ ◓ ◑ ◒
- `bounce` - ⠁ ⠂ ⠄ ⠂
- `arrow` - ← ↖ ↑ ↗ → ↘ ↓ ↙
- `runes` - ᛜ ᛃ ᛋ (Norse runes)

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

### Query History

Control query history tracking (default: `true`):

```bash
export VIBE_ENABLE_HISTORY=false  # Disable history
```

Set maximum number of history entries (default: `100`):

```bash
export VIBE_HISTORY_SIZE=200
```

Customize the history keybindings:

```bash
export VIBE_HISTORY_KEY="^R"      # Use Ctrl+R for history menu (default: ^Xh)
export VIBE_REGENERATE_KEY="^[r"  # Use Alt+R for quick regenerate (default: ^Xg)
```

**Note:** Avoid using `^H` (Ctrl+H) as it conflicts with Backspace.

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

Set generation temperature (default: `0.2`):

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
| `VIBE_PROVIDER` | _(inferred from `VIBE_API_URL`)_ | LLM provider. Hosted: `openai`, `anthropic`, `groq`, `openrouter`, `deepseek`, `google-openai`, `mistral`, `cohere`. Local: `ollama`, `lmstudio`, `vllm`. Recommended to set explicitly. |
| `VIBE_API_URL` | `http://localhost:11434/v1` | Endpoint URL for **local** providers only. Hosted providers ignore this. |
| `VIBE_API_KEY` | `""` | API key. Required for hosted providers; ignored by local providers. |
| `VIBE_MODEL` | `llama3:8b` | Model to use. Set this for hosted providers — the default only suits Ollama. |
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
| **History** | | |
| `VIBE_ENABLE_HISTORY` | `true` | Enable query history tracking |
| `VIBE_HISTORY_SIZE` | `100` | Maximum number of history entries |
| `VIBE_HISTORY_KEY` | `^Xh` (Ctrl+X H) | Keybinding for history menu |
| `VIBE_REGENERATE_KEY` | `^Xg` (Ctrl+X G) | Keybinding to regenerate last command |
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
export VIBE_PROVIDER="openai"
export VIBE_API_KEY="sk-..."
export VIBE_MODEL="gpt-4o"
export VIBE_TEMPERATURE=0.2
export VIBE_SHOW_EXPLANATION=true
export VIBE_SHOW_PROGRESS=true
export VIBE_PROGRESS_STYLE=dots
export VIBE_STREAM_OUTPUT=true
export VIBE_STREAM_DELAY=20ms
export VIBE_INTERACTIVE=false
export VIBE_ENABLE_CACHE=true
export VIBE_CACHE_TTL=24h
export VIBE_ENABLE_HISTORY=true
export VIBE_HISTORY_SIZE=100
export VIBE_HISTORY_KEY="^Xh"
export VIBE_REGENERATE_KEY="^Xg"
```
