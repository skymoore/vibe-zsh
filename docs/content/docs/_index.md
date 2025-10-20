---
title: Documentation
---

# vibe Documentation

Welcome to the vibe documentation! This guide will help you get started with transforming natural language into shell commands using AI.

## Quick Links

- [Installation](installation) - Get vibe installed on your system
- [Configuration](configuration) - Configure your LLM provider
- [Usage](usage) - Learn how to use vibe with examples
- [Troubleshooting](troubleshooting) - Common issues and solutions
- [Configuration Reference](reference) - Complete environment variable reference
- [API Reference](api) - Developer documentation

## What is vibe?

vibe is a zsh plugin that transforms natural language into shell commands using AI. Just type what you want to do in plain English, press `Ctrl+G`, and get the command with inline explanations.

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
