# vibe 🌊

Transform natural language into shell commands using AI. Works with any OpenAI-compatible API.

```bash
# Type your intent in natural language
list all docker containers

# Press Ctrl+G

# Get the command with explanations
docker ps -a
# docker: Docker command-line tool
# ps: List containers
# -a: Show all containers (not just running)
```

## Features

- 🧠 **Natural language to commands** - Just describe what you want
- 🖥️ **OS-aware generation** - Commands that work on YOUR system (macOS/Linux/Windows)
- ⚡ **Lightning fast** - Cached responses are 100-400x faster
- 🎬 **Streaming output** - Typewriter effect with progress indicators
- 🔌 **Universal compatibility** - Works with Ollama, LM Studio, OpenAI, Claude, and more
- 🛡️ **Safe by default** - Preview commands before execution
- 📚 **Learn while you work** - Inline explanations for every command
- 📜 **Query history** - Interactive menu to browse and re-run previous queries
- 🎯 **Zero dependencies** - Single compiled binary

## Installation

### Homebrew (Recommended)

```bash
brew tap skymoore/tap
brew install --cask vibe-zsh
```

The Homebrew installation automatically:
- Installs vibe to `~/.oh-my-zsh/custom/plugins/vibe`
- Adds `vibe` to your plugins list in `~/.zshrc`
- Creates a global `vibe-zsh` command for CLI usage

After installation, reload your shell:
```bash
source ~/.zshrc
```

**Note:** Requires Oh-My-Zsh. If not installed, the installer will provide instructions.

### Oh-My-Zsh Plugin

For oh-my-zsh integration:

```bash
curl -fsSL https://raw.githubusercontent.com/skymoore/vibe-zsh/main/install.sh | bash
```

Or using wget:

```bash
wget -qO- https://raw.githubusercontent.com/skymoore/vibe-zsh/main/install.sh | bash
```

This script downloads the latest release and installs it to `~/.oh-my-zsh/custom/plugins/vibe`. You'll need to manually add `vibe` to your plugins list in `~/.zshrc`.

### Manual Install

1. **Clone the repository:**
   ```bash
   git clone https://github.com/skymoore/vibe-zsh.git ~/.oh-my-zsh/custom/plugins/vibe
   ```

2. **Build the binary:**
   ```bash
   cd ~/.oh-my-zsh/custom/plugins/vibe
   make build
   ```

3. **Add to your `.zshrc`:**
   ```bash
   plugins=(... vibe)
   ```

4. **Reload your shell:**
   ```bash
   source ~/.zshrc
   ```

## Configuration

Add these to your `~/.zshrc` (all optional):

### Local LLM (Ollama)
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

### Anthropic (via OpenRouter)
```bash
export VIBE_API_URL="https://openrouter.ai/api/v1"
export VIBE_API_KEY="sk-or-..."
export VIBE_MODEL="anthropic/claude-3.5-sonnet"
```

### Other Providers

**LM Studio:**
```bash
export VIBE_API_URL="http://localhost:1234/v1"
export VIBE_MODEL="local-model"
```

**Groq:**
```bash
export VIBE_API_URL="https://api.groq.com/openai/v1"
export VIBE_API_KEY="gsk_..."
export VIBE_MODEL="llama-3.1-70b-versatile"
```

## Usage

### Plugin Usage (Recommended)

1. Type a natural language description in your terminal
2. Press `Ctrl+G`
3. Review the generated command (with explanations)
4. Press `Enter` to execute, or edit first

**History Menu:**

Access your query history in two ways:

1. **Keybinding**: Press `Ctrl+X` then `H`
2. **Command**: Type `vh` and press Enter

Both methods open an interactive menu where you can:
- Browse previous queries with arrow keys
- Search with `/` (filter mode)
- Press `Enter` to insert the selected command into your buffer
- Press `Esc` or `q` to cancel

The selected command appears on your command line, ready to execute!

**Note**: Don't run `vibe-zsh history` or `./vibe history` directly - use `vh` or the keybinding instead.

### Direct CLI Usage

