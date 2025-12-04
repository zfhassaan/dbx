package db

import (
	"bufio"
	"dbx/internal/logs"
	"dbx/internal/notify"
	"dbx/internal/utils"
	"fmt"
	"os"
	"os/exec"
	osuser "os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// BackupMongo runs mongodump to create a backup
func BackupMongo(uri, dbName, outDir string) error {
	if dbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	// Check if mongodump is available
	if _, err := exec.LookPath("mongodump"); err != nil {
		fmt.Println("\n❌ 'mongodump' not found in PATH.")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Would you like DBX to install MongoDB Tools for you? (y/N): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		if choice == "y" || choice == "yes" {
			if err := installMongoTools(); err != nil {
				fmt.Println("❌ Automatic installation failed:", err)
			} else {
				fmt.Println("✅ MongoDB Tools installed successfully!")
			}
		} else {
			showMongoInstallHelp()
		}

		// Recheck after installation
		if _, err := exec.LookPath("mongodump"); err != nil {
			return fmt.Errorf("mongodump still not found, aborting backup")
		}
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Timestamped folder
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	outPath := filepath.Join(outDir, fmt.Sprintf("%s_%s", dbName, timestamp))

	cmd := exec.Command("mongodump",
		"--uri="+uri,
		"--db="+dbName,
		"--out="+outPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	fmt.Println("🔄 Running MongoDB backup...")
	err := cmd.Run()
	
	defer func() {
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
		}
		logs.LogEntry("MongoDB", "Backup", status, start, err)
		
		// Send Slack notification if webhook is configured
		if webhook := os.Getenv("SLACK_WEBHOOK"); webhook != "" {
			duration := time.Since(start).Round(time.Second)
			hostname, _ := os.Hostname()
			username := "unknown"
			if u, e := osuser.Current(); e == nil {
				username = u.Username
			}
			
			message := fmt.Sprintf("MongoDB Backup %s\nDatabase: %s\nDuration: %s\nHost: %s\nUser: %s", 
				status, dbName, duration, hostname, username)
			if err != nil {
				message += fmt.Sprintf("\nError: %v", err)
			}
			_ = notify.SlackNotify(webhook, message)
		}
	}()
	
	if err != nil {
		return fmt.Errorf("mongodump failed: %w", err)
	}

	zipPath := outPath + ".zip"
	if err := utils.CompressFolder(outPath, zipPath); err == nil {
		defer func() { _ = os.RemoveAll(outPath) }()
		fmt.Println("🗜 Compressed to:", zipPath)
	} else {
		fmt.Println("⚠️ Compression failed:", err)
	}

	fmt.Println("✅ Backup completed:", outPath)
	return nil
}

// installMongoTools tries to install MongoDB Database Tools based on OS
func installMongoTools() error {
	switch runtime.GOOS {
	case "linux":
		if isCommandAvailable("apt") {
			fmt.Println("📦 Installing via apt (Ubuntu/Debian)...")
			cmd := exec.Command("sudo", "apt", "update")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
			cmd = exec.Command("sudo", "apt", "install", "-y", "mongodb-database-tools")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			return cmd.Run()
		} else if isCommandAvailable("yum") {
			fmt.Println("📦 Installing via yum (CentOS/RHEL)...")
			cmd := exec.Command("sudo", "yum", "install", "-y", "mongodb-database-tools")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			return cmd.Run()
		}
	case "darwin":
		fmt.Println("📦 Installing via Homebrew (macOS)...")
		cmd := exec.Command("brew", "install", "mongodb-database-tools")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		return cmd.Run()
	case "windows":
		fmt.Println("⚠️ Windows auto-install is not supported.")
		fmt.Println("Opening MongoDB Tools download page...")
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler",
			"https://www.mongodb.com/try/download/database-tools").Start()
		return nil
	}

	return fmt.Errorf("unsupported OS for automatic installation")
}

func showMongoInstallHelp() {
	fmt.Println("\n❌ 'mongodump' command not found in your system PATH.")
	switch runtime.GOOS {
	case "linux":
		fmt.Println("🔧 To install MongoDB Database Tools:")
		fmt.Println("👉 On Ubuntu / Debian:")
		fmt.Println("   sudo apt update && sudo apt install -y mongodb-database-tools")
		fmt.Println("👉 On CentOS / RHEL:")
		fmt.Println("   sudo yum install -y mongodb-database-tools")
	case "darwin":
		fmt.Println("🔧 To install MongoDB Tools on macOS:")
		fmt.Println("   brew install mongodb-database-tools")
	case "windows":
		fmt.Println("🔧 To install MongoDB Tools on Windows:")
		fmt.Println("1️⃣ Go to: https://www.mongodb.com/try/download/database-tools")
		fmt.Println("2️⃣ Extract the ZIP (e.g., C:\\Program Files\\MongoDB\\DatabaseTools\\bin)")
		fmt.Println("3️⃣ Add that folder to your PATH environment variable.")
	default:
		fmt.Println("⚠️ Unsupported OS — please manually install MongoDB Database Tools.")
	}
	fmt.Println("Once installed, re-run DBX to use the MongoDB backup feature.")
}

// helper to check command availability
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func RunMongobackup() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--- MongoDB Backup ---")

	fmt.Print("MongoDB URI [mongodb://localhost:27017]: ")
	uri, _ := reader.ReadString('\n')
	uri = strings.TrimSpace(uri)
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	fmt.Print("Database Name: ")
	dbName, _ := reader.ReadString('\n')
	dbName = strings.TrimSpace(dbName)

	fmt.Print("Backup Directory [./backups]: ")
	out, _ := reader.ReadString('\n')
	out = strings.TrimSpace(out)
	if out == "" {
		out = "./backups"
	}

	if err := BackupMongo(uri, dbName, out); err != nil {
		fmt.Println("❌ Backup failed:", err)
	} else {
		fmt.Println("✅ Backup successful!")
	}
}
