package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/khelechy/consolidate/internal/storage"
)

// GetDBPath returns the path to the consolidate database
func GetDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(homeDir, ".consolidate", "history.db"), nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	dbPath, err := GetDBPath()
	if err != nil {
		return err
	}
	configDir := filepath.Dir(dbPath)
	return os.MkdirAll(configDir, 0755)
}

// InitAndGetDB initializes the database and returns the path
func InitAndGetDB() (string, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return "", err
	}
	if err := storage.InitDB(dbPath); err != nil {
		return "", fmt.Errorf("initializing database: %w", err)
	}
	return dbPath, nil
}

// PrintCommands prints the commands to stdout, either as JSON or formatted text
func PrintCommands(commands []storage.Command, jsonOutput bool) error {
	if jsonOutput {
		jsonData, err := json.MarshalIndent(commands, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling to JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		for _, cmd := range commands {
			fmt.Printf("[%s] %s (exit: %d)\n", cmd.Timestamp, cmd.Command, cmd.ExitCode)
		}
	}
	return nil
}

// DetectShell detects the current shell environment
func DetectShell() string {
	// Check environment variables
	if shell := os.Getenv("ZSH_VERSION"); shell != "" {
		return "zsh"
	}
	if shell := os.Getenv("BASH_VERSION"); shell != "" {
		return "bash"
	}

	// Fallback to SHELL environment variable
	shellPath := os.Getenv("SHELL")
	if strings.Contains(shellPath, "zsh") {
		return "zsh"
	}
	if strings.Contains(shellPath, "bash") {
		return "bash"
	}

	// Check for Windows shells
	if os.Getenv("OS") == "Windows_NT" {
		if os.Getenv("PROMPT") == "" && os.Getenv("PSModulePath") != "" {
			return "powershell"
		}
		return "cmd"
	}

	// Default to bash
	return "bash"
}
