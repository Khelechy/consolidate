#!/bin/bash

# Consolidate hook script for bash/zsh
# Source this in your shell profile (e.g., .bashrc, .zshrc)

# Get the path to the consolidate binary
# This will be set by the hook installation
CONSOLIDATE_BIN="${CONSOLIDATE_BIN:-consolidate}"

# Function to log command after execution
_log_command() {
    local exit_code=$?
    local cwd=$(pwd)
    local session_id=$$

    # Get the last command from history
    local last_command
    if [[ -n "$ZSH_VERSION" ]]; then
        last_command="${history[$HISTCMD]}"
    else
        last_command=$(fc -ln -1 2>/dev/null | sed 's/^[[:space:]]*[0-9]*[[:space:]]*//')
    fi

    # Skip logging if command is empty or starts with space (bash histcontrol)
    [[ -z "$last_command" ]] && return
    [[ "$last_command" =~ ^[[:space:]] ]] && return

    # Skip logging consolidate commands to avoid recursion
    [[ "$last_command" =~ ^(\./)?consolidate(\.exe)? ]] && return

    # Log the command
    encoded_command=$(echo -n "$last_command" | base64)
    $CONSOLIDATE_BIN log "$encoded_command" --encoded --session "$session_id" --cwd "$cwd" --exit-code "$exit_code" 2>/dev/null || true
}

# Set up the hook
if [[ -n "$ZSH_VERSION" ]]; then
    # Zsh
    autoload -Uz add-zsh-hook
    add-zsh-hook precmd _log_command
elif [[ -n "$BASH_VERSION" ]]; then
    # Bash
    PROMPT_COMMAND="_log_command"
fi