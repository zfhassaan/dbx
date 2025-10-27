package db

import (
	"fmt"
	"os"
	"os/exec"
)

// RestorePostgres restores a PostgreSQL database from a backup file
func RestorePostgres(host, port, user, pass, dbName, backupFile string) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_restore failed: %w", err)
	}

	fmt.Println("✅ Restore completed successfully.")
	return nil
}
