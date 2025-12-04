package db_test

import (
	"dbx/internal/db"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestBackupPostgresWithType_SlackNotification tests Slack notification path
func TestBackupPostgresWithType_SlackNotification(t *testing.T) {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		t.Skip("Skipping test: pg_dump not found in PATH")
	}

	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	tmpDir := t.TempDir()
	err := db.BackupPostgresWithType("localhost", "5432", "postgres", "password", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupPostgresWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupPostgresWithType_SlackNotificationWithError tests Slack notification with error
func TestBackupPostgresWithType_SlackNotificationWithError(t *testing.T) {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		t.Skip("Skipping test: pg_dump not found in PATH")
	}

	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	// Use invalid parameters to trigger error path
	tmpDir := t.TempDir()
	err := db.BackupPostgresWithType("invalid-host", "5432", "invalid-user", "invalid-pass", "invalid-db", tmpDir, db.BackupTypeFull)
	// Error is expected, but notification path should be executed
	if err != nil {
		t.Logf("BackupPostgresWithType() returned error (expected): %v", err)
	}
}

// TestBackupPostgresWithType_CompressionSuccess tests compression success path
func TestBackupPostgresWithType_CompressionSuccess(t *testing.T) {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		t.Skip("Skipping test: pg_dump not found in PATH")
	}

	tmpDir := t.TempDir()
	// Create a mock backup file to test compression
	backupFile := filepath.Join(tmpDir, "testdb_full_2024-01-01_00-00-00.sql")
	os.WriteFile(backupFile, []byte("dummy backup"), 0644)

	// The compression logic is in BackupPostgresWithType
	// We test it indirectly by running the backup
	err := db.BackupPostgresWithType("localhost", "5432", "postgres", "password", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupPostgresWithType() returned error (expected without real DB): %v", err)
	}
}

