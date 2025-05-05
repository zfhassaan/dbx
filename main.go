package main

import (
	"bufio"
	"dbx/internal/db"
	"fmt"
	"os"
	"strings"
)

func main() {
	showMenu()
}

func showBanner() {
	fmt.Println(`
  ____  ____  __  __       
 |  _ \| __ )|  \/  | ___  
 | | | |  _ \| |\/| |/ _ \ 
 | |_| | |_) | |  | | (_) |
 |____/|____/|_|  |_|\___/  v0.0.1

   DBX — Dead-simple database backups.
 ----------------------------
      Author: @zfhassaan
 ----------------------------`)
}

func showMenu() {
	showBanner()

	fmt.Println("\nChoose an option:")
	fmt.Println("[1] MySQL Backup")
	fmt.Println("[0] Exit")

	fmt.Print("\nEnter choice: ")
	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		runMySQLBackup()
	case 0:
		fmt.Println("Exiting DBX.")
		os.Exit(0)
	default:
		fmt.Println("Invalid option.\n")
		showMenu()
	}
}

func runMySQLBackup() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("MySQL Host [localhost]: ")
	host, _ := reader.ReadString('\n')
	if host == "\n" {
		host = "localhost"
	}

	fmt.Print("MySQL User [root]: ")
	user, _ := reader.ReadString('\n')
	if user == "\n" {
		user = "root"
	}

	fmt.Print("MySQL Password: ")
	pass, _ := reader.ReadString('\n')

	fmt.Print("Database Name: ")
	dbname, _ := reader.ReadString('\n')

	fmt.Print("Backup Directory [./backups]: ")
	out, _ := reader.ReadString('\n')
	if out == "\n" {
		out = "./backups"
	}

	// Trim newline characters
	host = trim(host)
	user = trim(user)
	pass = trim(pass)
	dbname = trim(dbname)
	out = trim(out)

	err := db.BackupMySQL(host, user, pass, dbname, out)
	if err != nil {
		fmt.Println("\n❌ Backup failed:", err)
	} else {
		fmt.Println("\n✅ Backup successful!")
	}

	fmt.Print("\nPress ENTER to return to menu...")
	reader.ReadString('\n')
	showMenu()
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
