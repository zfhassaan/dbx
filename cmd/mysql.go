package cmd

import (
	"dbx/internal/db"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Variables are declared in backup.go

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Backup a MySQL database",
	RunE: func(cmd *cobra.Command, args []string) error {
		bt := db.BackupTypeFull
		if backupType == "incremental" {
			bt = db.BackupTypeIncremental
		} else if backupType == "differential" {
			bt = db.BackupTypeDifferential
		}

		err := db.BackupMySQLWithType(host, user, password, database, out, bt)
		if err != nil {
			fmt.Println("Backup failed:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Backup complete")
		return nil
	},
}

func init() {
	backupCmd.AddCommand(mysqlCmd)

	mysqlCmd.Flags().StringVar(&host, "host", "localhost", "MySQL host")
	mysqlCmd.Flags().StringVar(&user, "user", "root", "MySQL user")
	mysqlCmd.Flags().StringVar(&password, "password", "", "MySQL password")
	mysqlCmd.Flags().StringVar(&database, "database", "", "MySQL database name")
	mysqlCmd.Flags().StringVar(&out, "out", "./backups", "Output directory")
	mysqlCmd.Flags().StringVar(&backupType, "type", "full", "Backup type: full, incremental, or differential")

	mysqlCmd.MarkFlagRequired("database")
}
