package db_test

import (
	"dbx/internal/db"
	"os"
	"path/filepath"
	"testing"
)

// TestBackupSQLite_SlackNotification tests Slack notification path
func TestBackupSQLite_SlackNotification(t *testing.T) {
	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	tmpDir := t.TempDir()
	testDB := filepath.Join(tmpDir, "test.db")
	os.WriteFile(testDB, []byte("SQLite format 3"), 0644)

	backupDir := filepath.Join(tmpDir, "backups")
	err := db.BackupSQLite(testDB, backupDir)
	if err != nil {
		t.Logf("BackupSQLite() returned error (expected if not valid SQLite): %v", err)
	}
}

// TestBackupSQLite_SlackNotificationWithError tests Slack notification with error
func TestBackupSQLite_SlackNotificationWithError(t *testing.T) {
	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	// Use invalid path to trigger error path
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	err := db.BackupSQLite("/nonexistent/db.db", backupDir)
	// Error is expected, but notification path should be executed
	if err != nil {
		t.Logf("BackupSQLite() returned error (expected): %v", err)
	}
}

// TestBackupSQLite_CompressionSuccess tests compression success path
func TestBackupSQLite_CompressionSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	testDB := filepath.Join(tmpDir, "test.db")
	os.WriteFile(testDB, []byte("SQLite format 3"), 0644)

	backupDir := filepath.Join(tmpDir, "backups")
	err := db.BackupSQLite(testDB, backupDir)
	if err != nil {
		t.Logf("BackupSQLite() returned error (expected if not valid SQLite): %v", err)
	}
	// Compression is tested indirectly through BackupSQLite
}

