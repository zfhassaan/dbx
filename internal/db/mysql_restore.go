package db

import (
	"bufio"
	"dbx/internal/logs"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// RestoreMySQL restores a MySQL database from a .sql dump file
func RestoreMySQL(host, user, pass, dbName, backupFile string) error {
	start := time.Now()
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

	err = cmd.Run()
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("MySQL", "Restore", status, start, err)
	}()

	if err != nil {
		return fmt.Errorf("mysql restore failed: %w", err)
	}

	fmt.Println("✅ MySQL restore completed successfully.")
	return nil
}

// RestoreMySQLTable restores a specific table from a MySQL dump file
func RestoreMySQLTable(host, user, pass, dbName, backupFile, tableName string) error {
	start := time.Now()
	if _, err := exec.LookPath("mysql"); err != nil {
		return fmt.Errorf("mysql not found in PATH")
	}

	fmt.Printf("🔄 Restoring MySQL table '%s'...\n", tableName)

	// Extract table data from dump file
	file, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	var tableData strings.Builder
	inTable := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, fmt.Sprintf("CREATE TABLE `%s`", tableName)) ||
			strings.Contains(line, fmt.Sprintf("CREATE TABLE %s", tableName)) {
			inTable = true
			tableData.WriteString(line + "\n")
		} else if inTable {
			tableData.WriteString(line + "\n")
			if strings.HasPrefix(strings.TrimSpace(line), "UNLOCK TABLES") ||
				strings.HasPrefix(strings.TrimSpace(line), "LOCK TABLES") {
				if strings.Contains(line, tableName) {
					continue
				}
				break
			}
		}
	}

	if tableData.Len() == 0 {
		return fmt.Errorf("table '%s' not found in backup file", tableName)
	}

	cmd := exec.Command("mysql",
		"-h", host,
		"-u", user,
		"-p"+pass,
		dbName,
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.Stdin = strings.NewReader(tableData.String())

	err = cmd.Run()
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("MySQL", fmt.Sprintf("RestoreTable(%s)", tableName), status, start, err)
	}()

	if err != nil {
		return fmt.Errorf("mysql table restore failed: %w", err)
	}

	fmt.Printf("✅ MySQL table '%s' restore completed successfully.\n", tableName)
	return nil
}

