package db

import (
	"dbx/internal/logs"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// RestorePostgres restores a PostgreSQL database from a backup file
func RestorePostgres(host, port, user, pass, dbName, backupFile string) error {
	start := time.Now()
	if _, err := exec.LookPath("pg_restore"); err != nil {
		showPostgresInstallHelp()
		return fmt.Errorf("pg_restore not found in PATH")
	}

	// Set env for passwordless execution - must be set BEFORE cmd.Run()
	if pass != "" {
		os.Setenv("PGPASSWORD", pass)
	} else {
		os.Unsetenv("PGPASSWORD")
	}

	cmd := exec.Command("pg_restore",
		"-h", host,
		"-p", port,
		"-U", user,
		"-d", dbName,
		"-c", // clean before restore
		backupFile,
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	fmt.Println("🔄 Restoring PostgreSQL database...")
	err := cmd.Run()
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("PostgreSQL", "Restore", status, start, err)
	}()

	if err != nil {
		return fmt.Errorf("pg_restore failed: %w", err)
	}

	fmt.Println("✅ Restore completed successfully.")
	return nil
}

// RestorePostgresTable restores a specific table from a PostgreSQL backup
func RestorePostgresTable(host, port, user, pass, dbName, backupFile, tableName string) error {
	start := time.Now()
	if _, err := exec.LookPath("pg_restore"); err != nil {
		showPostgresInstallHelp()
		return fmt.Errorf("pg_restore not found in PATH")
	}

	// Verify backup file exists
	if _, err := os.Stat(backupFile); err != nil {
		return fmt.Errorf("backup file not found: %w", err)
	}

	// Verify table exists in backup using pg_restore --list
	// Set env for passwordless execution - must be set BEFORE cmd.Run()
	if pass != "" {
		os.Setenv("PGPASSWORD", pass)
	} else {
		os.Unsetenv("PGPASSWORD")
	}

	// List contents of backup to verify table exists
	listCmd := exec.Command("pg_restore", "--list", backupFile)
	listOutput, listErr := listCmd.Output()
	if listErr == nil {
		// Check if table name appears in the backup contents
		if !strings.Contains(string(listOutput), tableName) {
			return fmt.Errorf("table '%s' not found in backup file. Use 'pg_restore --list %s' to see available tables", tableName, backupFile)
		}
	}
	// If list command fails, continue anyway (backup might be in different format)

	cmd := exec.Command("pg_restore",
		"-h", host,
		"-p", port,
		"-U", user,
		"-d", dbName,
		"-t", tableName, // restore specific table
		"-c",            // clean before restore
		backupFile,
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	fmt.Printf("🔄 Restoring PostgreSQL table '%s'...\n", tableName)
	err := cmd.Run()
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("PostgreSQL", fmt.Sprintf("RestoreTable(%s)", tableName), status, start, err)
	}()

	if err != nil {
		return fmt.Errorf("pg_restore failed: %w", err)
	}

	fmt.Printf("✅ PostgreSQL table '%s' restore completed successfully.\n", tableName)
	return nil
}
