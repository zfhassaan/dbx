package db

import (
	"dbx/internal/logs"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// RestoreMongo restores a MongoDB database using mongorestore.
func RestoreMongo(uri, dbName, backupDir string) error {
	start := time.Now()
	if dbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	if _, err := exec.LookPath("mongorestore"); err != nil {
		showMongoRestoreHelp()
		return fmt.Errorf("mongorestore not found in PATH")
	}

	cmd := exec.Command("mongorestore",
		"--uri="+uri,
		"--db="+dbName,
		"--drop",
		backupDir,
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	fmt.Println("🔄 Restoring MongoDB database...")
	err := cmd.Run()
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("MongoDB", "Restore", status, start, err)
	}()

	if err != nil {
		return fmt.Errorf("mongorestore failed: %w", err)
	}

	fmt.Println("✅ MongoDB restore completed successfully.")
	return nil
}

// RestoreMongoCollection restores a specific collection from a MongoDB backup
func RestoreMongoCollection(uri, dbName, backupDir, collectionName string) error {
	start := time.Now()
	if dbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	if collectionName == "" {
		return fmt.Errorf("collection name cannot be empty")
	}

	if _, err := exec.LookPath("mongorestore"); err != nil {
		showMongoRestoreHelp()
		return fmt.Errorf("mongorestore not found in PATH")
	}

	// Find the collection directory in the backup
	collectionPath := filepath.Join(backupDir, dbName, collectionName+".bson")
	if _, err := os.Stat(collectionPath); err != nil {
		// Try alternative path structure
		collectionPath = filepath.Join(backupDir, collectionName+".bson")
		if _, err := os.Stat(collectionPath); err != nil {
			return fmt.Errorf("collection '%s' not found in backup directory", collectionName)
		}
	}

	cmd := exec.Command("mongorestore",
		"--uri="+uri,
		"--db="+dbName,
		"--collection="+collectionName,
		"--drop",
		collectionPath,
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	fmt.Printf("🔄 Restoring MongoDB collection '%s'...\n", collectionName)
	err := cmd.Run()
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("MongoDB", fmt.Sprintf("RestoreCollection(%s)", collectionName), status, start, err)
	}()

	if err != nil {
		return fmt.Errorf("mongorestore failed: %w", err)
	}

	fmt.Printf("✅ MongoDB collection '%s' restore completed successfully.\n", collectionName)
	return nil
}

func showMongoRestoreHelp() {
	fmt.Println("\n❌ 'mongorestore' command not found in PATH.")
	switch runtime.GOOS {
	case "linux":
		fmt.Println("👉 Ubuntu / Debian:")
		fmt.Println("   sudo apt install -y mongodb-database-tools")
	case "darwin":
		fmt.Println("👉 macOS:")
		fmt.Println("   brew install mongodb-database-tools")
	case "windows":
		fmt.Println("👉 Windows:")
		fmt.Println("   Download MongoDB Tools and add the bin folder to PATH.")
	default:
		fmt.Println("⚠️ Unsupported OS — please install MongoDB Tools manually.")
	}
}
