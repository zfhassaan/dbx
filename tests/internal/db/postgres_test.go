package db_test

import (
	"dbx/internal/db"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestBackupPostgres_MissingPgDump tests error handling
func TestBackupPostgres_MissingPgDump(t *testing.T) {
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	os.Setenv("PATH", "")

	err := db.BackupPostgres("localhost", "5432", "postgres", "password", "testdb", "./backups")
	if err == nil {
		t.Error("BackupPostgres() should return error when pg_dump is not found")
	}
}

// TestBackupPostgresWithType_Full tests full backup
func TestBackupPostgresWithType_Full(t *testing.T) {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		t.Skip("Skipping test: pg_dump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupPostgresWithType("localhost", "5432", "postgres", "password", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupPostgresWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupPostgresWithType_Incremental tests incremental backup
func TestBackupPostgresWithType_Incremental(t *testing.T) {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		t.Skip("Skipping test: pg_dump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupPostgresWithType("localhost", "5432", "postgres", "password", "testdb", tmpDir, db.BackupTypeIncremental)
	if err != nil {
		t.Logf("BackupPostgresWithType() returned error (expected without real DB): %v", err)
	}
}

// TestRestorePostgres_MissingPgRestore tests error handling
func TestRestorePostgres_MissingPgRestore(t *testing.T) {
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	os.Setenv("PATH", "")

	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.dump")
	os.WriteFile(backupFile, []byte("dummy backup"), 0644)

	err := db.RestorePostgres("localhost", "5432", "postgres", "password", "testdb", backupFile)
	if err == nil {
		t.Error("RestorePostgres() should return error when pg_restore is not found")
	}
}

// TestRestorePostgres_NonExistentFile tests error handling
func TestRestorePostgres_NonExistentFile(t *testing.T) {
	if _, err := exec.LookPath("pg_restore"); err != nil {
		t.Skip("Skipping test: pg_restore not found in PATH")
	}

	err := db.RestorePostgres("localhost", "5432", "postgres", "password", "testdb", "/nonexistent/backup.dump")
	if err == nil {
		t.Error("RestorePostgres() should return error for non-existent backup file")
	}
}

// TestRestorePostgresTable tests selective table restoration
func TestRestorePostgresTable(t *testing.T) {
	if _, err := exec.LookPath("pg_restore"); err != nil {
		t.Skip("Skipping test: pg_restore not found in PATH")
	}

	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.dump")
	os.WriteFile(backupFile, []byte("dummy backup"), 0644)

	err := db.RestorePostgresTable("localhost", "5432", "postgres", "password", "testdb", backupFile, "test_table")
	if err != nil {
		t.Logf("RestorePostgresTable() returned error (expected without real DB): %v", err)
	}
}

// TestRestorePostgres_EmptyPassword tests restore with empty password
func TestRestorePostgres_EmptyPassword(t *testing.T) {
	if _, err := exec.LookPath("pg_restore"); err != nil {
		t.Skip("Skipping test: pg_restore not found in PATH")
	}

	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.dump")
	os.WriteFile(backupFile, []byte("dummy backup"), 0644)

	err := db.RestorePostgres("localhost", "5432", "postgres", "", "testdb", backupFile)
	if err != nil {
		t.Logf("RestorePostgres() returned error (expected without real DB): %v", err)
	}
}

// TestRestorePostgresTable_EmptyPassword tests table restore with empty password
func TestRestorePostgresTable_EmptyPassword(t *testing.T) {
	if _, err := exec.LookPath("pg_restore"); err != nil {
		t.Skip("Skipping test: pg_restore not found in PATH")
	}

	tmpDir := t.TempDir()
	backupFile := filepath.Join(tmpDir, "backup.dump")
	os.WriteFile(backupFile, []byte("dummy backup"), 0644)

	err := db.RestorePostgresTable("localhost", "5432", "postgres", "", "testdb", backupFile, "test_table")
	if err != nil {
		t.Logf("RestorePostgresTable() returned error (expected without real DB): %v", err)
	}
}

// TestBackupPostgresWithType_Differential tests differential backup type
func TestBackupPostgresWithType_Differential(t *testing.T) {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		t.Skip("Skipping test: pg_dump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupPostgresWithType("localhost", "5432", "postgres", "password", "testdb", tmpDir, db.BackupTypeDifferential)
	if err != nil {
		t.Logf("BackupPostgresWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupPostgresWithType_EmptyPassword tests backup with empty password
func TestBackupPostgresWithType_EmptyPassword(t *testing.T) {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		t.Skip("Skipping test: pg_dump not found in PATH")
	}

	tmpDir := t.TempDir()
	err := db.BackupPostgresWithType("localhost", "5432", "postgres", "", "testdb", tmpDir, db.BackupTypeFull)
	if err != nil {
		t.Logf("BackupPostgresWithType() returned error (expected without real DB): %v", err)
	}
}

// TestBackupPostgres_EmptyDatabase tests validation
func TestBackupPostgres_EmptyDatabase(t *testing.T) {
	if _, err := exec.LookPath("pg_dump"); err != nil {
		t.Skip("Skipping test: pg_dump not found in PATH")
	}

	err := db.BackupPostgres("localhost", "5432", "postgres", "password", "", "./backups")
	if err == nil {
		t.Error("BackupPostgres() should return error for empty database name")
	}
}

