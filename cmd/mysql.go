package cmd

import (
	"dbx/internal/db"
	"github.com/spf13/cobra"
)

var (
	outDir string
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Backup a MySQL database",
	RunE: func(cmd *cobra.Command, args []string) error {
		return db.BackupMySQL(host, user, password, database, outDir)
	},
}

func init() {
	mysqlCmd.Flags().StringVar(&host, "host", "localhost", "MySQL host")
	mysqlCmd.Flags().StringVar(&user, "user", "root", "MySQL user")
	mysqlCmd.Flags().StringVar(&password, "password", "", "MySQL password")
	mysqlCmd.Flags().StringVar(&database, "database", "", "MySQL database name")
	mysqlCmd.Flags().StringVar(&outDir, "out", "./backups", "Output directory")

	mysqlCmd.MarkFlagRequired("database")

	rootCmd.AddCommand(mysqlCmd)
}
