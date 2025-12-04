package main

import (
	"bufio"
	"dbx/internal/cloud"
	"dbx/internal/db"
	"dbx/internal/scheduler"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/term"
)

type App struct {
	reader *bufio.Reader
}

func main() {
	app := &App{
		reader: bufio.NewReader(os.Stdin),
	}
	app.MainMenu()
}

// clearScreen clears the terminal screen using ANSI escape codes and flushes stdout.
func (a *App) clearScreen() {
	switch runtime.GOOS {
	case "windows":
		for i := 0; i < 40; i++ {
			fmt.Println()
		}
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			// Silently ignore clear screen errors - not critical for functionality
			return
		}
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			// Silently ignore clear screen errors - not critical for functionality
			return
		}
	}
}

func (a *App) showBanner() {
	fmt.Println(`
  ____  ____  __  __       
 |  _ \| __ )|  \/  | ___  
 | | | |  _ \| |\/| |/ _ \ 
 | |_| | |_) | |  | | (_) |
 |____/|____/|_|  |_|\___/  v0.2.0

 DBX — Dead-simple database backups.
 -----------------------------------
    Author: @zfhassaan
 ----------------------------------`)
}

func (a *App) promptInput(prompt, defaultVal string, hideInput bool) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	var input string
	if hideInput {
		// Use terminal password input (hides characters as user types)
		// Works cross-platform: Windows, Linux, macOS
		fd := int(os.Stdin.Fd())
		bytePassword, err := term.ReadPassword(fd)
		if err != nil {
			// Fallback to normal input if password reading fails
			fmt.Println("\n⚠️ Warning: Could not hide input, using normal input mode")
			input, _ = a.reader.ReadString('\n')
		} else {
			input = string(bytePassword)
			fmt.Println() // Add newline after hidden input
		}
	} else {
		// Ignore ReadString error - input will be empty string on error, which is handled below
		input, _ = a.reader.ReadString('\n')
	}

	input = strings.TrimSpace(input)
	if input == "" && defaultVal != "" {
		return defaultVal
	}
	return input
}

func (a *App) MainMenu() {
	scheduler.Init()
	a.clearScreen()
	a.showBanner()
	fmt.Println("===== DBX: Database Backup Utility =====")
	fmt.Println("[1] 🔄 Backup Menu")
	fmt.Println("[2] 🔁 Restore Menu")
	fmt.Println("[3] 📜 View Logs")
	fmt.Println("[4] 🧩 Test Database Connection")
	fmt.Println("[5] 🕒 Schedule Backups")
	fmt.Println("[6] ☁️ Cloud Storage Help")
	fmt.Println("[0] ❌  Exit")
	fmt.Print("Enter your choice: ")

	choice := a.readInt()

	switch choice {
	case 1:
		a.BackupMenu()
	case 2:
		a.RestoreMenu()
	case 3:
		a.ViewLogs()
	case 4:
		a.TestConnectionMenu()
	case 5:
		a.ScheduleMenu()
	case 6:
		a.CloudHelp()
	case 0:
		fmt.Println("👋 Exiting DBX.")
		os.Exit(0)
	default:
		fmt.Println("⚠️ Invalid option. Try again.")
		a.MainMenu()
	}
}

func (a *App) TestConnectionMenu() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("--- Test Database Connection ---")
	fmt.Println("[1] MySQL")
	fmt.Println("[2] PostgreSQL")
	fmt.Println("[3] MongoDB")
	fmt.Println("[4] SQLite")
	fmt.Println("[0] Back to Main Menu")
	fmt.Print("Choose database: ")

	switch a.readInt() {
	case 1:
		a.TestMySQLConnection()
	case 2:
		a.TestPostgresConnection()
	case 3:
		a.TestMongoConnection()
	case 4:
		a.TestSQLiteConnection()
	case 0:
		a.MainMenu()
	default:
		fmt.Println("Invalid choice.")
		a.TestConnectionMenu()
	}
}

