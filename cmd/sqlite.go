package cmd

import (
	"dbx/internal/cloud"
	"dbx/internal/db"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		// Handle cloud upload if requested
		if uploadCloud {
			dbName := filepath.Base(sqlitePath)
			if err := handleCloudUpload(dbName, out, "sqlite"); err != nil {
				fmt.Printf("⚠️  Cloud upload failed: %v\n", err)
			}
		}

		return nil
	},
}

func init() {
	backupCmd.AddCommand(sqliteCmd)

	sqliteCmd.Flags().StringVar(&sqlitePath, "path", "", "Path to SQLite database file")
	sqliteCmd.Flags().StringVar(&out, "out", "./backups", "Output directory")
	
	// Cloud upload flags
	sqliteCmd.Flags().BoolVar(&uploadCloud, "upload", false, "Upload backup to cloud storage")
	sqliteCmd.Flags().StringVar(&cloudProvider, "cloud", "s3", "Cloud provider: s3, gcs, or azure")
	sqliteCmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket name (or set DBX_S3_BUCKET env var)")
	sqliteCmd.Flags().StringVar(&s3Prefix, "s3-prefix", "dbx/", "S3 prefix/folder path")
	sqliteCmd.Flags().StringVar(&gcsBucket, "gcs-bucket", "", "GCS bucket name")
	sqliteCmd.Flags().StringVar(&gcsPrefix, "gcs-prefix", "dbx/", "GCS prefix/folder path")
	sqliteCmd.Flags().StringVar(&azureAccount, "azure-account", "", "Azure storage account name")
	sqliteCmd.Flags().StringVar(&azureContainer, "azure-container", "", "Azure container name")
	sqliteCmd.Flags().StringVar(&azureBlob, "azure-blob", "", "Azure blob name (optional)")

	sqliteCmd.MarkFlagRequired("path")
}

