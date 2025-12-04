package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupType represents the type of backup
type BackupType string

const (
	BackupTypeFull        BackupType = "full"
	BackupTypeIncremental BackupType = "incremental"
	BackupTypeDifferential BackupType = "differential"
)

// BackupMetadata stores information about backups for incremental/differential tracking
type BackupMetadata struct {
	LastFullBackup     time.Time `json:"last_full_backup"`
	LastIncrementalBackup time.Time `json:"last_incremental_backup"`
	BackupPath         string    `json:"backup_path"`
	DBType             string    `json:"db_type"`
	Database           string    `json:"database"`
}

// GetMetadataPath returns the path to the metadata file for a database
func GetMetadataPath(outDir, dbType, dbName string) string {
	return filepath.Join(outDir, fmt.Sprintf(".%s_%s_metadata.json", dbType, dbName))
}

// LoadMetadata loads backup metadata from disk
func LoadMetadata(metadataPath string) (*BackupMetadata, error) {
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return &BackupMetadata{}, nil
	}

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata BackupMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &metadata, nil
}

// SaveMetadata saves backup metadata to disk
func SaveMetadata(metadataPath string, metadata *BackupMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

