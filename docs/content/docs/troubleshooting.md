---
title: "Troubleshooting"
weight: 4
---


## Command Not Working After Install

**Symptoms:** Pressing `Ctrl+G` does nothing, or you get "command not found" errors.

**Solutions:**

1. **Reload your shell:**
   ```bash
   source ~/.zshrc
   ```

2. **Verify plugin is loaded:**
   ```bash
   echo $plugins
   # Should include "vibe"
   ```

3. **Check binary exists:**
   ```bash
   ls ~/.oh-my-zsh/custom/plugins/vibe/vibe
   # Should show the binary
   ```

4. **Verify function is loaded:**
   ```bash
   which vibe
   # Should show "vibe" as a shell function
   ```

5. **Rebuild the binary:**
   ```bash
   cd ~/.oh-my-zsh/custom/plugins/vibe
   make clean build
   ```

## Slow Responses

**Symptoms:** vibe takes several seconds to respond.

**Solutions:**

1. **Enable caching:**
   ```bash
   export VIBE_ENABLE_CACHE=true
   ```
   Cached responses are 100-400x faster!

2. **Use a faster model:**
   ```bash
   export VIBE_MODEL="llama3:8b"  # Instead of larger models
   ```

3. **Use a local LLM:**
   ```bash
   export VIBE_API_URL="http://localhost:11434/v1"
   export VIBE_MODEL="llama3:8b"
   ```

4. **Check network connection:**
   ```bash
   curl -I $VIBE_API_URL
   ```

5. **Increase timeout if needed:**
   ```bash
   export VIBE_TIMEOUT=60s
   ```

## Bad Command Suggestions

**Symptoms:** Generated commands don't match your intent or are incorrect.

**Solutions:**

1. **Try a different/better model:**
   ```bash
   export VIBE_MODEL="gpt-4"  # Or other capable models
   ```

2. **Make your query more specific:**
   - ❌ "list files"
   - ✅ "list all files in current directory including hidden ones with sizes"

3. **Check model supports structured output:**
   Some models may not follow the expected output format.

4. **Lower temperature for more deterministic results:**
   ```bash
   export VIBE_TEMPERATURE=0.3
   ```

5. **Clear cache if getting stale results:**
   ```bash
   rm -rf ~/.cache/vibe/*
   ```

## Ctrl+G Does Nothing

**Symptoms:** Pressing `Ctrl+G` has no effect.

**Solutions:**

1. **Check if another plugin uses Ctrl+G:**
   ```bash
   bindkey | grep '"\^G"'
   ```

2. **Rebind to different key:**
   ```bash
   bindkey '^X' vibe  # Use Ctrl+X instead
   ```

3. **Ensure plugin is loaded:**
   ```bash
   echo $plugins | grep vibe
   ```

4. **Reload shell:**
   ```bash
   source ~/.zshrc
   ```

## API Connection Errors

**Symptoms:** Errors like "connection refused" or "timeout."

**Solutions:**

1. **For Ollama - check if it's running:**
   ```bash
   curl http://localhost:11434/v1/models
   ```

2. **Start Ollama if needed:**
   ```bash
   ollama serve
   ```

3. **Verify API URL is correct:**
   ```bash
   echo $VIBE_API_URL
   ```

4. **Check API key is set (if using OpenAI/Claude):**
   ```bash
   echo $VIBE_API_KEY
   # Should not be empty
   ```

5. **Test API directly:**
   ```bash
   curl -H "Authorization: Bearer $VIBE_API_KEY" $VIBE_API_URL/models
   ```

## Cache Issues

**Symptoms:** Getting outdated commands or cache-related errors.

**Solutions:**

1. **Clear the cache:**
   ```bash
   rm -rf ~/.cache/vibe/*
   ```

2. **Disable cache temporarily:**
   ```bash
   export VIBE_ENABLE_CACHE=false
   ```

3. **Reduce cache TTL:**
   ```bash
   export VIBE_CACHE_TTL=1h
   ```

## Model Not Found

**Symptoms:** Error messages about model not being available.

**Solutions:**

1. **For Ollama - pull the model:**
   ```bash
   ollama pull llama3:8b
   ```

2. **List available models:**
   ```bash
   # Ollama
   ollama list
   
   # OpenAI
   curl https://api.openai.com/v1/models \
     -H "Authorization: Bearer $VIBE_API_KEY"
   ```

3. **Use a different model:**
   ```bash
   export VIBE_MODEL="llama3:8b"
   ```

## Permission Denied

**Symptoms:** "Permission denied" when trying to execute vibe binary.

**Solutions:**

1. **Make binary executable:**
   ```bash
   chmod +x ~/.oh-my-zsh/custom/plugins/vibe/vibe
   ```

2. **Rebuild binary:**
   ```bash
   cd ~/.oh-my-zsh/custom/plugins/vibe
   make clean build
   ```

## Still Having Issues?

If you're still experiencing problems:

1. **Check the logs:**
   ```bash
   # vibe logs errors to stderr
   vibe "test query" 2>&1
   ```

2. **Enable debug mode (if available):**
   ```bash
   export VIBE_DEBUG=true
   ```

3. **Report the issue:**
   - Visit: https://github.com/skymoore/vibe-zsh/issues
   - Include:
     - Your OS and version
     - Zsh version: `zsh --version`
     - Oh-My-Zsh version
     - vibe configuration (environment variables)
     - Error messages or unexpected behavior

4. **Join the community:**
   - Check existing issues for solutions
   - Ask questions in discussions
   - Contribute fixes if you find them!
