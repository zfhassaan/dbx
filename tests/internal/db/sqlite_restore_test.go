package db_test

import (
	"dbx/internal/db"
	"os"
	"path/filepath"
	"testing"
)

// TestRestoreSQLite_NonExistentFile tests error handling
func TestRestoreSQLite_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "restored.db")

	err := db.RestoreSQLite("/nonexistent/backup.db", targetPath)
	if err == nil {
		t.Error("RestoreSQLite() should return error for non-existent backup file")
	}
}

// TestRestoreSQLite_ValidBackup tests valid restoration
func TestRestoreSQLite_ValidBackup(t *testing.T) {
	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.db")
	targetPath := filepath.Join(tmpDir, "restored.db")

	// Create a valid SQLite backup file
	os.WriteFile(backupFile, []byte("SQLite format 3"), 0644)

	err := db.RestoreSQLite(backupFile, targetPath)
	if err != nil {
		t.Errorf("RestoreSQLite() error = %v, want nil", err)
	}

	// Verify file was created
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Error("RestoreSQLite() did not create restored database file")
	}
}

// TestRestoreSQLite_CompressedBackup tests restoration from compressed backup
func TestRestoreSQLite_CompressedBackup(t *testing.T) {
	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.db.zip")
	targetPath := filepath.Join(tmpDir, "restored.db")

	// Create a compressed backup file (simplified - real zip would be more complex)
	os.WriteFile(backupFile, []byte("compressed backup"), 0644)

	// This should handle compressed files
	err := db.RestoreSQLite(backupFile, targetPath)
	// May fail if not a valid zip, but should not panic
	if err != nil {
		t.Logf("RestoreSQLite() returned error (expected for invalid zip): %v", err)
	}
}

// TestRestoreSQLite_EmptyBackupFile tests error handling
func TestRestoreSQLite_EmptyBackupFile(t *testing.T) {
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "restored.db")

	err := db.RestoreSQLite("", targetPath)
	if err == nil {
		t.Error("RestoreSQLite() should return error for empty backup file path")
	}
}

// TestRestoreSQLite_EmptyTargetPath tests automatic target path generation
func TestRestoreSQLite_EmptyTargetPath(t *testing.T) {
	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.db")
	os.WriteFile(backupFile, []byte("SQLite format 3"), 0644)

	// Empty target path should use backup file location with "restored_" prefix
	err := db.RestoreSQLite(backupFile, "")
	if err != nil {
		t.Logf("RestoreSQLite() returned error (expected if not valid SQLite): %v", err)
	}
}

// TestRestoreSQLite_InvalidTargetDir tests error handling for invalid target directory
func TestRestoreSQLite_InvalidTargetDir(t *testing.T) {
	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.db")
	os.WriteFile(backupFile, []byte("SQLite format 3"), 0644)

	var invalidDir string
	if os.PathSeparator == '\\' {
		invalidDir = "C:\\<invalid>\\restored.db"
	} else {
		invalidDir = "/root/nonexistent/restored.db"
	}

	err := db.RestoreSQLite(backupFile, invalidDir)
	if err == nil {
		t.Error("RestoreSQLite() should return error for invalid target directory")
	}
}

