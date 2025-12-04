package cmd

import (
	"dbx/internal/db"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Backup a PostgreSQL database",
	RunE: func(cmd *cobra.Command, args []string) error {
		bt := db.BackupTypeFull
		if backupType == "incremental" {
			bt = db.BackupTypeIncremental
		} else if backupType == "differential" {
			bt = db.BackupTypeDifferential
		}

		err := db.BackupPostgresWithType(host, port, user, password, database, out, bt)
		if err != nil {
			fmt.Println("Backup failed:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Backup complete")
		return nil
	},
}

func init() {
	backupCmd.AddCommand(postgresCmd)

	postgresCmd.Flags().StringVar(&host, "host", "localhost", "PostgreSQL host")
	postgresCmd.Flags().StringVar(&port, "port", "5432", "PostgreSQL port")
	postgresCmd.Flags().StringVar(&user, "user", "postgres", "PostgreSQL user")
	postgresCmd.Flags().StringVar(&password, "password", "", "PostgreSQL password")
	postgresCmd.Flags().StringVar(&database, "database", "", "Database name")
	postgresCmd.Flags().StringVar(&out, "out", "./backups", "Output directory")
	postgresCmd.Flags().StringVar(&backupType, "type", "full", "Backup type: full, incremental, or differential")

	postgresCmd.MarkFlagRequired("database")
}

