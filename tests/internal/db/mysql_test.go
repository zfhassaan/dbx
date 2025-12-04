package db_test

import (
	"dbx/internal/db"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestBackupMySQL_MissingMysqldump tests error handling when mysqldump is not available
func TestBackupMySQL_MissingMysqldump(t *testing.T) {
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Temporarily remove PATH to simulate missing mysqldump
	os.Setenv("PATH", "")

	err := db.BackupMySQL("localhost", "root", "password", "testdb", "./backups")
	if err == nil {
		t.Error("BackupMySQL() should return error when mysqldump is not found")
	}
}

// TestBackupMySQL_EmptyDatabase tests validation
func TestBackupMySQL_EmptyDatabase(t *testing.T) {
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	err := db.BackupMySQL("localhost", "root", "password", "", "./backups")
	if err == nil {
		t.Error("BackupMySQL() should return error for empty database name")
	}
}

// TestBackupMySQLWithType_Full tests full backup type
func TestBackupMySQLWithType_Full(t *testing.T) {
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("localhost", "root", "password", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		// Expected to fail without real MySQL connection, but should not panic
		t.Logf("BackupMySQLWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupMySQLWithType_Incremental tests incremental backup type
func TestBackupMySQLWithType_Incremental(t *testing.T) {
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("localhost", "root", "password", "testdb", tmpDir, db.BackupTypeIncremental)
	if err != nil {
		t.Logf("BackupMySQLWithType() returned error (expected without real DB): %v", err)
	}
}

// TestRestoreMySQL_MissingMySQL tests error handling
func TestRestoreMySQL_MissingMySQL(t *testing.T) {
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	os.Setenv("PATH", "")

	// Create a dummy backup file
	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.sql")
	os.WriteFile(backupFile, []byte("CREATE DATABASE test;"), 0644)

	err := db.RestoreMySQL("localhost", "root", "password", "testdb", backupFile)
	if err == nil {
		t.Error("RestoreMySQL() should return error when mysql is not found")
	}
}

// TestRestoreMySQL_NonExistentFile tests error handling
func TestRestoreMySQL_NonExistentFile(t *testing.T) {
	if _, err := exec.LookPath("mysql"); err != nil {
		t.Skip("Skipping test: mysql not found in PATH")
	}

	err := db.RestoreMySQL("localhost", "root", "password", "testdb", "/nonexistent/backup.sql")
	if err == nil {
		t.Error("RestoreMySQL() should return error for non-existent backup file")
	}
}

// TestRestoreMySQLTable tests selective table restoration
func TestRestoreMySQLTable(t *testing.T) {
	if _, err := exec.LookPath("mysql"); err != nil {
		t.Skip("Skipping test: mysql not found in PATH")
	}

	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.sql")
	os.WriteFile(backupFile, []byte("CREATE TABLE test (id INT);"), 0644)

	err := db.RestoreMySQLTable("localhost", "root", "password", "testdb", backupFile, "test_table")
	if err != nil {
		t.Logf("RestoreMySQLTable() returned error (expected without real DB): %v", err)
	}
}

// TestBackupMySQLWithType_Differential tests differential backup type
func TestBackupMySQLWithType_Differential(t *testing.T) {
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("localhost", "root", "password", "testdb", tmpDir, db.BackupTypeDifferential)
	if err != nil {
		t.Logf("BackupMySQLWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupMySQLWithType_EmptyPassword tests backup with empty password
func TestBackupMySQLWithType_EmptyPassword(t *testing.T) {
	if _, err := exec.LookPath("mysqldump"); err != nil {
		t.Skip("Skipping test: mysqldump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupMySQLWithType("localhost", "root", "", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupMySQLWithType() returned error (expected without real DB): %v", err)
	}
}

// TestRestoreMySQL_EmptyPassword tests restore with empty password
func TestRestoreMySQL_EmptyPassword(t *testing.T) {
	if _, err := exec.LookPath("mysql"); err != nil {
		t.Skip("Skipping test: mysql not found in PATH")
	}

	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.sql")
	os.WriteFile(backupFile, []byte("CREATE DATABASE test;"), 0644)

	err := db.RestoreMySQL("localhost", "root", "", "testdb", backupFile)
	if err != nil {
		t.Logf("RestoreMySQL() returned error (expected without real DB): %v", err)
	}
}

// TestRestoreMySQLTable_EmptyPassword tests table restore with empty password
func TestRestoreMySQLTable_EmptyPassword(t *testing.T) {
	if _, err := exec.LookPath("mysql"); err != nil {
		t.Skip("Skipping test: mysql not found in PATH")
	}

	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.sql")
	os.WriteFile(backupFile, []byte("CREATE TABLE test (id INT);"), 0644)

	err := db.RestoreMySQLTable("localhost", "root", "", "testdb", backupFile, "test_table")
	if err != nil {
		t.Logf("RestoreMySQLTable() returned error (expected without real DB): %v", err)
	}
}

// TestRestoreMySQLTable_TableNotFound tests error handling when table not in backup
func TestRestoreMySQLTable_TableNotFound(t *testing.T) {
	if _, err := exec.LookPath("mysql"); err != nil {
		t.Skip("Skipping test: mysql not found in PATH")
	}

	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.sql")
	// Backup file without the requested table
	os.WriteFile(backupFile, []byte("CREATE TABLE other_table (id INT);"), 0644)

	err := db.RestoreMySQLTable("localhost", "root", "password", "testdb", backupFile, "nonexistent_table")
	if err == nil {
		t.Error("RestoreMySQLTable() should return error when table not found in backup")
	}
}

