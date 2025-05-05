package cmd

import (
	"dbx/internal/db"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	dbType, host, user, password, database, out string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a database",
	Run: func(cmd *cobra.Command, args []string) {
		if dbType != "mysql" {
			fmt.Println("Only MySQL is supported in this MVP")
			return
		}

		err := db.BackupMySQL(host, user, password, database, out)
		if err != nil {
			fmt.Println("Backup failed:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Backup complete")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVar(&dbType, "db", "", "Database type (mysql)")
	backupCmd.Flags().StringVar(&host, "host", "localhost", "Database host")
	backupCmd.Flags().StringVar(&user, "user", "", "Database username")
	backupCmd.Flags().StringVar(&password, "password", "", "Database password")
	backupCmd.Flags().StringVar(&database, "database", "", "Database name")
	backupCmd.Flags().StringVar(&out, "out", "./backups", "Output directory")

	backupCmd.MarkFlagRequired("db")
	backupCmd.MarkFlagRequired("user")
	backupCmd.MarkFlagRequired("password")
	backupCmd.MarkFlagRequired("database")
}
