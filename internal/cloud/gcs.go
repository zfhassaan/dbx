package cloud

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// UploadToGCS uploads the given file/folder to Google Cloud Storage using gsutil.
func UploadToGCS(localPath, bucket, prefix string) error {
	if _, err := exec.LookPath("gsutil"); err != nil {
		return errors.New("gsutil not found in PATH — install it first: https://cloud.google.com/storage/docs/gsutil_install")
	}

	if bucket == "" {
		return errors.New("GCS bucket name required")
	}

	key := filepath.Join(prefix, filepath.Base(localPath))
	fmt.Println("☁️  Uploading to GCS:", bucket+"/"+key)

	cmd := exec.Command("gsutil", "cp", localPath, fmt.Sprintf("gs://%s/%s", bucket, key))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("GCS upload failed: %w", err)
	}

	fmt.Println("✅ Uploaded successfully to gs://" + bucket + "/" + key)
	return nil
}

