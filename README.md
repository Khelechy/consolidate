# Consolidate

Consolidate is a lightwieght cross-platform command-line utility that automatically logs every command you run in your terminal, providing a searchable and persistent history. Unlike built-in shell history, Consolidate works across different shells and sessions, with advanced search capabilities and secure storage.

## Use Cases

- **Developers**: Track and search through complex command sequences during development.
- **System Administrators**: Maintain audit trails of commands executed on servers, with the ability to clean up old history.
- **Anyone**: Never lose a useful command again; search through your entire command history and clean up when needed.

## Features

- **Automatic Logging**: Hooks into bash, zsh, and PowerShell to capture commands after execution.
- **Secure Storage**: Uses SQLite with optional encryption for sensitive data.
- **Fast Search**: Full-text search through command history with regex-like queries.
- **History Cleanup**: Remove old or unwanted commands with flexible date-based filtering or delete all history.
- **Cross-Platform**: Works on Windows, Linux, and macOS.
- **CLI Interface**: Simple commands for logging, searching, and managing history.
- **JSON Export**: Export history for analysis or backup.
- **Session Tracking**: Associates commands with sessions, working directories, and exit codes.

**Note**: This tool logs commands after execution to avoid interfering with command behavior. It captures the command as run, including any shell expansions.

## Installation

### Prerequisites

- Go 1.19 or later

### Option 1: Go Install (Recommended)

Install directly from the repository:

```bash
go install github.com/khelechy/consolidate@latest
```

This will download, build, and install consolidate to your `$GOPATH/bin` or `$GOBIN`.

### Option 2: Download Pre-built Binary

1. Go to the [Releases](https://github.com/khelechy/consolidate/releases) page.
2. Download the appropriate binary for your platform.
3. Extract and place `consolidate` in your PATH.

### Option 3: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/khelechy/consolidate.git
   cd consolidate
   ```

2. Build the binary:
   ```bash
   go build -o consolidate
   ```

3. (Optional) Install to system PATH:
   ```bash
   sudo mv consolidate /usr/local/bin/
   ```

### Setup

1. Initialize the database:
   ```bash
   >> consolidate init
   ```

2. Install shell hooks for automatic logging:
   ```bash
   >> consolidate hook
   ```

3. Restart your shell or source your profile to activate hooks.

## Usage

### Basic Commands

#### Initialize Consolidate
Initialize the database and configuration.
```bash
>> consolidate init
```
Sets up the database and configuration files.


#### Install Hooks
Install shell hooks, automatically detects your shell and installs hooks for command logging.
```bash
>> consolidate hook
```


#### View History
Displays all logged commands, ordered from recent to oldest.
```bash
>> consolidate history

>> consolidate history --json > history.json
```
- Flags:
  - `--limit int`: Maximum commands (default 100)
  - `--json`: Output in JSON format


#### Search History
Searches for commands containing "git".
```bash
>> consolidate search "git"

>> consolidate search "docker.*build"
```
#### `consolidate search [query]`
- Flags:
  - `--limit int`: Maximum results (default 10)
  - `--json`: Output in JSON format


#### Manual Logging
Manually log a command (useful for testing or scripting).
```bash
>> consolidate log "echo hello world" --session mysession --cwd /home/user
```
- Flags:
  - `--session string`: Session ID
  - `--cwd string`: Current working directory
  - `--exit-code int`: Exit code (default 0)
  - `--metadata string`: Additional metadata


#### Clean History
Remove commands from history based on date ranges or delete all commands.
```bash
# Delete all commands from history
>> consolidate clean --all

# Delete commands from a specific date range
>> consolidate clean --from 2023-01-01 --to 2023-12-31

# Delete commands from a specific date onwards
>> consolidate clean --from 2023-01-01

# Delete commands up to a specific date
>> consolidate clean --to 2023-12-31

# Preview what would be deleted (dry run)
>> consolidate clean --all --dry-run
>> consolidate clean --from 2023-01-01 --dry-run
```
- Flags:
  - `--all`: Delete all commands from history (cannot be used with --from or --to)
  - `--from string`: Start datetime (RFC3339 or YYYY-MM-DD format, e.g., 2023-01-01 or 2023-01-01T00:00:00Z)
  - `--to string`: End datetime (RFC3339 or YYYY-MM-DD format, e.g., 2023-12-31 or 2023-12-31T23:59:59Z)
  - `--dry-run`: Show what would be deleted without actually deleting


#### `consolidate help [command]`
Get help for any command.


## Configuration

Consolidate stores data in `~/.consolidate/`:
- `history.db`: SQLite database with command history
- Configuration is minimal; most settings are command-line flags

## Security

- Commands are stored locally; no data is sent to external servers.
- Sensitive information in commands (e.g., passwords) should be avoided.
- Future versions may include encryption options.

## Troubleshooting

### Hooks Not Working
- Ensure you've run `consolidate hook` and restarted your shell.
- Check that the hook script is sourced in your shell profile.

### Database Errors
- Run `consolidate init` to recreate the database.
- Ensure write permissions in `~/.consolidate/`.

### Commands Not Logging
- Verify hooks are installed: `consolidate hook`
- Check for errors in your shell after running commands.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

### Development Setup

```bash
git clone https://github.com/khelechy/consolidate.git
cd consolidate
go mod tidy
go test ./...
go build
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Issues: [GitHub Issues](https://github.com/khelechy/consolidate/issues)
- Discussions: [GitHub Discussions](https://github.com/khelechy/consolidate/discussions)

---