You can also use vibe directly from the command line:

```bash
# Homebrew installation
vibe-zsh "list all docker containers"

# Manual installation (Oh-My-Zsh plugin)
~/.oh-my-zsh/custom/plugins/vibe/vibe "list all docker containers"
```

**CLI Flags:**
```bash
vibe-zsh --help                    # Show all available flags
vibe-zsh --debug "query"           # Enable debug logging
vibe-zsh --temperature 0.1 "query" # Override temperature
vibe-zsh --interactive "query"     # Confirm before execution
```

**History Commands:**
```bash
vh                                 # Interactive menu (recommended)
vibe-zsh history list              # List history in plain text
vibe-zsh history clear             # Clear all history

# Note: Use 'vh' or Ctrl+X H for interactive menu
# Don't use 'vibe-zsh history' directly (it just outputs to stdout)
```

### Examples

**Query History:**
```bash
# Open history menu with keybinding
Ctrl+X H

# Or use the command
vh

# List history in plain text
vibe-zsh history list

# Clear all history
vibe-zsh history clear
```

**File Operations:**
```bash
show me all hidden files including their sizes
# → ls -lah
```

**Docker:**
```bash
show logs of nginx container and follow them
# → docker logs -f nginx
```

**Git:**
```bash
show me commits from last week
# → git log --since="1 week ago"
```

**Find & Search:**
```bash
find all python files modified today
# → find . -name "*.py" -mtime 0
```

### Advanced Features

**Tab Completion:**
```bash
vibe <TAB><TAB>
# Shows common query suggestions
```

**Hide Explanations:**
```bash
export VIBE_SHOW_EXPLANATION=false
```

**Interactive Mode** (confirm before inserting):
```bash
export VIBE_INTERACTIVE=true
```

**Disable Cache:**
```bash
export VIBE_ENABLE_CACHE=false
```

**Disable Auto-Updates:**
```bash
export VIBE_AUTO_UPDATE=false
```

**Customize History Keybinding:**
```bash
export VIBE_HISTORY_KEY="^R"   # Use Ctrl+R instead
# Note: Avoid ^H (Ctrl+H) as it conflicts with Backspace
```

**Disable History:**
```bash
export VIBE_ENABLE_HISTORY=false
```

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
| **History** | | |
| `VIBE_ENABLE_HISTORY` | `true` | Enable query history tracking |
| `VIBE_HISTORY_SIZE` | `100` | Maximum number of history entries |
| `VIBE_HISTORY_KEY` | `^Xh` (Ctrl+X H) | Keybinding for history menu |
| **Updates & Debugging** | | |
| `VIBE_AUTO_UPDATE` | `true` | Enable auto-update checks |
| `VIBE_UPDATE_CHECK_INTERVAL` | `7d` | How often to check for updates |
| `VIBE_DEBUG_LOGS` | `false` | Enable debug logging for troubleshooting |

## How It Works

1. **Capture** - You type natural language and press `Ctrl+G`
2. **Context** - vibe detects your OS (macOS/Linux/Windows) and shell (zsh/bash/etc.)
3. **Generate** - vibe sends your query with system context to the configured LLM
4. **Parse** - Response is structured as command + explanations
5. **Cache** - Response is cached for 24 hours (configurable)
6. **Stream** - Command is streamed to your terminal with typewriter effect
7. **Insert** - Command appears in your buffer for review
8. **Execute** - You press Enter to run (or edit first)

### OS-Aware Command Generation

vibe automatically detects your operating system and shell, ensuring generated commands work on your system:

- **macOS**: Uses BSD utilities (e.g., `find` without `-printf`, `sed -i ''`)
- **Linux**: Uses GNU utilities (e.g., `find -printf`, `sed -i`)
- **Shell-specific**: Generates syntax appropriate for zsh, bash, etc.

This means you get commands that actually work on your system, not generic Linux commands that fail on macOS!

## Performance

- **First query:** ~500ms-2s (depends on LLM)
- **Cached query:** ~5-10ms (100-400x faster!)
- **Binary size:** ~8MB
- **Memory usage:** <10MB

