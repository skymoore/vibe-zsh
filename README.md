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
- ⚡ **Lightning fast** - Cached responses are 100-400x faster
- 🔌 **Universal compatibility** - Works with Ollama, OpenAI, Claude, and more
- 🛡️ **Safe by default** - Preview commands before execution
- 📚 **Learn while you work** - Inline explanations for every command
- 🎯 **Zero dependencies** - Single compiled binary

## Installation

### Homebrew (Recommended)

```bash
brew tap skymoore/tap
brew install vibe
```

Add to your `~/.zshrc`:
```bash
source $(brew --prefix)/opt/vibe/libexec/vibe.plugin.zsh
```

Then reload your shell:
```bash
source ~/.zshrc
```

### Oh-My-Zsh Plugin

For oh-my-zsh integration:

```bash
curl -fsSL https://raw.githubusercontent.com/skymoore/vibe-zsh/main/install.sh | bash
```

Or using wget:

```bash
wget -qO- https://raw.githubusercontent.com/skymoore/vibe-zsh/main/install.sh | bash
```

This installs vibe as an oh-my-zsh plugin and automatically adds it to your plugins list.

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

### Basic Usage

1. Type a natural language description in your terminal
2. Press `Ctrl+G`
3. Review the generated command (with explanations)
4. Press `Enter` to execute, or edit first

### Examples

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
| `VIBE_AUTO_UPDATE` | `true` | Enable auto-update checks |
| `VIBE_UPDATE_CHECK_INTERVAL` | `7d` | How often to check for updates |

## How It Works

1. **Capture** - You type natural language and press `Ctrl+G`
2. **Generate** - vibe sends your query to the configured LLM
3. **Parse** - Response is structured as command + explanations
4. **Cache** - Response is cached for 24 hours (configurable)
5. **Insert** - Command replaces your input for review
6. **Execute** - You press Enter to run (or edit first)

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
brew uninstall vibe
brew untap skymoore/tap  # Optional: remove the tap
```

Then remove the `source` line from your `~/.zshrc`.

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