func (a *App) CloudHelp() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("--- Cloud Storage Setup ---")
	fmt.Println()
	fmt.Println("🌐 AWS S3:")
	fmt.Println("1️⃣ Install AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html")
	fmt.Println("2️⃣ Run: aws configure")
	fmt.Println("   Enter your Access Key, Secret, and Region.")
	fmt.Println("3️⃣ DBX will use your default AWS profile automatically.")
	fmt.Println()
	fmt.Println("☁️ Google Cloud Storage (GCS):")
	fmt.Println("1️⃣ Install gsutil: https://cloud.google.com/storage/docs/gsutil_install")
	fmt.Println("2️⃣ Run: gcloud auth login")
	fmt.Println("3️⃣ Configure your project: gcloud config set project YOUR_PROJECT_ID")
	fmt.Println()
	fmt.Println("🔷 Azure Blob Storage:")
	fmt.Println("1️⃣ Install Azure CLI: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli")
	fmt.Println("2️⃣ Run: az login")
	fmt.Println("3️⃣ Set your storage account: az storage account show --name YOUR_ACCOUNT")
	fmt.Println()
	fmt.Println("✅ Tip: You can set env vars DBX_S3_BUCKET & DBX_S3_PREFIX for auto-upload.")
	fmt.Print("\nPress ENTER to return...")
	a.reader.ReadString('\n')
	a.MainMenu()
}

func (a *App) BackupMenu() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("--- MySQL Backup Menu ---")
	fmt.Println("[1] Run MySQL Backup")
	fmt.Println("[2] Run MongoDB Backup")
	fmt.Println("[3] Run PostgreSQL Backup")
	fmt.Println("[0] Back to Main Menu")
	fmt.Print("Enter your choice: ")

	choice := a.readInt()
	switch choice {
	case 1:
		a.RunMySQLBackup()
	case 2:
		a.RunMongoBackup()
	case 3:
		a.RunPostgresBackup()
	case 0:
		a.MainMenu()
	default:
		fmt.Println("Invalid choice.")
		a.BackupMenu()
	}
}

//func (a *App) RestoreMenu() {
//	a.clearScreen()
//	a.showBanner()
//	fmt.Println("--- Restore Menu (Coming Soon) ---")
//	fmt.Println("[0] Back to Main Menu")
//	choice := a.readInt()
//	if choice == 0 {
//		a.MainMenu()
//	} else {
//		a.RestoreMenu()
//	}
//}

func (a *App) RestoreMenu() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("--- Restore Menu ---")
	fmt.Println("[1] Restore MySQL Backup")
	fmt.Println("[2] Restore PostgreSQL Backup")
	fmt.Println("[3] Restore MongoDB Backup")
	fmt.Println("[4] Restore SQLite Backup")
	fmt.Println("[0] Back to Main Menu")
	fmt.Print("Enter your choice: ")

	switch a.readInt() {
	case 1:
		a.RunMySQLRestore()
	case 2:
		a.RunPostgresRestore()
	case 3:
		a.RunMongoRestore()
	case 4:
		a.RunSQLiteRestore()
	case 0:
		a.MainMenu()
	default:
		fmt.Println("Invalid choice.")
		a.RestoreMenu()
	}
}

func (a *App) RunMongoBackup() {
	a.clearScreen()
	a.showBanner()
	uri := a.promptInput("MongoDB URI", "mongodb://localhost:27017", false)
	dbname := a.promptInput("Database Name", "", false)
	out := a.promptInput("Backup Directory", "./backups", false)

	err := db.BackupMongo(uri, dbname, out)
	if err != nil {
		fmt.Println("\n❌ Backup failed:", err)
	} else {
		fmt.Println("\n✅ Backup successful!")
	}

	upload := a.promptInput("Upload to AWS S3? (y/N)", "N", false)
	if strings.ToLower(upload) == "y" {
		bucket := a.promptInput("S3 Bucket Name", "my-db-backups", false)
		prefix := a.promptInput("S3 Prefix (folder path)", "dbx/", false)
		// Find the backup file - MongoDB backups create a directory with timestamp, then zip it
		// The zip file will be in the output directory with pattern: dbname_timestamp.zip
		backupPattern := filepath.Join(out, dbname+"_*.zip")
		matches, _ := filepath.Glob(backupPattern)
		var backupFile string
		if len(matches) > 0 {
			backupFile = matches[len(matches)-1] // Use the most recent one
		} else {
			// Fallback: try to find any zip file in the output directory
			allZips, _ := filepath.Glob(filepath.Join(out, "*.zip"))
			if len(allZips) > 0 {
				backupFile = allZips[len(allZips)-1]
			} else {
				fmt.Println("⚠️ No backup zip file found to upload")
				return
			}
		}
		if err := cloud.UploadToS3(backupFile, bucket, prefix); err != nil {
			fmt.Println("❌ Upload failed:", err)
		} else {
			fmt.Println("☁️  Backup uploaded to S3 successfully!")
		}
	}

	fmt.Print("\nPress ENTER to return to Backup Menu...")
	// Ignore ReadString error - always return to menu regardless
	if _, err := a.reader.ReadString('\n'); err == nil {
		a.BackupMenu()
	} else {
		// Even if read fails, return to menu
		a.BackupMenu()
	}
}

