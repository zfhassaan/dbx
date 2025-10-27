package cmd

import (
	"bufio"
	"dbx/internal/db"
	"fmt"
	"os"
	"strings"
)

func RunMongoBackup() {
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

	if err := db.BackupMongo(uri, dbName, out); err != nil {
		fmt.Println("❌ Backup failed:", err)
	} else {
		fmt.Println("✅ Backup successful!")
	}
}
