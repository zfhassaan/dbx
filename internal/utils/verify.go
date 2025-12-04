package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// VerifyBackupFile verifies a backup file exists and is readable
func VerifyBackupFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("backup file not found or inaccessible: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("backup file is empty")
	}
	return nil
}

// CalculateChecksum calculates SHA256 checksum of a file
func CalculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// VerifyBackupIntegrity verifies backup file exists and optionally calculates checksum
func VerifyBackupIntegrity(filePath string, calculateChecksum bool) (string, error) {
	if err := VerifyBackupFile(filePath); err != nil {
		return "", err
	}

	if calculateChecksum {
		checksum, err := CalculateChecksum(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to calculate checksum: %w", err)
		}
		return checksum, nil
	}

	return "", nil
}

