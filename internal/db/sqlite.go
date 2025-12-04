package db

import (
	"dbx/internal/logs"
	"dbx/internal/notify"
	"dbx/internal/utils"
	"fmt"
	"io"
	"os"
	osuser "os/user"
	"path/filepath"
	"time"
)

// BackupSQLite creates a backup of a SQLite database
func BackupSQLite(dbPath, outDir string) error {
	start := time.Now()

	if dbPath == "" {
		return fmt.Errorf("sqlite database path cannot be empty")
	}

	// Check if database file exists
	if _, err := os.Stat(dbPath); err != nil {
		return fmt.Errorf("sqlite database file not found: %w", err)
	}

	// Prepare output directory
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	dbName := filepath.Base(dbPath)
	dbNameWithoutExt := dbName[:len(dbName)-len(filepath.Ext(dbName))]
	outFile := filepath.Join(outDir, fmt.Sprintf("%s_%s.db", dbNameWithoutExt, timestamp))

	// Open source database
	src, err := os.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer func() { _ = src.Close() }()

	// Create destination database
	dst, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer func() { _ = dst.Close() }()

	// Copy database file
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy database: %w", err)
	}

	// Compress the backup
	zipPath := outFile + ".zip"
	if err := utils.CompressFile(outFile, zipPath); err == nil {
		// Remove uncompressed file after successful compression
		_ = os.Remove(outFile)
		fmt.Println("🗜 Compressed to:", zipPath)
	} else {
		fmt.Println("⚠️ Compression failed, keeping uncompressed backup:", err)
	}

	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("SQLite", "Backup", status, start, err)
		
		// Send Slack notification if webhook is configured
		if webhook := os.Getenv("SLACK_WEBHOOK"); webhook != "" {
			duration := time.Since(start).Round(time.Second)
			hostname, _ := os.Hostname()
			username := "unknown"
			if u, e := osuser.Current(); e == nil {
				username = u.Username
			}
			
			message := fmt.Sprintf("SQLite Backup %s\nDatabase: %s\nDuration: %s\nHost: %s\nUser: %s", 
				status, filepath.Base(dbPath), duration, hostname, username)
			if err != nil {
				message += fmt.Sprintf("\nError: %v", err)
			}
			_ = notify.SlackNotify(webhook, message)
		}
	}()

	fmt.Println("✅ SQLite backup completed:", outFile)
	return nil
}

