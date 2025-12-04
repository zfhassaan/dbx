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

		// Handle cloud upload if requested
		if uploadCloud {
			if err := handleCloudUpload(database, out, "mongodb"); err != nil {
				fmt.Printf("⚠️  Cloud upload failed: %v\n", err)
			}
		}

		return nil
	},
}

func init() {
	backupCmd.AddCommand(mongodbCmd)

	mongodbCmd.Flags().StringVar(&uri, "uri", "mongodb://localhost:27017", "MongoDB URI")
	mongodbCmd.Flags().StringVar(&database, "database", "", "Database name")
	mongodbCmd.Flags().StringVar(&out, "out", "./backups", "Output directory")
	
	// Cloud upload flags
	mongodbCmd.Flags().BoolVar(&uploadCloud, "upload", false, "Upload backup to cloud storage")
	mongodbCmd.Flags().StringVar(&cloudProvider, "cloud", "s3", "Cloud provider: s3, gcs, or azure")
	mongodbCmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket name (or set DBX_S3_BUCKET env var)")
	mongodbCmd.Flags().StringVar(&s3Prefix, "s3-prefix", "dbx/", "S3 prefix/folder path")
	mongodbCmd.Flags().StringVar(&gcsBucket, "gcs-bucket", "", "GCS bucket name")
	mongodbCmd.Flags().StringVar(&gcsPrefix, "gcs-prefix", "dbx/", "GCS prefix/folder path")
	mongodbCmd.Flags().StringVar(&azureAccount, "azure-account", "", "Azure storage account name")
	mongodbCmd.Flags().StringVar(&azureContainer, "azure-container", "", "Azure container name")
	mongodbCmd.Flags().StringVar(&azureBlob, "azure-blob", "", "Azure blob name (optional)")

	mongodbCmd.MarkFlagRequired("database")
}
