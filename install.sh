#!/usr/bin/env bash

set -e

REPO="skymoore/vibe-zsh"
INSTALL_DIR="${HOME}/.oh-my-zsh/custom/plugins/vibe"
TMP_DIR=$(mktemp -d)

echo "🚀 Installing vibe-zsh..."

if ! command -v zsh &> /dev/null; then
    echo "❌ Error: zsh is not installed"
    exit 1
fi

if [ ! -d "${HOME}/.oh-my-zsh" ]; then
    echo "❌ Error: oh-my-zsh is not installed"
    echo "Install oh-my-zsh first: https://ohmyz.sh"
    exit 1
fi

echo "📦 Downloading latest release..."
cd "$TMP_DIR"

ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$OS" = "darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        BINARY="vibe-darwin-arm64"
    else
        BINARY="vibe-darwin-amd64"
    fi
elif [ "$OS" = "linux" ]; then
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        BINARY="vibe-linux-arm64"
    else
        BINARY="vibe-linux-amd64"
    fi
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "⚠️  No release found, cloning repository..."
    git clone "https://github.com/${REPO}.git" vibe-zsh
    cd vibe-zsh
    
    if ! command -v go &> /dev/null; then
        echo "❌ Error: Go is not installed. Install Go 1.21+ to build from source."
        exit 1
    fi
    
    echo "🔨 Building from source..."
    go build -o vibe main.go
    
    mkdir -p "$INSTALL_DIR"
    cp vibe "$INSTALL_DIR/"
    cp vibe.plugin.zsh "$INSTALL_DIR/"
    cp _vibe "$INSTALL_DIR/"
else
    echo "📥 Downloading release ${LATEST_RELEASE}..."
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/${BINARY}"
    CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/checksums.txt"
    
    if ! curl -fsSL "$DOWNLOAD_URL" -o vibe; then
        echo "❌ Failed to download binary"
        exit 1
    fi
    
    echo "🔐 Verifying checksum..."
    if curl -fsSL "$CHECKSUMS_URL" -o checksums.txt 2>/dev/null; then
        if command -v sha256sum &> /dev/null; then
            if echo "$(grep "$BINARY" checksums.txt)" | sha256sum -c --status 2>/dev/null; then
                echo "✓ Checksum verified"
            else
                echo "⚠️  Checksum verification failed, but continuing..."
            fi
        elif command -v shasum &> /dev/null; then
            EXPECTED=$(grep "$BINARY" checksums.txt | awk '{print $1}')
            ACTUAL=$(shasum -a 256 vibe | awk '{print $1}')
            if [ "$EXPECTED" = "$ACTUAL" ]; then
                echo "✓ Checksum verified"
            else
                echo "⚠️  Checksum verification failed, but continuing..."
            fi
        fi
        rm -f checksums.txt
    else
        echo "⚠️  Could not download checksums, skipping verification"
    fi
    
    chmod +x vibe
    
    curl -fsSL "https://raw.githubusercontent.com/${REPO}/${LATEST_RELEASE}/vibe.plugin.zsh" -o vibe.plugin.zsh
    curl -fsSL "https://raw.githubusercontent.com/${REPO}/${LATEST_RELEASE}/_vibe" -o _vibe
    
    mkdir -p "$INSTALL_DIR"
    cp vibe "$INSTALL_DIR/"
    cp vibe.plugin.zsh "$INSTALL_DIR/"
    cp _vibe "$INSTALL_DIR/"
fi

cd - > /dev/null
rm -rf "$TMP_DIR"

echo "✅ vibe-zsh installed to ${INSTALL_DIR}"
echo ""
echo "📝 Next steps:"
echo "1. Add 'vibe' to your plugins list in ~/.zshrc:"
echo "   plugins=(... vibe)"
echo ""
echo "2. Configure your LLM provider (example for local Ollama):"
echo "   export VIBE_API_URL=\"http://localhost:11434/v1\""
echo "   export VIBE_MODEL=\"llama3:8b\""
echo ""
echo "3. Reload your shell:"
echo "   source ~/.zshrc"
echo ""
echo "4. Try it out:"
echo "   - Type: list all files"
echo "   - Press: Ctrl+G"
echo ""
echo "🎉 Happy vibing!"
