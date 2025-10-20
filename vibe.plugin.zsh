VIBE_PLUGIN_DIR="${0:A:h}"
VIBE_BINARY="${VIBE_BINARY:-${VIBE_PLUGIN_DIR}/vibe}"

function vibe() {
  local request="$BUFFER"
  
  if [[ -z "$request" ]]; then
    return
  fi
  
  local cmd=$("$VIBE_BINARY" "$request" 2>&1)
  local exit_code=$?
  
  if [[ $exit_code -eq 0 && -n "$cmd" ]]; then
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
    
    BUFFER=""
    cursor=0
    zle reset-prompt
    BUFFER="$cmd"
  else
    zle -M "vibe: Failed to generate command${cmd:+ - }${cmd}"
  fi
}

zle -N vibe

bindkey '^G' vibe

fpath+="${VIBE_PLUGIN_DIR}"
autoload -Uz compinit
compinit