//func (a *App) ViewLogs() {
//	a.clearScreen()
//	a.showBanner()
//	fmt.Println("--- Backup Logs (Coming Soon) ---")
//	fmt.Println("[0] Back to Main Menu")
//	choice := a.readInt()
//	if choice == 0 {
//		a.MainMenu()
//	} else {
//		a.ViewLogs()
//	}
//}

func (a *App) ViewLogs() {
	a.clearScreen()
	a.showBanner()
	logFile := "./logs/dbx.log"

	data, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Println("⚠️ No logs found.")
	} else {
		fmt.Println("--- Backup & Restore Logs ---")
		fmt.Println(string(data))
	}

	fmt.Print("\nPress ENTER to return to Main Menu...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.MainMenu()
}

func (a *App) RunMySQLBackup() {
	a.clearScreen()
	a.showBanner()
	host := a.promptInput("MySQL Host", "localhost", false)
	user := a.promptInput("MySQL User", "root", false)
	pass := a.promptInput("MySQL Password", "", true)
	dbname := a.promptInput("Database Name", "", false)
	out := a.promptInput("Backup Directory", "./backups", false)

	err := db.BackupMySQL(host, user, pass, dbname, out)
	if err != nil {
		fmt.Println("\n❌ Backup failed:", err)
	} else {
		fmt.Println("\n✅ Backup successful!")
	}

	upload := a.promptInput("Upload to AWS S3? (y/N)", "N", false)
	if strings.ToLower(upload) == "y" {
		bucket := a.promptInput("S3 Bucket Name", "my-db-backups", false)
		prefix := a.promptInput("S3 Prefix (folder path)", "dbx/", false)
		if err := cloud.UploadToS3(out+".zip", bucket, prefix); err != nil {
			fmt.Println("❌ Upload failed:", err)
		} else {
			fmt.Println("☁️  Backup uploaded to S3 successfully!")
		}
	}

	fmt.Print("\nPress ENTER to return to Backup Menu...")
	// Ignore ReadString error - always return to menu
	_, err = a.reader.ReadString('\n')
	if err != nil {
		// Input error is non-critical, continue to menu
	}
	a.BackupMenu()
}

func (a *App) readInt() int {
	var choice int
	for {
		_, err := fmt.Fscanln(a.reader, &choice)
		if err == nil {
			break
		}
		// If input is not an int, clear the buffer and prompt again
		// Ignore ReadString error - return 0 to indicate invalid input
		_, err = a.reader.ReadString('\n')
		if err != nil {
			// If buffer clear fails, return 0 to indicate invalid input
			return 0
		}
		fmt.Print("Please enter a valid number: ")
	}
	return choice
}

func (a *App) RunPostgresBackup() {
	a.clearScreen()
	a.showBanner()
	host := a.promptInput("PostgreSQL Host", "localhost", false)
	port := a.promptInput("PostgreSQL Port", "5432", false)
	user := a.promptInput("PostgreSQL User", "postgres", false)
	pass := a.promptInput("PostgreSQL Password", "", true)
	dbname := a.promptInput("Database Name", "", false)
	out := a.promptInput("Backup Directory", "./backups", false)

	err := db.BackupPostgres(host, port, user, pass, dbname, out)
	if err != nil {
		fmt.Println("\n❌ Backup failed:", err)
	} else {
		fmt.Println("\n✅ Backup successful!")
	}

	upload := a.promptInput("Upload to AWS S3? (y/N)", "N", false)
	if strings.ToLower(upload) == "y" {
		bucket := a.promptInput("S3 Bucket Name", "my-db-backups", false)
		prefix := a.promptInput("S3 Prefix (folder path)", "dbx/", false)
		if err := cloud.UploadToS3(out+".zip", bucket, prefix); err != nil {
			fmt.Println("❌ Upload failed:", err)
		} else {
			fmt.Println("☁️  Backup uploaded to S3 successfully!")
		}
	}

	fmt.Print("\nPress ENTER to return to Backup Menu...")
	// Ignore ReadString error - always return to menu
	_, err = a.reader.ReadString('\n')
	if err != nil {
		// Input error is non-critical, continue to menu
	}
	a.BackupMenu()
}

