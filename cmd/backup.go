package cmd

import (
	"dbx/internal/cloud"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Shared variables for all commands
var (
	dbType, host, user, password, database, out, port, uri string
	backupType, sqlitePath string
	outDir string // Alias for out, used in some commands
	// Cloud upload flags
	uploadCloud      bool
	cloudProvider    string // s3, gcs, azure
	s3Bucket, s3Prefix string
	gcsBucket, gcsPrefix string
	azureAccount, azureContainer, azureBlob string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a database",
	Long:  "Backup a database. Use subcommands (mysql, postgres, mongo, sqlite) for specific database types.",
}

func init() {
	rootCmd.AddCommand(backupCmd)
	// Subcommands are added in their respective files
}

// handleCloudUpload handles cloud upload for backup files
func handleCloudUpload(dbName, outDir, dbType string) error {
	// Find the most recent backup file
	pattern := filepath.Join(outDir, dbName+"*")
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return fmt.Errorf("no backup file found in %s", outDir)
	}
	
	// Use the most recent file
	backupFile := matches[len(matches)-1]
	
	switch strings.ToLower(cloudProvider) {
	case "s3":
		bucket := s3Bucket
		if bucket == "" {
			bucket = os.Getenv("DBX_S3_BUCKET")
		}
		if bucket == "" {
			return fmt.Errorf("S3 bucket name required (use --s3-bucket or set DBX_S3_BUCKET env var)")
		}
		prefix := s3Prefix
		if prefix == "" {
			prefix = os.Getenv("DBX_S3_PREFIX")
			if prefix == "" {
				prefix = "dbx/"
			}
		}
		return cloud.UploadToS3(backupFile, bucket, prefix)
	case "gcs":
		if gcsBucket == "" {
			return fmt.Errorf("GCS bucket name required (use --gcs-bucket)")
		}
		return cloud.UploadToGCS(backupFile, gcsBucket, gcsPrefix)
	case "azure":
		if azureAccount == "" || azureContainer == "" {
			return fmt.Errorf("Azure account and container required (use --azure-account and --azure-container)")
		}
		return cloud.UploadToAzure(backupFile, azureAccount, azureContainer, azureBlob)
	default:
		return fmt.Errorf("unsupported cloud provider: %s (use s3, gcs, or azure)", cloudProvider)
	}
}
