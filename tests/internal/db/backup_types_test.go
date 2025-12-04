package db_test

import (
	"dbx/internal/db"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestBackupType_Constants tests backup type constants
func TestBackupType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		backupType db.BackupType
		want     string
	}{
		{"Full backup", db.BackupTypeFull, "full"},
		{"Incremental backup", db.BackupTypeIncremental, "incremental"},
		{"Differential backup", db.BackupTypeDifferential, "differential"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.backupType) != tt.want {
				t.Errorf("BackupType = %v, want %v", tt.backupType, tt.want)
			}
		})
	}
}

// TestGetMetadataPath tests metadata path generation
func TestGetMetadataPath(t *testing.T) {
	tests := []struct {
		name   string
		outDir string
		dbType string
		dbName string
		want   string
	}{
		{
			name:   "MySQL metadata path",
			outDir: "./backups",
			dbType: "mysql",
			dbName: "testdb",
			want:   filepath.Join(".", "backups", ".mysql_testdb_metadata.json"),
		},
		{
			name:   "PostgreSQL metadata path",
			outDir: "/var/backups",
			dbType: "postgres",
			dbName: "mydb",
			want:   filepath.Join("/var", "backups", ".postgres_mydb_metadata.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := db.GetMetadataPath(tt.outDir, tt.dbType, tt.dbName)
			if got != tt.want {
				t.Errorf("GetMetadataPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLoadMetadata_NonExistent tests loading non-existent metadata
func TestLoadMetadata_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "nonexistent.json")

	metadata, err := db.LoadMetadata(metadataPath)
	if err != nil {
		t.Errorf("LoadMetadata() should return empty metadata for non-existent file, got error: %v", err)
	}

	if metadata == nil {
		t.Error("LoadMetadata() should return empty metadata struct, got nil")
	}
}

// TestSaveAndLoadMetadata tests save and load cycle
func TestSaveAndLoadMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	original := &db.BackupMetadata{
		LastFullBackup:        time.Now(),
		LastIncrementalBackup: time.Now().Add(-1 * time.Hour),
		BackupPath:            "/backups/test",
		DBType:                "mysql",
		Database:              "testdb",
	}

	// Save metadata
	if err := db.SaveMetadata(metadataPath, original); err != nil {
		t.Fatalf("SaveMetadata() error = %v", err)
	}

	// Load metadata
	loaded, err := db.LoadMetadata(metadataPath)
	if err != nil {
		t.Fatalf("LoadMetadata() error = %v", err)
	}

	// Verify fields
	if loaded.DBType != original.DBType {
		t.Errorf("LoadMetadata() DBType = %v, want %v", loaded.DBType, original.DBType)
	}
	if loaded.Database != original.Database {
		t.Errorf("LoadMetadata() Database = %v, want %v", loaded.Database, original.Database)
	}
	if loaded.BackupPath != original.BackupPath {
		t.Errorf("LoadMetadata() BackupPath = %v, want %v", loaded.BackupPath, original.BackupPath)
	}
}

// TestSaveMetadata_InvalidPath tests error handling for invalid path
func TestSaveMetadata_InvalidPath(t *testing.T) {
	// Use a platform-appropriate invalid path
	var invalidPath string
	if os.PathSeparator == '\\' {
		// Windows: use a path with invalid characters
		invalidPath = "C:\\<invalid>\\metadata.json"
	} else {
		// Unix: use a path that requires root permissions
		invalidPath = "/root/nonexistent/metadata.json"
	}
	metadata := &db.BackupMetadata{}

	err := db.SaveMetadata(invalidPath, metadata)
	if err == nil {
		// On some systems, the parent directory might be created automatically
		// So we check if the file was actually created at the invalid path
		if _, statErr := os.Stat(invalidPath); statErr == nil {
			t.Error("SaveMetadata() should not create file at invalid path")
		} else {
			// File wasn't created at invalid path, which is acceptable
			t.Log("SaveMetadata() did not create file at invalid path (acceptable behavior)")
		}
	}
	// If error occurred, that's the expected behavior
}

// TestLoadMetadata_CorruptedFile tests handling of corrupted JSON
func TestLoadMetadata_CorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "corrupted.json")

	// Write invalid JSON
	os.WriteFile(metadataPath, []byte("invalid json content"), 0644)

	_, err := db.LoadMetadata(metadataPath)
	if err == nil {
		t.Error("LoadMetadata() should return error for corrupted JSON file")
	}
}

// TestBackupMetadata_Serialization tests JSON serialization/deserialization
func TestBackupMetadata_Serialization(t *testing.T) {
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "serialization_test.json")

	now := time.Now()
	original := &db.BackupMetadata{
		LastFullBackup:        now,
		LastIncrementalBackup: now.Add(-2 * time.Hour),
		BackupPath:            "/backups/db",
		DBType:                "postgres",
		Database:              "production",
	}

	// Save
	if err := db.SaveMetadata(metadataPath, original); err != nil {
		t.Fatalf("SaveMetadata() failed: %v", err)
	}

	// Load
	loaded, err := db.LoadMetadata(metadataPath)
	if err != nil {
		t.Fatalf("LoadMetadata() failed: %v", err)
	}

	// Verify time fields are preserved (within 1 second tolerance)
	if loaded.LastFullBackup.Sub(original.LastFullBackup) > time.Second {
		t.Errorf("LastFullBackup not preserved correctly")
	}
	if loaded.LastIncrementalBackup.Sub(original.LastIncrementalBackup) > time.Second {
		t.Errorf("LastIncrementalBackup not preserved correctly")
	}
}

