---
title: "Installation"
weight: 1
---

## Homebrew (Recommended)

The easiest way to install vibe on macOS and Linux is via Homebrew:

```bash
brew tap skymoore/tap
brew install --cask vibe-zsh
```

Add the following line to your `~/.zshrc`:

```bash
source "$(brew --prefix)/opt/vibe/libexec/vibe.plugin.zsh"
```

Then reload your shell:

```bash
source ~/.zshrc
```

Zsh completions are installed automatically.

## Oh-My-Zsh Plugin

If you use Oh-My-Zsh and prefer plugin-based installation, use our automated installation script:

{{< tabs >}}

{{< tab name="curl" >}}
```bash
curl -fsSL https://raw.githubusercontent.com/skymoore/vibe-zsh/main/install.sh | bash
```
{{< /tab >}}

{{< tab name="wget" >}}
```bash
wget -qO- https://raw.githubusercontent.com/skymoore/vibe-zsh/main/install.sh | bash
```
{{< /tab >}}

{{< /tabs >}}

This script will:
1. Clone the repository to `~/.oh-my-zsh/custom/plugins/vibe`
2. Build the binary for your platform
3. Add the plugin to your `.zshrc`
4. Reload your shell

## Manual Install

If you prefer to install manually:

### 1. Clone the repository

```bash
git clone https://github.com/skymoore/vibe-zsh.git ~/.oh-my-zsh/custom/plugins/vibe
```

### 2. Build the binary

```bash
cd ~/.oh-my-zsh/custom/plugins/vibe
make build
```

### 3. Add to your `.zshrc`

Edit your `~/.zshrc` and add `vibe` to the plugins array:

```bash
plugins=(... vibe)
```

### 4. Reload your shell

```bash
source ~/.zshrc
```

## Building from Source

If you want to build from source:

```bash
git clone https://github.com/skymoore/vibe-zsh.git
cd vibe-zsh
make build        # Build for current platform
make build-all    # Build for all platforms
make test         # Run tests
make install      # Install to oh-my-zsh
```

## Verification

After installation, verify that vibe is working:

```bash
# Check that the binary exists
which vibe

# Type something and press Ctrl+G
echo "test"
# Then press Ctrl+G
```

## Uninstalling

### Homebrew

To uninstall vibe installed via Homebrew:

```bash
brew uninstall --cask vibe-zsh
brew untap skymoore/tap  # Optional: remove the tap
```

Then remove the `source` line from your `~/.zshrc` and reload your shell:

```bash
source ~/.zshrc
```

### Oh-My-Zsh Plugin

To uninstall vibe installed as an Oh-My-Zsh plugin:

{{< tabs >}}

{{< tab name="curl" >}}
```bash
curl -fsSL https://raw.githubusercontent.com/skymoore/vibe-zsh/main/uninstall.sh | bash
```
{{< /tab >}}

{{< tab name="wget" >}}
```bash
wget -qO- https://raw.githubusercontent.com/skymoore/vibe-zsh/main/uninstall.sh | bash
```
{{< /tab >}}

{{< /tabs >}}

Then remove `vibe` from your plugins list in `~/.zshrc` and reload:

```bash
source ~/.zshrc
```

### Manual Uninstall

If you installed manually, simply remove the plugin directory:

```bash
rm -rf ~/.oh-my-zsh/custom/plugins/vibe
```

Then remove `vibe` from your plugins list in `~/.zshrc` and reload your shell.

## Next Steps

After installation, you'll want to [configure your API settings]({{< relref "configuration" >}}).
