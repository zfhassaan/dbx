package main

import (
	"bufio"
	"dbx/internal/db"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
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
			return
		}
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
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
 |____/|____/|_|  |_|\___/  v0.0.1

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
	if hideInput {
		// For simplicity, we use normal input; for real password hiding, use a library.
	}
	input, _ := a.reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" && defaultVal != "" {
		return defaultVal
	}
	return input
}

func (a *App) MainMenu() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("===== DBX: Database Backup Utility =====")
	fmt.Println("[1] 🔄 Backup Menu")
	fmt.Println("[2] 🔁 Restore Menu")
	fmt.Println("[3] 📜 View Logs")
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
	case 0:
		fmt.Println("👋 Exiting DBX.")
		os.Exit(0)
	default:
		fmt.Println("⚠️ Invalid option. Try again.")
		a.MainMenu()
	}
}

func (a *App) BackupMenu() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("--- MySQL Backup Menu ---")
	fmt.Println("[1] Run MySQL Backup")
	fmt.Println("[0] Back to Main Menu")
	fmt.Print("Enter your choice: ")

	choice := a.readInt()
	switch choice {
	case 1:
		a.RunMySQLBackup()
	case 0:
		a.MainMenu()
	default:
		fmt.Println("Invalid choice.")
		a.BackupMenu()
	}
}

func (a *App) RestoreMenu() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("--- Restore Menu (Coming Soon) ---")
	fmt.Println("[0] Back to Main Menu")
	choice := a.readInt()
	if choice == 0 {
		a.MainMenu()
	} else {
		a.RestoreMenu()
	}
}

func (a *App) ViewLogs() {
	a.clearScreen()
	a.showBanner()
	fmt.Println("--- Backup Logs (Coming Soon) ---")
	fmt.Println("[0] Back to Main Menu")
	choice := a.readInt()
	if choice == 0 {
		a.MainMenu()
	} else {
		a.ViewLogs()
	}
}

func (a *App) RunMySQLBackup() {
	a.clearScreen()
	a.showBanner()
	host := a.promptInput("MySQL Host", "localhost", false)
	user := a.promptInput("MySQL User", "root", false)
	pass := a.promptInput("MySQL Password", "", false)
	dbname := a.promptInput("Database Name", "", false)
	out := a.promptInput("Backup Directory", "./backups", false)

	err := db.BackupMySQL(host, user, pass, dbname, out)
	if err != nil {
		fmt.Println("\n❌ Backup failed:", err)
	} else {
		fmt.Println("\n✅ Backup successful!")
	}

	fmt.Print("\nPress ENTER to return to Backup Menu...")
	_, err = a.reader.ReadString('\n')
	if err != nil {
		return
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
		_, err = a.reader.ReadString('\n')
		if err != nil {
			return 0
		}
		fmt.Print("Please enter a valid number: ")
	}
	return choice
}
