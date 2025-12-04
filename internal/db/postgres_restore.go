package db

import (
	"dbx/internal/logs"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// RestorePostgres restores a PostgreSQL database from a backup file
func RestorePostgres(host, port, user, pass, dbName, backupFile string) error {
	start := time.Now()
	if _, err := exec.LookPath("pg_restore"); err != nil {
		showPostgresInstallHelp()
		return fmt.Errorf("pg_restore not found in PATH")
	}

	defer func() { _ = os.Setenv("PGPASSWORD", pass) }()

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

	defer func() { _ = os.Setenv("PGPASSWORD", pass) }()

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