func (a *App) RunMySQLRestore() {
	a.clearScreen()
	a.showBanner()
	host := a.promptInput("MySQL Host", "localhost", false)
	user := a.promptInput("MySQL User", "root", false)
	pass := a.promptInput("MySQL Password", "", true)
	dbname := a.promptInput("Database Name", "", false)
	file := a.promptInput("Path to .sql backup file", "./backups/backup.sql", false)

	if err := db.RestoreMySQL(host, user, pass, dbname, file); err != nil {
		fmt.Println("\n❌ Restore failed:", err)
	} else {
		fmt.Println("\n✅ Restore successful!")
	}

	fmt.Print("\nPress ENTER to return to Restore Menu...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.RestoreMenu()
}

func (a *App) RunPostgresRestore() {
	a.clearScreen()
	a.showBanner()
	host := a.promptInput("PostgreSQL Host", "localhost", false)
	port := a.promptInput("PostgreSQL Port", "5432", false)
	user := a.promptInput("PostgreSQL User", "postgres", false)
	pass := a.promptInput("PostgreSQL Password", "", true)
	dbname := a.promptInput("Database Name", "", false)
	file := a.promptInput("Path to backup file", "./backups/backup.dump", false)

	if err := db.RestorePostgres(host, port, user, pass, dbname, file); err != nil {
		fmt.Println("\n❌ Restore failed:", err)
	} else {
		fmt.Println("\n✅ Restore successful!")
	}

	fmt.Print("\nPress ENTER to return to Restore Menu...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.RestoreMenu()
}

func (a *App) RunMongoRestore() {
	a.clearScreen()
	a.showBanner()
	uri := a.promptInput("MongoDB URI", "mongodb://localhost:27017", false)
	dbname := a.promptInput("Database Name", "", false)
	backupDir := a.promptInput("Path to backup folder", "./backups/dbname_timestamp", false)

	if err := db.RestoreMongo(uri, dbname, backupDir); err != nil {
		fmt.Println("\n❌ Restore failed:", err)
	} else {
		fmt.Println("\n✅ Restore successful!")
	}

	fmt.Print("\nPress ENTER to return to Restore Menu...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.RestoreMenu()
}

func (a *App) RunSQLiteRestore() {
	a.clearScreen()
	a.showBanner()
	backupFile := a.promptInput("Path to backup file", "./backups/backup.db.zip", false)
	targetPath := a.promptInput("Target database path (optional)", "", false)

	if err := db.RestoreSQLite(backupFile, targetPath); err != nil {
		fmt.Println("\n❌ Restore failed:", err)
	} else {
		fmt.Println("\n✅ Restore successful!")
	}

	fmt.Print("\nPress ENTER to return to Restore Menu...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.RestoreMenu()
}

func (a *App) TestMySQLConnection() {
	a.clearScreen()
	a.showBanner()
	host := a.promptInput("MySQL Host", "localhost", false)
	user := a.promptInput("MySQL User", "root", false)
	pass := a.promptInput("MySQL Password", "", true)
	dbname := a.promptInput("Database Name", "", false)

	params := map[string]string{"host": host, "user": user, "pass": pass, "dbname": dbname}
	if err := db.TestConnection("mysql", params); err != nil {
		fmt.Println("❌ Connection failed:", err)
	} else {
		fmt.Println("✅ Connection successful!")
	}
	fmt.Print("\nPress ENTER to return...")
	// Ignore ReadString error - always return to menu
	_, err := a.reader.ReadString('\n')
	if err != nil {
		// Input error is non-critical, continue to menu
	}
	a.TestConnectionMenu()
}

func (a *App) TestPostgresConnection() {
	a.clearScreen()
	a.showBanner()
	host := a.promptInput("PostgreSQL Host", "localhost", false)
	port := a.promptInput("PostgreSQL Port", "5432", false)
	user := a.promptInput("PostgreSQL User", "postgres", false)
	pass := a.promptInput("PostgreSQL Password", "", true)
	dbname := a.promptInput("Database Name", "", false)

	params := map[string]string{"host": host, "port": port, "user": user, "pass": pass, "dbname": dbname}
	if err := db.TestConnection("postgres", params); err != nil {
		fmt.Println("❌ Connection failed:", err)
	} else {
		fmt.Println("✅ Connection successful!")
	}
	fmt.Print("\nPress ENTER to return...")
	// Ignore ReadString error - always return to menu
	_, err := a.reader.ReadString('\n')
	if err != nil {
		// Input error is non-critical, continue to menu
	}
	a.TestConnectionMenu()
}

func (a *App) TestMongoConnection() {
	a.clearScreen()
	a.showBanner()
	uri := a.promptInput("MongoDB URI", "mongodb://localhost:27017", false)

	params := map[string]string{"uri": uri}
	if err := db.TestConnection("mongodb", params); err != nil {
		fmt.Println("❌ Connection failed:", err)
	} else {
		fmt.Println("✅ Connection successful!")
	}
	fmt.Print("\nPress ENTER to return...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.TestConnectionMenu()
}

func (a *App) TestSQLiteConnection() {
	a.clearScreen()
	a.showBanner()
	dbPath := a.promptInput("Path to SQLite .db file", "./database.db", false)

	params := map[string]string{"path": dbPath}
	if err := db.TestConnection("sqlite", params); err != nil {
		fmt.Println("❌ Connection failed:", err)
	} else {
		fmt.Println("✅ Connection successful!")
	}
	fmt.Print("\nPress ENTER to return...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.TestConnectionMenu()
}

func (a *App) ScheduleMenu() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("--- Backup Scheduler ---")
	fmt.Println("[1] Add New Scheduled Backup")
	fmt.Println("[2] View Scheduled Jobs")
	fmt.Println("[0] Back to Main Menu")
	fmt.Print("Enter choice: ")

	switch a.readInt() {
	case 1:
		a.AddScheduledBackup()
	case 2:
		a.ViewScheduledJobs()
	case 0:
		a.MainMenu()
	default:
		fmt.Println("Invalid choice.")
		a.ScheduleMenu()
	}
}

func (a *App) AddScheduledBackup() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("Choose DB Type:")
	fmt.Println("[1] MySQL")
	fmt.Println("[2] PostgreSQL")
	fmt.Println("[3] MongoDB")
	fmt.Println("[4] SQLite")
	fmt.Print("Select: ")

	dbChoice := a.readInt()
	var dbType string
	params := make(map[string]string)

	switch dbChoice {
	case 1:
		dbType = "mysql"
		params["host"] = a.promptInput("Host", "localhost", false)
		params["user"] = a.promptInput("User", "root", false)
		params["pass"] = a.promptInput("Password", "", true)
		params["dbname"] = a.promptInput("Database Name", "", false)
		params["out"] = a.promptInput("Backup Dir", "./backups", false)
	case 2:
		dbType = "postgres"
		params["host"] = a.promptInput("Host", "localhost", false)
		params["port"] = a.promptInput("Port", "5432", false)
		params["user"] = a.promptInput("User", "postgres", false)
		params["pass"] = a.promptInput("Password", "", true)
		params["dbname"] = a.promptInput("Database Name", "", false)
		params["out"] = a.promptInput("Backup Dir", "./backups", false)
	case 3:
		dbType = "mongodb"
		params["uri"] = a.promptInput("Mongo URI", "mongodb://localhost:27017", false)
		params["dbname"] = a.promptInput("Database Name", "", false)
		params["out"] = a.promptInput("Backup Dir", "./backups", false)
	case 4:
		dbType = "sqlite"
		params["path"] = a.promptInput("SQLite file path", "./database.db", false)
		params["out"] = a.promptInput("Backup Dir", "./backups", false)
	default:
		fmt.Println("Invalid DB type.")
		a.ScheduleMenu()
		return
	}

	schedule := a.promptInput("Cron schedule (e.g. @daily, @hourly, */30 * * * *)", "@daily", false)

	if err := scheduler.AddJob(dbType, schedule, params); err != nil {
		fmt.Println("❌ Failed to schedule job:", err)
	} else {
		fmt.Println("✅ Backup job scheduled successfully!")
	}

	fmt.Print("\nPress ENTER to return...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.ScheduleMenu()
}

func (a *App) ViewScheduledJobs() {
	a.clearScreen()
	a.showBanner()
	jobs := scheduler.ListJobs()
	if len(jobs) == 0 {
		fmt.Println("⚠️ No scheduled jobs found.")
	} else {
		fmt.Println("--- Scheduled Jobs ---")
		for _, j := range jobs {
			fmt.Printf("[%s] %s @ %s → %v\n", j.CreatedAt.Format("2006-01-02 15:04"), j.DBType, j.Schedule, j.Params)
		}
	}
	fmt.Print("\nPress ENTER to return...")
	// Ignore ReadString error - always return to menu
	a.reader.ReadString('\n')
	a.ScheduleMenu()
}
