package db_test

import (
	"dbx/internal/db"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestBackupMongo_MissingMongodump tests error handling
func TestBackupMongo_MissingMongodump(t *testing.T) {
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	os.Setenv("PATH", "")

	err := db.BackupMongo("mongodb://localhost:27017", "testdb", "./backups")
	if err == nil {
		t.Error("BackupMongo() should return error when mongodump is not found")
	}
}

// TestBackupMongo_EmptyDatabase tests validation
func TestBackupMongo_EmptyDatabase(t *testing.T) {
	if _, err := exec.LookPath("mongodump"); err != nil {
		t.Skip("Skipping test: mongodump not found in PATH")
	}

	err := db.BackupMongo("mongodb://localhost:27017", "", "./backups")
	if err == nil {
		t.Error("BackupMongo() should return error for empty database name")
	}
}

// TestRestoreMongo_MissingMongorestore tests error handling
func TestRestoreMongo_MissingMongorestore(t *testing.T) {
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	os.Setenv("PATH", "")

	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backup")
	os.MkdirAll(backupDir, 0755)

	err := db.RestoreMongo("mongodb://localhost:27017", "testdb", backupDir)
	if err == nil {
		t.Error("RestoreMongo() should return error when mongorestore is not found")
	}
}

// TestRestoreMongo_NonExistentDir tests error handling
func TestRestoreMongo_NonExistentDir(t *testing.T) {
	if _, err := exec.LookPath("mongorestore"); err != nil {
		t.Skip("Skipping test: mongorestore not found in PATH")
	}

	err := db.RestoreMongo("mongodb://localhost:27017", "testdb", "/nonexistent/backup")
	if err == nil {
		t.Error("RestoreMongo() should return error for non-existent backup directory")
	}
}

// TestRestoreMongoCollection tests selective collection restoration
func TestRestoreMongoCollection(t *testing.T) {
	if _, err := exec.LookPath("mongorestore"); err != nil {
		t.Skip("Skipping test: mongorestore not found in PATH")
	}

	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backup")
	os.MkdirAll(backupDir, 0755)

	err := db.RestoreMongoCollection("mongodb://localhost:27017", "testdb", backupDir, "test_collection")
	if err != nil {
		t.Logf("RestoreMongoCollection() returned error (expected without real DB): %v", err)
	}
}

// TestRestoreMongoCollection_EmptyCollection tests error handling
func TestRestoreMongoCollection_EmptyCollection(t *testing.T) {
	if _, err := exec.LookPath("mongorestore"); err != nil {
		t.Skip("Skipping test: mongorestore not found in PATH")
	}

	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backup")
	os.MkdirAll(backupDir, 0755)

	err := db.RestoreMongoCollection("mongodb://localhost:27017", "testdb", backupDir, "")
	if err == nil {
		t.Error("RestoreMongoCollection() should return error for empty collection name")
	}
}

// TestRestoreMongoCollection_EmptyDatabase tests error handling
func TestRestoreMongoCollection_EmptyDatabase(t *testing.T) {
	if _, err := exec.LookPath("mongorestore"); err != nil {
		t.Skip("Skipping test: mongorestore not found in PATH")
	}

	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backup")
	os.MkdirAll(backupDir, 0755)

	err := db.RestoreMongoCollection("mongodb://localhost:27017", "", backupDir, "test_collection")
	if err == nil {
		t.Error("RestoreMongoCollection() should return error for empty database name")
	}
}

// TestRestoreMongo_EmptyDatabase tests error handling
func TestRestoreMongo_EmptyDatabase(t *testing.T) {
	if _, err := exec.LookPath("mongorestore"); err != nil {
		t.Skip("Skipping test: mongorestore not found in PATH")
	}

	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backup")
	os.MkdirAll(backupDir, 0755)

	err := db.RestoreMongo("mongodb://localhost:27017", "", backupDir)
	if err == nil {
		t.Error("RestoreMongo() should return error for empty database name")
	}
}

// TestBackupMongo_EmptyURI tests validation
func TestBackupMongo_EmptyURI(t *testing.T) {
	if _, err := exec.LookPath("mongodump"); err != nil {
		t.Skip("Skipping test: mongodump not found in PATH")
	}

	err := db.BackupMongo("", "testdb", "./backups")
	if err == nil {
		t.Error("BackupMongo() should return error for empty URI")
	}
}

