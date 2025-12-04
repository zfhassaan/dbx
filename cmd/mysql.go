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

		// Handle cloud upload if requested
		if uploadCloud {
			if err := handleCloudUpload(database, out, "mysql"); err != nil {
				fmt.Printf("⚠️  Cloud upload failed: %v\n", err)
				// Don't fail the backup if upload fails
			}
		}

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
	
	// Cloud upload flags
	mysqlCmd.Flags().BoolVar(&uploadCloud, "upload", false, "Upload backup to cloud storage")
	mysqlCmd.Flags().StringVar(&cloudProvider, "cloud", "s3", "Cloud provider: s3, gcs, or azure")
	mysqlCmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket name (or set DBX_S3_BUCKET env var)")
	mysqlCmd.Flags().StringVar(&s3Prefix, "s3-prefix", "dbx/", "S3 prefix/folder path")
	mysqlCmd.Flags().StringVar(&gcsBucket, "gcs-bucket", "", "GCS bucket name")
	mysqlCmd.Flags().StringVar(&gcsPrefix, "gcs-prefix", "dbx/", "GCS prefix/folder path")
	mysqlCmd.Flags().StringVar(&azureAccount, "azure-account", "", "Azure storage account name")
	mysqlCmd.Flags().StringVar(&azureContainer, "azure-container", "", "Azure container name")
	mysqlCmd.Flags().StringVar(&azureBlob, "azure-blob", "", "Azure blob name (optional)")

	mysqlCmd.MarkFlagRequired("database")
}
