package logs_test

import (
	"dbx/internal/logs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLogEntry_Success tests successful log entry creation
func TestLogEntry_Success(t *testing.T) {
	// Note: LogEntry uses hardcoded "./logs" directory
	// In a real scenario, this should be configurable via dependency injection
	logDir := "./logs"
	os.MkdirAll(logDir, 0755)
	defer os.RemoveAll(logDir)

	startTime := time.Now().Add(-5 * time.Second)
	logs.LogEntry("mysql", "Backup", "SUCCESS", startTime, nil)

	// Verify log file was created
	logFile := filepath.Join(logDir, "dbx.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("LogEntry() did not create log file")
	}

	// Verify log content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("LogEntry() did not write content to log file")
	}

	// Verify log contains expected information
	logStr := string(content)
	if !contains(logStr, "mysql") {
		t.Error("LogEntry() log does not contain database type")
	}
	if !contains(logStr, "Backup") {
		t.Error("LogEntry() log does not contain operation")
	}
	if !contains(logStr, "SUCCESS") {
		t.Error("LogEntry() log does not contain status")
	}
}

// TestLogEntry_WithError tests log entry with error
func TestLogEntry_WithError(t *testing.T) {
	logDir := t.TempDir()
	os.Setenv("DBX_LOG_DIR", logDir)

	startTime := time.Now().Add(-2 * time.Second)
	testErr := os.ErrNotExist
	logs.LogEntry("postgres", "Restore", "FAILED", startTime, testErr)

	logFile := filepath.Join(logDir, "dbx.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logStr := string(content)
	if !contains(logStr, "FAILED") {
		t.Error("LogEntry() log does not contain FAILED status")
	}
	if !contains(logStr, "Error") {
		t.Error("LogEntry() log does not contain error information")
	}
}

// TestLogEntry_MultipleEntries tests multiple log entries
func TestLogEntry_MultipleEntries(t *testing.T) {
	logDir := t.TempDir()
	os.Setenv("DBX_LOG_DIR", logDir)

	startTime := time.Now()
	logs.LogEntry("mysql", "Backup", "SUCCESS", startTime, nil)
	logs.LogEntry("postgres", "Backup", "SUCCESS", startTime, nil)
	logs.LogEntry("mongodb", "Backup", "SUCCESS", startTime, nil)

	logFile := filepath.Join(logDir, "dbx.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logStr := string(content)
	// Count newlines to verify multiple entries
	entryCount := countLines(logStr)
	if entryCount < 3 {
		t.Errorf("LogEntry() wrote %d entries, want at least 3", entryCount)
	}
}

// TestLogEntry_ConcurrentWrites tests concurrent log writes (scalability)
func TestLogEntry_ConcurrentWrites(t *testing.T) {
	logDir := t.TempDir()
	os.Setenv("DBX_LOG_DIR", logDir)

	numGoroutines := 20
	done := make(chan bool, numGoroutines)

	startTime := time.Now()
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			logs.LogEntry("mysql", "Backup", "SUCCESS", startTime, nil)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	logFile := filepath.Join(logDir, "dbx.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Verify all entries were written
	entryCount := countLines(string(content))
	if entryCount < numGoroutines {
		t.Errorf("Concurrent LogEntry() wrote %d entries, want at least %d", entryCount, numGoroutines)
	}
}

// TestLogEntry_InvalidLogDir tests error handling for invalid log directory
func TestLogEntry_InvalidLogDir(t *testing.T) {
	// Set an invalid log directory (path with invalid characters on Windows)
	var invalidDir string
	if os.PathSeparator == '\\' {
		invalidDir = "C:\\<invalid>\\logs"
	} else {
		invalidDir = "/root/nonexistent/logs"
	}

	os.Setenv("DBX_LOG_DIR", invalidDir)
	defer os.Unsetenv("DBX_LOG_DIR")

	startTime := time.Now()
	// Should not panic even with invalid directory
	logs.LogEntry("mysql", "Backup", "SUCCESS", startTime, nil)
	t.Log("LogEntry() handled invalid log directory gracefully")
}

// TestLogEntry_EmptyFields tests logging with empty fields
func TestLogEntry_EmptyFields(t *testing.T) {
	logDir := t.TempDir()
	os.Setenv("DBX_LOG_DIR", logDir)

	startTime := time.Now()
	logs.LogEntry("", "", "", startTime, nil)

	logFile := filepath.Join(logDir, "dbx.log")
	content, err := os.ReadFile(logFile)
	if err == nil {
		logStr := string(content)
		if len(logStr) == 0 {
			t.Error("LogEntry() should write log even with empty fields")
		}
	}
}

// TestLogEntry_WriteError tests handling of write errors
func TestLogEntry_WriteError(t *testing.T) {
	logDir := t.TempDir()
	os.Setenv("DBX_LOG_DIR", logDir)

	// Create log file and make it read-only to simulate write error
	logFile := filepath.Join(logDir, "dbx.log")
	os.WriteFile(logFile, []byte("existing"), 0444) // Read-only

	startTime := time.Now()
	// Should not panic on write error
	logs.LogEntry("mysql", "Backup", "SUCCESS", startTime, nil)
	t.Log("LogEntry() handled write error gracefully")

	// Restore permissions for cleanup
	os.Chmod(logFile, 0644)
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func countLines(s string) int {
	count := 0
	for _, char := range s {
		if char == '\n' {
			count++
		}
	}
	return count
}

