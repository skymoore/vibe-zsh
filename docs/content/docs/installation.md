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

The Homebrew installation automatically:
1. Installs vibe to `~/.oh-my-zsh/custom/plugins/vibe`
2. Adds `vibe` to your plugins list in `~/.zshrc`
3. Creates a global `vibe-zsh` command for CLI usage

After installation, reload your shell:

```bash
source ~/.zshrc
```

**Note:** Homebrew installation requires Oh-My-Zsh. If you don't have it installed, the installer will provide instructions.

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
1. Download the latest release for your platform (or clone and build from source if no release is available)
2. Install to `~/.oh-my-zsh/custom/plugins/vibe`
3. Display instructions for adding `vibe` to your plugins list in `.zshrc`

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
