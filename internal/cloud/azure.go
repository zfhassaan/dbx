package cloud

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// UploadToAzure uploads the given file/folder to Azure Blob Storage using az CLI.
func UploadToAzure(localPath, accountName, containerName, blobName string) error {
	if _, err := exec.LookPath("az"); err != nil {
		return errors.New("Azure CLI not found in PATH — install it first: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli")
	}

	if accountName == "" || containerName == "" {
		return errors.New("Azure storage account name and container name required")
	}

	if blobName == "" {
		blobName = filepath.Base(localPath)
	}

	fmt.Println("☁️  Uploading to Azure Blob Storage:", containerName+"/"+blobName)

	cmd := exec.Command("az", "storage", "blob", "upload",
		"--account-name", accountName,
		"--container-name", containerName,
		"--name", blobName,
		"--file", localPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Azure upload failed: %w", err)
	}

	fmt.Println("✅ Uploaded successfully to Azure Blob Storage:", containerName+"/"+blobName)
	return nil
}

