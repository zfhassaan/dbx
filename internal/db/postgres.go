package db

import (
	"dbx/internal/logs"
	"dbx/internal/notify"
	"dbx/internal/utils"
	"fmt"
	"os"
	"os/exec"
	osuser "os/user"
	"path/filepath"
	"time"
)

// BackupPostgres runs pg_dump to create a backup of a PostgreSQL database.
func BackupPostgres(host, port, user, pass, dbName, outDir string) error {
	return BackupPostgresWithType(host, port, user, pass, dbName, outDir, BackupTypeFull)
}

// BackupPostgresWithType runs pg_dump to create a backup with specified type
func BackupPostgresWithType(host, port, user, pass, dbName, outDir string, backupType BackupType) error {
	if dbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	if _, err := exec.LookPath("pg_dump"); err != nil {
		showPostgresInstallHelp()
		return fmt.Errorf("pg_dump not found in PATH")
	}

	// Prepare directory
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupSuffix := string(backupType)
	if backupType != BackupTypeFull {
		backupSuffix = string(backupType) + "_" + timestamp
	}
	outFile := filepath.Join(outDir, fmt.Sprintf("%s_%s_%s.sql", dbName, backupSuffix, timestamp))

	// Set env for passwordless execution
	defer func() { _ = os.Setenv("PGPASSWORD", pass) }()
	//os.Setenv("PGPASSWORD", pass)

	args := []string{
		"-h", host,
		"-p", port,
		"-U", user,
		"-F", "c", // custom format (compressed binary)
		"-f", outFile,
	}
	
	// For incremental backups, use WAL-based approach
	// Note: True incremental requires continuous archiving setup
	if backupType == BackupTypeIncremental || backupType == BackupTypeDifferential {
		args = append(args, "--verbose")
	}
	
	args = append(args, dbName)
	cmd := exec.Command("pg_dump", args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	start := time.Now()
	fmt.Println("🔄 Running PostgreSQL backup...")
	err := cmd.Run()
	
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("PostgreSQL", "Backup", status, start, err)
		
		// Send Slack notification if webhook is configured
		if webhook := os.Getenv("SLACK_WEBHOOK"); webhook != "" {
			duration := time.Since(start).Round(time.Second)
			hostname, _ := os.Hostname()
			username := "unknown"
			if u, e := osuser.Current(); e == nil {
				username = u.Username
			}
			
			message := fmt.Sprintf("PostgreSQL Backup %s\nDatabase: %s\nDuration: %s\nHost: %s\nUser: %s", 
				status, dbName, duration, hostname, username)
			if err != nil {
				message += fmt.Sprintf("\nError: %v", err)
			}
			_ = notify.SlackNotify(webhook, message)
		}
	}()
	
	if err != nil {
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	fmt.Println("✅ Backup completed:", outFile)

	// Optional: compress final .sql/.dump
	zipPath := outFile + ".zip"
	if err := utils.CompressFolder(outDir, zipPath); err == nil {
		fmt.Println("🗜 Compressed to:", zipPath)
	}

	return nil
}

// showPostgresInstallHelp prints guidance if pg_dump missing
func showPostgresInstallHelp() {
	fmt.Println("\n❌ 'pg_dump' not found.")
	fmt.Println("👉 Ubuntu / Debian:")
	fmt.Println("   sudo apt install -y postgresql-client")
	fmt.Println("👉 CentOS / RHEL:")
	fmt.Println("   sudo yum install -y postgresql")
	fmt.Println("👉 macOS:")
	fmt.Println("   brew install postgresql")
	fmt.Println("👉 Windows:")
	fmt.Println("   Install PostgreSQL and ensure its /bin folder is in PATH.")
}
