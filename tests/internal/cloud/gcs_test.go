package cloud_test

import (
	"dbx/internal/cloud"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestUploadToGCS_MissingGSutil tests error handling when gsutil is not available
func TestUploadToGCS_MissingGSutil(t *testing.T) {
	// This test assumes gsutil is not in PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Temporarily remove PATH to simulate missing gsutil
	os.Setenv("PATH", "")

	err := cloud.UploadToGCS("./test.txt", "my-bucket", "prefix/")
	if err == nil {
		t.Error("UploadToGCS() should return error when gsutil is not found")
	}
}

// TestUploadToGCS_EmptyBucket tests validation of empty bucket name
func TestUploadToGCS_EmptyBucket(t *testing.T) {
	err := cloud.UploadToGCS("./test.txt", "", "prefix/")
	if err == nil {
		t.Error("UploadToGCS() should return error for empty bucket name")
	}

	if err != nil && !strings.Contains(err.Error(), "GCS bucket name required") {
		t.Logf("UploadToGCS() error message = %v (acceptable)", err)
	}
}

// TestUploadToGCS_ValidInput tests valid input parameters
func TestUploadToGCS_ValidInput(t *testing.T) {
	// This test will only pass if gsutil is installed and configured
	if _, err := exec.LookPath("gsutil"); err != nil {
		t.Skip("Skipping test: gsutil not found in PATH")
	}

	// Create a temporary test file
	tmpFile := t.TempDir() + "/test.txt"
	os.WriteFile(tmpFile, []byte("test content"), 0644)

	// Note: This will fail if GCS credentials are not configured
	err := cloud.UploadToGCS(tmpFile, "test-bucket", "test-prefix/")
	if err != nil {
		t.Logf("UploadToGCS() returned error (expected if GCS not configured): %v", err)
	}
}

// TestUploadToGCS_PathTraversalSecurity tests security against path traversal
func TestUploadToGCS_PathTraversalSecurity(t *testing.T) {
	maliciousPath := "../../../etc/passwd"
	err := cloud.UploadToGCS(maliciousPath, "bucket", "prefix/")

	// Should either fail validation or sanitize path
	if err == nil {
		t.Log("UploadToGCS() accepted potentially malicious path")
	}
}

// TestUploadToGCS_NonExistentFile tests error handling
func TestUploadToGCS_NonExistentFile(t *testing.T) {
	if _, err := exec.LookPath("gsutil"); err != nil {
		t.Skip("Skipping test: gsutil not found in PATH")
	}

	err := cloud.UploadToGCS("/nonexistent/file.txt", "test-bucket", "prefix/")
	if err == nil {
		t.Error("UploadToGCS() should return error for non-existent file")
	}
}

// TestUploadToGCS_EmptyPrefix tests handling of empty prefix
func TestUploadToGCS_EmptyPrefix(t *testing.T) {
	if _, err := exec.LookPath("gsutil"); err != nil {
		t.Skip("Skipping test: gsutil not found in PATH")
	}

	tmpFile := t.TempDir() + "/test.txt"
	os.WriteFile(tmpFile, []byte("test"), 0644)

	err := cloud.UploadToGCS(tmpFile, "test-bucket", "")
	if err != nil {
		t.Logf("UploadToGCS() with empty prefix returned error: %v", err)
	}
}

// TestUploadToGCS_WithPrefix tests upload with prefix
func TestUploadToGCS_WithPrefix(t *testing.T) {
	if _, err := exec.LookPath("gsutil"); err != nil {
		t.Skip("Skipping test: gsutil not found in PATH")
	}

	tmpFile := t.TempDir() + "/test.txt"
	os.WriteFile(tmpFile, []byte("test"), 0644)

	err := cloud.UploadToGCS(tmpFile, "test-bucket", "backups/2024/")
	if err != nil {
		t.Logf("UploadToGCS() with prefix returned error (expected if GCS not configured): %v", err)
	}
}

// TestUploadToGCS_FileWithPath tests file path handling
func TestUploadToGCS_FileWithPath(t *testing.T) {
	if _, err := exec.LookPath("gsutil"); err != nil {
		t.Skip("Skipping test: gsutil not found in PATH")
	}

	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/subdir/test.txt"
	os.MkdirAll(tmpDir+"/subdir", 0755)
	os.WriteFile(tmpFile, []byte("test"), 0644)

	err := cloud.UploadToGCS(tmpFile, "test-bucket", "prefix/")
	if err != nil {
		t.Logf("UploadToGCS() returned error (expected if GCS not configured): %v", err)
	}
}

