package db

import (
	"bytes"
	"dbx/internal/logs"
	"dbx/internal/notify"
	"fmt"
	"io"
	"os"
	"os/exec"
	osuser "os/user"
	"path/filepath"
	"time"
)

// BackupMySQL creates a backup of a MySQL database
func BackupMySQL(host, user, password, database, outDir string) error {
	return BackupMySQLWithType(host, user, password, database, outDir, BackupTypeFull)
}

// BackupMySQLWithType creates a backup of a MySQL database with specified backup type
func BackupMySQLWithType(host, user, password, database, outDir string, backupType BackupType) error {
	start := time.Now()

	ts := time.Now().Format("2006-01-02_15-04")
	backupSuffix := string(backupType)
	if backupType != BackupTypeFull {
		backupSuffix = string(backupType) + "_" + ts
	}
	outFile := filepath.Join(outDir, fmt.Sprintf("%s-%s_%s.sql", database, backupSuffix, ts))

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	args := []string{"-h", host, "-u", user}
	if password != "" {
		args = append(args, "-p"+password)
	}
	
	// For incremental backups, use --master-data and --flush-logs
	if backupType == BackupTypeIncremental {
		args = append(args, "--master-data=2", "--flush-logs", "--single-transaction")
	} else if backupType == BackupTypeDifferential {
		// Differential backup: backup since last full backup
		args = append(args, "--master-data=2", "--single-transaction")
	}
	
	args = append(args, database)

	cmd := exec.Command("mysqldump", args...)
	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	var stderrBuf bytes.Buffer
	go func() {
		io.Copy(&stderrBuf, stderr)
	}()

	// Start mysqldump process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("mysqldump failed to start: %v", err)
	}

	// Buffer the dump output in memory
	var outputBuf bytes.Buffer
	if _, err := io.Copy(&outputBuf, stdout); err != nil {
		return err
	}

	// Wait for mysqldump to finish and check for error
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("mysqldump failed: %v\n%s", err, stderrBuf.String())
	}

	// Write to .sql file only after successful dump
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("MySQL", "Backup", status, start, err)
		
		// Send Slack notification if webhook is configured
		if webhook := os.Getenv("SLACK_WEBHOOK"); webhook != "" {
			duration := time.Since(start).Round(time.Second)
			hostname, _ := os.Hostname()
			username := "unknown"
			if u, e := osuser.Current(); e == nil {
				username = u.Username
			}
			
			message := fmt.Sprintf("MySQL Backup %s\nDatabase: %s\nDuration: %s\nHost: %s\nUser: %s", 
				status, database, duration, hostname, username)
			if err != nil {
				message += fmt.Sprintf("\nError: %v", err)
			}
			_ = notify.SlackNotify(webhook, message)
		}
	}()

	if _, err := outputBuf.WriteTo(file); err != nil {
		return err
	}

	return nil
}
