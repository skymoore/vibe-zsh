VIBE_PLUGIN_DIR="${0:A:h}"
VIBE_BINARY="${VIBE_BINARY:-${VIBE_PLUGIN_DIR}/vibe}"

function vibe() {
  local request="$BUFFER"
  
  if [[ -z "$request" ]]; then
    return
  fi
  
  # Call vibe binary
  # - stdout: Command output (captured)
  # - stderr: Progress spinner (displayed to user, not captured)
  local output=$("$VIBE_BINARY" "$request")
  local exit_code=$?
  
  if [[ $exit_code -eq 0 && -n "$output" ]]; then
    # Extract only the first line (the actual command)
    # Everything else is explanations that were already displayed during streaming
    local cmd="${output%%$'\n'*}"
    
    if [[ "$VIBE_INTERACTIVE" == "true" ]]; then
      echo "\n---"
      echo "$cmd"
      echo "---"
      read "confirm?Execute this command? [Y/n] "
      if [[ "$confirm" =~ ^[Nn] ]]; then
        zle -M "Command cancelled"
        return
      fi
    fi
    
    # Clear buffer and reset prompt
    BUFFER=""
    CURSOR=0
    zle reset-prompt
    
    # Set new buffer and move cursor to end (only the command, not explanations)
    BUFFER="$cmd"
    CURSOR=${#BUFFER}
    zle redisplay
  else
    zle -M "vibe: Failed to generate command"
  fi
  
  # Check for updates in background (silent)
  (nohup "$VIBE_BINARY" check-update >/dev/null 2>&1 &)
}

zle -N vibe

bindkey '^G' vibe

fpath+="${VIBE_PLUGIN_DIR}"
autoload -Uz compinit
compinit
