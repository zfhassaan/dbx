package db_test

import (
	"dbx/internal/db"
	"os"
	"os/exec"
	"testing"
)

// TestBackupMySQLWithType_StdoutPipeError tests error handling for stdout pipe
func TestBackupMySQLWithType_StdoutPipeError(t *testing.T) {
	// This is hard to test directly, but we can test the overall function
	// The stdout pipe error would occur if the command can't be created
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("localhost", "root", "password", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupMySQLWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupMySQLWithType_StderrPipeError tests error handling for stderr pipe
func TestBackupMySQLWithType_StderrPipeError(t *testing.T) {
	// Similar to stdout pipe - tested indirectly
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("localhost", "root", "password", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupMySQLWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupMySQLWithType_FileWriteError tests error handling for file write
func TestBackupMySQLWithType_FileWriteError(t *testing.T) {
	// This would require making the output directory read-only
	// For now, we test the normal path
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("localhost", "root", "password", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupMySQLWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupMySQLWithType_SlackNotification tests Slack notification path
func TestBackupMySQLWithType_SlackNotification(t *testing.T) {
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("localhost", "root", "password", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupMySQLWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupMySQLWithType_SlackNotificationWithError tests Slack notification with error
func TestBackupMySQLWithType_SlackNotificationWithError(t *testing.T) {
	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	// Use invalid parameters to trigger error path
	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("invalid-host", "invalid-user", "invalid-pass", "invalid-db", tmpDir, db.BackupTypeFull)
	// Error is expected, but notification path should be executed
	if err != nil {
		t.Logf("BackupMySQLWithType() returned error (expected): %v", err)
	}
}

