package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the SQLite database and creates the necessary tables
func InitDB(dbPath string) error {
	if db != nil {
		db.Close()
	}
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Create table for command history
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		command TEXT NOT NULL,
		session_id TEXT,
		cwd TEXT,
		exit_code INTEGER,
		metadata TEXT
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// SaveCommand saves a command to the database
func SaveCommand(command, sessionID, cwd string, exitCode int, metadata string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := db.Exec(
		"INSERT INTO commands (command, session_id, cwd, exit_code, metadata) VALUES (?, ?, ?, ?, ?)",
		command, sessionID, cwd, exitCode, metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to save command: %w", err)
	}

	return nil
}

// SearchCommands searches for commands matching the query
func SearchCommands(query string, limit int) ([]Command, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.Query(`
		SELECT id, timestamp, command, session_id, cwd, exit_code, metadata
		FROM commands
		WHERE command LIKE ?
		ORDER BY id DESC
		LIMIT ?
	`, "%"+query+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search commands: %w", err)
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(&cmd.ID, &cmd.Timestamp, &cmd.Command, &cmd.SessionID, &cmd.CWD, &cmd.ExitCode, &cmd.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to scan command: %w", err)
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

// CleanHistory removes commands from history based on datetime range or all commands
func CleanHistory(fromTime, toTime *time.Time, all, dryRun bool) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	var query string
	var args []interface{}

	if all {
		query = "DELETE FROM commands"
	} else if fromTime != nil && toTime != nil {
		query = "DELETE FROM commands WHERE timestamp BETWEEN ? AND ?"
		args = []interface{}{fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339)}
	} else if fromTime != nil {
		query = "DELETE FROM commands WHERE timestamp >= ?"
		args = []interface{}{fromTime.Format(time.RFC3339)}
	} else if toTime != nil {
		query = "DELETE FROM commands WHERE timestamp <= ?"
		args = []interface{}{toTime.Format(time.RFC3339)}
	} else {
		return 0, fmt.Errorf("at least one datetime range must be specified or use --all flag")
	}

	if dryRun {
		// For dry run, count instead of delete
		var countQuery string
		if all {
			countQuery = "SELECT COUNT(*) FROM commands"
		} else {
			countQuery = strings.Replace(query, "DELETE FROM commands", "SELECT COUNT(*) FROM commands", 1)
		}
		var count int64
		err := db.QueryRow(countQuery, args...).Scan(&count)
		return count, err
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to clean history: %w", err)
	}

	return result.RowsAffected()
}

// Command represents a stored command
type Command struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Command   string `json:"command"`
	SessionID string `json:"session_id"`
	CWD       string `json:"cwd"`
	ExitCode  int    `json:"exit_code"`
	Metadata  string `json:"metadata"`
}
