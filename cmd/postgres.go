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

		// Handle cloud upload if requested
		if uploadCloud {
			if err := handleCloudUpload(database, out, "postgres"); err != nil {
				fmt.Printf("⚠️  Cloud upload failed: %v\n", err)
			}
		}

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
	
	// Cloud upload flags
	postgresCmd.Flags().BoolVar(&uploadCloud, "upload", false, "Upload backup to cloud storage")
	postgresCmd.Flags().StringVar(&cloudProvider, "cloud", "s3", "Cloud provider: s3, gcs, or azure")
	postgresCmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket name (or set DBX_S3_BUCKET env var)")
	postgresCmd.Flags().StringVar(&s3Prefix, "s3-prefix", "dbx/", "S3 prefix/folder path")
	postgresCmd.Flags().StringVar(&gcsBucket, "gcs-bucket", "", "GCS bucket name")
	postgresCmd.Flags().StringVar(&gcsPrefix, "gcs-prefix", "dbx/", "GCS prefix/folder path")
	postgresCmd.Flags().StringVar(&azureAccount, "azure-account", "", "Azure storage account name")
	postgresCmd.Flags().StringVar(&azureContainer, "azure-container", "", "Azure container name")
	postgresCmd.Flags().StringVar(&azureBlob, "azure-blob", "", "Azure blob name (optional)")

	postgresCmd.MarkFlagRequired("database")
}

