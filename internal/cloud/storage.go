package cloud

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// UploadToS3 uploads the given file/folder to AWS S3 using the AWS CLI.
func UploadToS3(localPath, bucket, prefix string) error {
	if _, err := exec.LookPath("aws"); err != nil {
		return errors.New("aws CLI not found in PATH — install it first: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html")
	}

	if bucket == "" {
		return errors.New("S3 bucket name required")
	}

	key := filepath.Join(prefix, filepath.Base(localPath))
	fmt.Println("☁️  Uploading to S3:", bucket+"/"+key)

	cmd := exec.Command("aws", "s3", "cp", localPath, fmt.Sprintf("s3://%s/%s", bucket, key))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("S3 upload failed: %w", err)
	}

	fmt.Println("✅ Uploaded successfully to s3://" + bucket + "/" + key)
	return nil
}
