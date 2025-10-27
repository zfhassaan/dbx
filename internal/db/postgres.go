package db

import (
	"dbx/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// BackupPostgres runs pg_dump to create a backup of a PostgreSQL database.
func BackupPostgres(host, port, user, pass, dbName, outDir string) error {
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
	outFile := filepath.Join(outDir, fmt.Sprintf("%s_%s.sql", dbName, timestamp))

	// Set env for passwordless execution
	defer func() { _ = os.Setenv("PGPASSWORD", pass) }()
	//os.Setenv("PGPASSWORD", pass)

	cmd := exec.Command("pg_dump",
		"-h", host,
		"-p", port,
		"-U", user,
		"-F", "c", // custom format (compressed binary)
		"-f", outFile,
		dbName,
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	fmt.Println("🔄 Running PostgreSQL backup...")
	if err := cmd.Run(); err != nil {
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
