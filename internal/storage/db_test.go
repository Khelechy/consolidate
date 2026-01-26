package storage

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	// Use in-memory database for testing
	dbPath := ":memory:"

	err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
}

func TestSaveCommand(t *testing.T) {
	dbPath := ":memory:"

	err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	// Test saving a command
	err = SaveCommand("ls -la", "session1", "/home", 0, "")
	if err != nil {
		t.Errorf("SaveCommand failed: %v", err)
	}

	// Test saving another command
	err = SaveCommand("echo hello", "session1", "/home", 0, "test metadata")
	if err != nil {
		t.Errorf("SaveCommand failed: %v", err)
	}

	// Test saving with empty fields
	err = SaveCommand("", "", "", 1, "")
	if err != nil {
		t.Errorf("SaveCommand with empty fields failed: %v", err)
	}
}

func TestSearchCommands(t *testing.T) {
	dbPath := ":memory:"

	err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	// Save some test commands
	commands := []struct {
		cmd  string
		sess string
		cwd  string
		exit int
		meta string
	}{
		{"git status", "sess1", "/repo", 0, ""},
		{"ls -la", "sess1", "/home", 0, ""},
		{"echo hello world", "sess2", "/tmp", 1, "failed"},
		{"docker build", "sess1", "/project", 0, ""},
	}

	for _, c := range commands {
		err := SaveCommand(c.cmd, c.sess, c.cwd, c.exit, c.meta)
		if err != nil {
			t.Fatalf("SaveCommand failed: %v", err)
		}
	}

	// Test search for "git"
	results, err := SearchCommands("git", 10)
	if err != nil {
		t.Errorf("SearchCommands failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'git', got %d", len(results))
	}
	if results[0].Command != "git status" {
		t.Errorf("Expected 'git status', got '%s'", results[0].Command)
	}

	// Test search for "echo"
	results, err = SearchCommands("echo", 10)
	if err != nil {
		t.Errorf("SearchCommands failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'echo', got %d", len(results))
	}

	// Test search for "ls"
	results, err = SearchCommands("ls", 10)
	if err != nil {
		t.Errorf("SearchCommands failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'ls', got %d", len(results))
	}

	// Test search with no matches
	results, err = SearchCommands("nonexistent", 10)
	if err != nil {
		t.Errorf("SearchCommands failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'nonexistent', got %d", len(results))
	}

	// Test limit
	results, err = SearchCommands("", 2) // Empty query should match all, but limit to 2
	if err != nil {
		t.Errorf("SearchCommands failed: %v", err)
	}
	if len(results) > 2 {
		t.Errorf("Expected at most 2 results, got %d", len(results))
	}

	// Test with limit 0 (should return nothing)
	results, err = SearchCommands("git", 0)
	if err != nil {
		t.Errorf("SearchCommands failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results with limit 0, got %d", len(results))
	}
}

func TestSearchCommandsOrder(t *testing.T) {
	dbPath := ":memory:"

	err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	// Save commands (assuming timestamps are sequential)
	SaveCommand("first", "s", "/", 0, "")
	SaveCommand("second", "s", "/", 0, "")
	SaveCommand("third", "s", "/", 0, "")

	results, err := SearchCommands("", 10)
	if err != nil {
		t.Errorf("SearchCommands failed: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// Should be ordered by id DESC (newest first)
	if results[0].Command != "third" || results[1].Command != "second" || results[2].Command != "first" {
		t.Errorf("Results not in correct order: %v", results)
	}
}
