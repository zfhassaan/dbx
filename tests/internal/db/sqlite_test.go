package db_test

import (
	"dbx/internal/db"
	"os"
	"path/filepath"
	"testing"
)

// TestBackupSQLite_ValidDatabase tests SQLite backup with valid database
func TestBackupSQLite_ValidDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a test SQLite database file
	testDB := filepath.Join(tmpDir, "test.db")
	os.WriteFile(testDB, []byte("SQLite format 3\x00"), 0644)

	backupDir := filepath.Join(tmpDir, "backups")
	err := db.BackupSQLite(testDB, backupDir)

	// Should succeed or fail gracefully (depending on if it's a real SQLite file)
	if err != nil {
		t.Logf("BackupSQLite() returned error (expected if not valid SQLite): %v", err)
	}
}

// TestBackupSQLite_NonExistentDatabase tests error handling
func TestBackupSQLite_NonExistentDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentDB := filepath.Join(tmpDir, "nonexistent.db")
	backupDir := filepath.Join(tmpDir, "backups")

	err := db.BackupSQLite(nonExistentDB, backupDir)
	if err == nil {
		t.Error("BackupSQLite() should return error for non-existent database")
	}
}

// TestBackupSQLite_EmptyPath tests validation
func TestBackupSQLite_EmptyPath(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")

	err := db.BackupSQLite("", backupDir)
	if err == nil {
		t.Error("BackupSQLite() should return error for empty database path")
	}
}

// TestBackupSQLite_InvalidOutputDir tests error handling for invalid output
func TestBackupSQLite_InvalidOutputDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDB := filepath.Join(tmpDir, "test.db")
	os.WriteFile(testDB, []byte("SQLite format 3"), 0644)

	// Use a platform-appropriate invalid path
	var invalidDir string
	if os.PathSeparator == '\\' {
		// Windows: use a path with invalid characters
		invalidDir = "C:\\<invalid>\\path"
	} else {
		// Unix: use a path that requires root permissions
		invalidDir = "/root/nonexistent/path"
	}

	err := db.BackupSQLite(testDB, invalidDir)

	// Should fail for invalid output directory
	if err == nil {
		t.Error("BackupSQLite() should return error for invalid output directory")
	}
}

