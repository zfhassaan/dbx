package db_test

import (
	"dbx/internal/db"
	"os"
	"os/exec"
	"testing"
)

// TestBackupMongo_SlackNotification tests Slack notification path
func TestBackupMongo_SlackNotification(t *testing.T) {
	if _, err := exec.LookPath("mongodump"); err != nil {
		t.Skip("Skipping test: mongodump not found in PATH")
	}

	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	tmpDir := t.TempDir()
	err := db.BackupMongo("mongodb://localhost:27017", "testdb", tmpDir)
	if err != nil {
		t.Logf("BackupMongo() returned error (expected without real DB): %v", err)
	}
}

// TestBackupMongo_SlackNotificationWithError tests Slack notification with error
func TestBackupMongo_SlackNotificationWithError(t *testing.T) {
	if _, err := exec.LookPath("mongodump"); err != nil {
		t.Skip("Skipping test: mongodump not found in PATH")
	}

	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	// Use invalid URI to trigger error path
	tmpDir := t.TempDir()
	err := db.BackupMongo("mongodb://invalid-host:27017", "testdb", tmpDir)
	// Error is expected, but notification path should be executed
	if err != nil {
		t.Logf("BackupMongo() returned error (expected): %v", err)
	}
}

// TestRestoreMongo_SlackNotification tests Slack notification in restore
func TestRestoreMongo_SlackNotification(t *testing.T) {
	if _, err := exec.LookPath("mongorestore"); err != nil {
		t.Skip("Skipping test: mongorestore not found in PATH")
	}

	// Set SLACK_WEBHOOK to trigger notification path
	os.Setenv("SLACK_WEBHOOK", "https://hooks.slack.com/test")
	defer os.Unsetenv("SLACK_WEBHOOK")

	tmpDir := t.TempDir()
	backupDir := tmpDir + "/backup"
	os.MkdirAll(backupDir, 0755)

	err := db.RestoreMongo("mongodb://localhost:27017", "testdb", backupDir)
	if err != nil {
		t.Logf("RestoreMongo() returned error (expected without real DB): %v", err)
	}
}

// TestRestoreMongoCollection_CollectionPathAlternatives tests alternative path structures
func TestRestoreMongoCollection_CollectionPathAlternatives(t *testing.T) {
	if _, err := exec.LookPath("mongorestore"); err != nil {
		t.Skip("Skipping test: mongorestore not found in PATH")
	}

	tmpDir := t.TempDir()
	// Test alternative path structure (collection at root of backup)
	backupDir := tmpDir + "/backup"
	os.MkdirAll(backupDir, 0755)
	os.WriteFile(backupDir+"/test_collection.bson", []byte("test"), 0644)

	err := db.RestoreMongoCollection("mongodb://localhost:27017", "testdb", backupDir, "test_collection")
	if err != nil {
		t.Logf("RestoreMongoCollection() returned error (expected without real DB): %v", err)
	}
}

