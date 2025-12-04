package db

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// RestoreMongo restores a MongoDB database using mongorestore.
func RestoreMongo(uri, dbName, backupDir string) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mongorestore failed: %w", err)
	}

	fmt.Println("✅ MongoDB restore completed successfully.")
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
