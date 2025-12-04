package cmd

import (
	"dbx/internal/db"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Variables are declared in backup.go

var sqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "Backup a SQLite database",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := db.BackupSQLite(sqlitePath, out)
		if err != nil {
			fmt.Println("Backup failed:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Backup complete")
		return nil
	},
}

func init() {
	backupCmd.AddCommand(sqliteCmd)

	sqliteCmd.Flags().StringVar(&sqlitePath, "path", "", "Path to SQLite database file")
	sqliteCmd.Flags().StringVar(&out, "out", "./backups", "Output directory")

	sqliteCmd.MarkFlagRequired("path")
}

