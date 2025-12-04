package db

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// RestoreSQLite restores a SQLite database from a backup file
func RestoreSQLite(backupFile, targetPath string) error {
	if backupFile == "" {
		return fmt.Errorf("backup file path cannot be empty")
	}

	// Check if backup file exists
	if _, err := os.Stat(backupFile); err != nil {
		return fmt.Errorf("backup file not found: %w", err)
	}

	// If target path is not provided, use the backup file name
	if targetPath == "" {
		targetPath = filepath.Join(filepath.Dir(backupFile), "restored_"+filepath.Base(backupFile))
	}

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Open backup file
	src, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() { _ = src.Close() }()

	// Create target database file
	dst, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target database: %w", err)
	}
	defer func() { _ = dst.Close() }()

	// Copy backup to target
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to restore database: %w", err)
	}

	fmt.Println("✅ SQLite restore completed successfully:", targetPath)
	return nil
}

