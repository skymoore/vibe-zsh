#!/usr/bin/env bash

set -e

INSTALL_DIR="${HOME}/.oh-my-zsh/custom/plugins/vibe"

echo "🗑️  Uninstalling vibe-zsh..."

if [ ! -d "$INSTALL_DIR" ]; then
    echo "⚠️  vibe-zsh is not installed at ${INSTALL_DIR}"
    exit 0
fi

echo "📦 Removing installation directory..."
rm -rf "$INSTALL_DIR"

echo "✅ vibe-zsh uninstalled successfully"
echo ""
echo "📝 Next steps:"
echo "1. Remove 'vibe' from your plugins list in ~/.zshrc"
echo "2. Remove VIBE_* environment variables from your shell config (if any)"
echo "3. Reload your shell:"
echo "   source ~/.zshrc"
