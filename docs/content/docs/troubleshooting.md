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

## Corrupted or Garbage Output

**Symptoms:** Response contains garbage characters, escape sequences, or Unicode pollution like:
```
AWS?…..??……… ...…..??……....??…... …………â[201~[200~¦………
```

**Solutions:**

vibe automatically handles corrupted LLM responses with multi-layer parsing. If you're still seeing issues:

1. **Enable debug logs to see what's happening:**
   ```bash
   export VIBE_DEBUG_LOGS=true
   ```
   This will show:
   - Raw LLM responses
   - Parsing attempts and which layer succeeded/failed
   - Extracted JSON content

2. **Increase retry attempts:**
   ```bash
   export VIBE_MAX_RETRIES=5  # Default is 3
   ```

3. **Verify JSON extraction is enabled:**
   ```bash
   export VIBE_ENABLE_JSON_EXTRACTION=true  # Default
   ```

4. **Try lower temperature for cleaner output:**
   ```bash
   export VIBE_TEMPERATURE=0.3  # More deterministic
   ```

5. **Check which parsing layer is working:**
   Debug logs will show messages like:
   - `[SUCCESS] Layer 'structured_output' succeeded on attempt 1`
   - `[SUCCESS] Layer 'enhanced_parsing' succeeded on attempt 2`
   - `[JSON_EXTRACT] Success` - JSON was found in corrupted response

## Parsing Errors

**Symptoms:** Errors like "failed to parse text response" or "unable to extract valid JSON."

**Solutions:**

1. **Enable JSON extraction (should be default):**
   ```bash
   export VIBE_ENABLE_JSON_EXTRACTION=true
   ```

2. **Temporarily disable strict validation for testing:**
   ```bash
   export VIBE_STRICT_VALIDATION=false
   ```
   This bypasses validation of command/explanation fields. Useful for debugging.

3. **Check debug logs to see parsing attempts:**
   ```bash
   export VIBE_DEBUG_LOGS=true
   vibe "your query" 2>&1 | grep -E "\[ATTEMPT|SUCCESS|EXTRACT\]"
   ```

4. **Try with structured output disabled:**
   ```bash
   export VIBE_USE_STRUCTURED_OUTPUT=false
   ```
   This forces text parsing mode, which may work better for some models.

5. **Some models produce cleaner JSON with lower temperature:**
   ```bash
   export VIBE_TEMPERATURE=0.3
   ```

6. **Report persistent parsing issues:**
   If a specific query consistently fails, please report it with:
   - The query that failed
   - Debug logs (with sensitive info redacted)
   - Your model and provider

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

## Understanding Debug Output

When `VIBE_DEBUG_LOGS=true`, you'll see structured log messages:

```
[VIBE] [ATTEMPT 1][structured_output] Parsing failed: ...
[VIBE] [ATTEMPT 2][enhanced_parsing] Parsing failed: ...
[VIBE] [JSON_EXTRACT] Success
Trimmed prefix: "AWS\x1b[200~..."
Trimmed suffix: "...\x1b[201~"
[VIBE] [SUCCESS] Layer 'enhanced_parsing' succeeded on attempt 2
```

**Key log types:**
- `[ATTEMPT N][layer]` - Shows which parsing layer is being tried
- `[JSON_EXTRACT]` - Shows JSON extraction from corrupted response
- `[SUCCESS]` - Shows which layer and attempt succeeded
- `Raw response (first 500 chars)` - Shows actual LLM output

**Multi-layer fallback order:**
1. `structured_output` - JSON schema mode (best)
2. `enhanced_parsing` - Retry with JSON extraction
3. `explicit_json_prompt` - Extra strict prompt with lower temperature
4. `emergency_fallback` - Returns helpful error message

## Still Having Issues?

If you're still experiencing problems:

1. **Check the logs with debug enabled:**
   ```bash
   export VIBE_DEBUG_LOGS=true
   vibe "test query" 2>&1
   ```

2. **Test with minimal config:**
   ```bash
   # Reset to defaults
   unset VIBE_TEMPERATURE VIBE_MAX_TOKENS VIBE_USE_STRUCTURED_OUTPUT
   export VIBE_DEBUG_LOGS=true
   vibe "list files"
   ```

3. **Report the issue:**
   - Visit: https://github.com/skymoore/vibe-zsh/issues
   - Include:
     - Your OS and version
     - Zsh version: `zsh --version`
     - Oh-My-Zsh version
     - vibe configuration (environment variables)
     - Error messages or unexpected behavior
     - Debug logs (redact any sensitive info)
     - Which model/provider you're using

4. **Join the community:**
   - Check existing issues for solutions
   - Ask questions in discussions
   - Contribute fixes if you find them!
