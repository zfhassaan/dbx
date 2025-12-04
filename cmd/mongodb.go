package cmd

import (
	"dbx/internal/db"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var mongodbCmd = &cobra.Command{
	Use:   "mongo",
	Short: "Backup a MongoDB database",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := db.BackupMongo(uri, database, out)
		if err != nil {
			fmt.Println("Backup failed:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Backup complete")
		return nil
	},
}

func init() {
	backupCmd.AddCommand(mongodbCmd)

	mongodbCmd.Flags().StringVar(&uri, "uri", "mongodb://localhost:27017", "MongoDB URI")
	mongodbCmd.Flags().StringVar(&database, "database", "", "Database name")
	mongodbCmd.Flags().StringVar(&out, "out", "./backups", "Output directory")

	mongodbCmd.MarkFlagRequired("database")
}