## Safety

- ✅ Commands are never executed automatically
- ✅ Full preview before execution
- ✅ Edit commands before running
- ✅ Warnings for dangerous commands
- ✅ Optional interactive confirmation mode
- ✅ Local-first with Ollama (your data stays private)

## Updates

vibe automatically checks for updates once a week in the background (zero impact on performance). When an update is available, you'll see a notification:

```
⚠️  vibe v1.2.4 available (current: v1.2.3)
   Run: vibe --update
```

**Manual Commands:**
```bash
vibe --version              # Show current version
vibe --update               # Download and install latest version
```

**Configuration:**
```bash
export VIBE_AUTO_UPDATE=false           # Disable auto-update checks
export VIBE_UPDATE_CHECK_INTERVAL=14d   # Check every 2 weeks instead
```

**Security:**
- All downloads are verified using SHA256 checksums
- Checksums are published with each release
- Original binary is backed up before replacement
- Safe atomic file replacement

The update check runs in a background process after each command, so it never slows down your workflow. Updates are never installed automatically - you always control when to update.

## Troubleshooting

**Command not working after install:**
- Reload your shell: `source ~/.zshrc`
- Verify plugin is in list: `echo $plugins`
- Check binary exists: `ls ~/.oh-my-zsh/custom/plugins/vibe/vibe`

**Slow responses:**
- Enable caching: `export VIBE_ENABLE_CACHE=true`
- Use a faster model or local LLM
- Check your network connection

**Bad command suggestions:**
- Try a different/better model
- Make your query more specific
- Check model supports structured output

**Corrupted or garbage output:**
- vibe now automatically handles corrupted LLM responses with multi-layer parsing
- Enable debug logs to see what's happening: `export VIBE_DEBUG_LOGS=true`
- Check logs for raw responses and parsing attempts
- Try increasing retries: `export VIBE_MAX_RETRIES=5`

**Parsing errors:**
- Enable JSON extraction: `export VIBE_ENABLE_JSON_EXTRACTION=true` (default)
- Disable strict validation temporarily: `export VIBE_STRICT_VALIDATION=false`
- Check debug logs to see which parsing layer succeeded/failed
- Some models produce cleaner JSON with lower temperature: `export VIBE_TEMPERATURE=0.3`

**Ctrl+G does nothing:**
- Ensure plugin is loaded: `which vibe` (should show a function)
- Check no other plugin uses Ctrl+G
- Try rebinding: `bindkey '^G' vibe`

## Requirements

- Zsh with Oh-My-Zsh
- OpenAI-compatible API endpoint (Ollama, OpenAI, etc.)
- macOS or Linux
- No external dependencies!

## Uninstalling

### Homebrew

```bash
brew uninstall --cask vibe-zsh
brew untap skymoore/tap  # Optional: remove the tap
```

Then remove `vibe` from your plugins list in `~/.zshrc` and reload:
```bash
source ~/.zshrc
```

### Oh-My-Zsh Plugin

```bash
curl -fsSL https://raw.githubusercontent.com/skymoore/vibe-zsh/main/uninstall.sh | bash
```

Or using wget:

```bash
wget -qO- https://raw.githubusercontent.com/skymoore/vibe-zsh/main/uninstall.sh | bash
```

Then remove `vibe` from your plugins list in `~/.zshrc` and reload:
```bash
source ~/.zshrc
```

## Building from Source

```bash
git clone https://github.com/skymoore/vibe-zsh.git
cd vibe-zsh
make build        # Build for current platform
make build-all    # Build for all platforms
make test         # Run tests
make install      # Install to oh-my-zsh
```

## Contributing

Contributions welcome! This project uses:
- Go 1.21+ (stdlib only, no external deps)
- Standard zsh scripting
- OpenAI-compatible API spec

## License

GPLv3 License - see [LICENSE](LICENSE)

## Acknowledgments

Inspired by [LoganPederson/vibe](https://github.com/LoganPederson/vibe)

Built for terminal productivity enthusiasts 🚀
