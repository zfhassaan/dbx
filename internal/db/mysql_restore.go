package db

import (
	"fmt"
	"os"
	"os/exec"
)

// RestoreMySQL restores a MySQL database from a .sql dump file
func RestoreMySQL(host, user, pass, dbName, backupFile string) error {
	if _, err := exec.LookPath("mysql"); err != nil {
		fmt.Println("❌ 'mysql' command not found in PATH.")
		fmt.Println("👉 Install MySQL client tools:")
		fmt.Println("   sudo apt install -y mysql-client")
		return fmt.Errorf("mysql not found in PATH")
	}

	fmt.Println("🔄 Restoring MySQL database...")

	cmd := exec.Command("mysql",
		"-h", host,
		"-u", user,
		"-p"+pass,
		dbName,
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	// Pipe backup file into mysql command
	file, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() { _ = file.Close() }()

	cmd.Stdin = file

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql restore failed: %w", err)
	}

	fmt.Println("✅ MySQL restore completed successfully.")
	return nil
}
