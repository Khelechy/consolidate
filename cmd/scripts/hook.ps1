# Consolidate hook script for PowerShell
# Add this to your PowerShell profile (e.g., $PROFILE)

# Get the path to the consolidate binary
# This will be set by the hook installation
$ConsolidateBin = if ($env:CONSOLIDATE_BIN) { $env:CONSOLIDATE_BIN } else { "consolidate" }

# Function to log command after execution
function Log-Command {
    param([string]$LastCommand, [int]$ExitCode, [string]$Cwd, [string]$SessionId)

    # Skip logging if command is empty
    if ([string]::IsNullOrWhiteSpace($LastCommand)) { return }

    # Skip logging consolidate commands to avoid recursion
    if ($LastCommand -match '^(\./)?consolidate(\.exe)?') { return }

    # Encode the command to avoid parsing issues
    $encodedCommand = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes($LastCommand))

    # Log the command
    try {
        & $ConsolidateBin log $encodedCommand --encoded --session $SessionId --cwd $Cwd --exit-code $ExitCode 2>$null
    } catch {
        # Ignore errors
    }
}

# Set up the hook using Register-EngineEvent for command execution
$null = Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action {
    # This runs when PowerShell exits, but we need per-command
    # For per-command, we might need to override prompt or use a different approach
    # For now, this is a placeholder
}

# Alternative: Override the prompt function to log after each command
$originalPrompt = $function:prompt
function prompt {
    # Log the last command
    if ($?) {
        $exitCode = 0
    } else {
        $exitCode = 1
    }
    
    # Get the last command, preferring PSReadLine for full command text
    $lastCommand = if (Get-Module PSReadLine -ErrorAction SilentlyContinue) {
        try {
            [Microsoft.PowerShell.PSConsoleReadLine]::GetHistoryItems() | Select-Object -Last 1 -ExpandProperty CommandLine
        } catch {
            (Get-History -Count 1).CommandLine
        }
    } else {
        (Get-History -Count 1).CommandLine
    }
    
    Log-Command -LastCommand $lastCommand -ExitCode $exitCode -Cwd (Get-Location).Path -SessionId $PID

    # Call original prompt
    & $originalPrompt
}