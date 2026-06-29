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

    # Note: Interactive confirmation is now handled in the Go binary
    # when VIBE_INTERACTIVE=true is set

    # The binary wrote the spinner/explanations to the terminal (stderr) while
    # this widget was active, so ZLE's view of the cursor position is now stale.
    # zle -I invalidates the display and tells ZLE that external output occurred;
    # without it, reset-prompt redraws relative to the wrong position and the
    # prompt (especially RPROMPT) ends up visually shifted. Set the buffer first,
    # then do a single reset-prompt so the prompt is redrawn cleanly once.
    zle -I

    BUFFER="$cmd"
    CURSOR=${#BUFFER}
    zle reset-prompt
  else
    zle -M "vibe: Failed to generate command"
  fi
  
  # Check for updates in background (silent)
  (nohup "$VIBE_BINARY" check-update >/dev/null 2>&1 &)
}

zle -N vibe

bindkey '^G' vibe

# History widget - shows interactive history menu
function vibe-history-widget() {
  # Call vibe-zsh history command
  local output=$("$VIBE_BINARY" history 2>/dev/null)
  local exit_code=$?
  
  if [[ $exit_code -eq 0 && -n "$output" ]]; then
    # The history menu is a full-screen TUI rendered to the terminal, so ZLE's
    # display state is stale on return. Invalidate it before redrawing so the
    # prompt isn't left visually shifted.
    zle -I

    # Check if this is a regeneration request
    if [[ "$output" == REGENERATE:* ]]; then
      # Extract the query (remove REGENERATE: prefix)
      local query="${output#REGENERATE:}"

      # Set the query in the buffer and trigger vibe
      BUFFER="$query"
      CURSOR=${#BUFFER}
      zle reset-prompt

      # Call the vibe widget to regenerate
      vibe
    elif [[ "$output" == EDIT:* ]]; then
      # Extract the query (remove EDIT: prefix)
      local query="${output#EDIT:}"

      # Put query in buffer for editing
      BUFFER="$query"
      CURSOR=${#BUFFER}
      zle reset-prompt
    else
      # Normal command insertion
      BUFFER="$output"
      CURSOR=${#BUFFER}
      zle reset-prompt
    fi
  fi
}

zle -N vibe-history-widget

# Configurable keybinding for history (default: Ctrl+X H)
# Users can override with: export VIBE_HISTORY_KEY="^R" for Ctrl+R, etc.
local history_key="${VIBE_HISTORY_KEY:-^Xh}"
bindkey "$history_key" vibe-history-widget

# Regenerate last command widget
# Fetches the most recent query and regenerates a new command
function vibe-regenerate-last-widget() {
  # Get the most recent query from history
  local query=$("$VIBE_BINARY" history last 2>/dev/null)
  local exit_code=$?
  
  if [[ $exit_code -eq 0 && -n "$query" ]]; then
    # Invalidate ZLE's display in case the lookup wrote to the terminal.
    zle -I

    # Set the query in the buffer and trigger vibe
    BUFFER="$query"
    CURSOR=${#BUFFER}
    zle reset-prompt

    # Call the vibe widget to regenerate
    vibe
  else
    zle -M "vibe: No history found"
  fi
}

zle -N vibe-regenerate-last-widget

# Configurable keybinding for regenerate last (default: Ctrl+X G)
# Users can override with: export VIBE_REGENERATE_KEY="^Xg"
local regenerate_key="${VIBE_REGENERATE_KEY:-^Xg}"
bindkey "$regenerate_key" vibe-regenerate-last-widget

# Shell function for direct command-line use
# Usage: vh
# 
# This provides the same interactive history experience as Ctrl+X H
# but can be invoked by typing 'vh' on the command line.
# 
# The selected command is inserted into your ZSH buffer using print -z,
# making it appear on your command line ready for execution.
# Press 'g' in the menu to regenerate a command from its original query.
function vh() {
  local output=$("$VIBE_BINARY" history 2>/dev/null)
  local exit_code=$?
  
  if [[ $exit_code -eq 0 && -n "$output" ]]; then
    # Check if this is a regeneration request
    if [[ "$output" == REGENERATE:* ]]; then
      # Extract the query (remove REGENERATE: prefix)
      local query="${output#REGENERATE:}"
      
      # Generate new command from the query
      local new_cmd=$("$VIBE_BINARY" "$query")
      if [[ -n "$new_cmd" ]]; then
        # Extract only the first line (the actual command)
        new_cmd="${new_cmd%%$'\n'*}"
        print -z "$new_cmd"
      fi
    elif [[ "$output" == EDIT:* ]]; then
      # Extract the query (remove EDIT: prefix)
      local query="${output#EDIT:}"
      
      # Put query in buffer for editing
      print -z "$query"
    else
      # Use print -z to add command to the editing buffer
      # This makes it appear on the command line ready for execution
      print -z "$output"
    fi
  fi
}

# Add this plugin's directory to fpath so the _vibe completion is discoverable.
# Do NOT run compinit here: the completion system is initialized by the shell
# framework (e.g. oh-my-zsh) or the user's own zshrc. Calling compinit a second
# time at plugin-load time re-initializes completions after the prompt/theme is
# already set up, which can leave the prompt visually shifted (especially with a
# right-side prompt/RPROMPT).
fpath+="${VIBE_PLUGIN_DIR}"
