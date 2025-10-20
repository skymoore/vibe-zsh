---
title: "Installation"
weight: 1
---

## One-Line Install

The easiest way to install vibe is using our automated installation script:

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

## Next Steps

After installation, you'll want to [configure your API settings]({{< relref "configuration" >}}).
