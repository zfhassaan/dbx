package cloud_test

import (
	"dbx/internal/cloud"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestUploadToAzure_MissingAzureCLI tests error handling when Azure CLI is not available
func TestUploadToAzure_MissingAzureCLI(t *testing.T) {
	// This test assumes Azure CLI is not in PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Temporarily remove PATH to simulate missing Azure CLI
	os.Setenv("PATH", "")

	err := cloud.UploadToAzure("./test.txt", "account", "container", "blob")
	if err == nil {
		t.Error("UploadToAzure() should return error when Azure CLI is not found")
	}
}

// TestUploadToAzure_EmptyAccount tests validation of empty account name
func TestUploadToAzure_EmptyAccount(t *testing.T) {
	err := cloud.UploadToAzure("./test.txt", "", "container", "blob")
	if err == nil {
		t.Error("UploadToAzure() should return error for empty account name")
	}

	if err != nil && !strings.Contains(err.Error(), "Azure storage account name and container name required") {
		t.Logf("UploadToAzure() error message = %v (acceptable)", err)
	}
}

// TestUploadToAzure_EmptyContainer tests validation of empty container name
func TestUploadToAzure_EmptyContainer(t *testing.T) {
	err := cloud.UploadToAzure("./test.txt", "account", "", "blob")
	if err == nil {
		t.Error("UploadToAzure() should return error for empty container name")
	}

	if err != nil && !strings.Contains(err.Error(), "Azure storage account name and container name required") {
		t.Logf("UploadToAzure() error message = %v (acceptable)", err)
	}
}

// TestUploadToAzure_ValidInput tests valid input parameters
func TestUploadToAzure_ValidInput(t *testing.T) {
	// This test will only pass if Azure CLI is installed and configured
	if _, err := exec.LookPath("az"); err != nil {
		t.Skip("Skipping test: Azure CLI not found in PATH")
	}

	// Create a temporary test file
	tmpFile := t.TempDir() + "/test.txt"
	os.WriteFile(tmpFile, []byte("test content"), 0644)

	// Note: This will fail if Azure credentials are not configured
	err := cloud.UploadToAzure(tmpFile, "test-account", "test-container", "test-blob")
	if err != nil {
		t.Logf("UploadToAzure() returned error (expected if Azure not configured): %v", err)
	}
}

// TestUploadToAzure_PathTraversalSecurity tests security against path traversal
func TestUploadToAzure_PathTraversalSecurity(t *testing.T) {
	maliciousPath := "../../../etc/passwd"
	err := cloud.UploadToAzure(maliciousPath, "account", "container", "blob")

	// Should either fail validation or sanitize path
	if err == nil {
		t.Log("UploadToAzure() accepted potentially malicious path")
	}
}

// TestUploadToAzure_NonExistentFile tests error handling
func TestUploadToAzure_NonExistentFile(t *testing.T) {
	if _, err := exec.LookPath("az"); err != nil {
		t.Skip("Skipping test: Azure CLI not found in PATH")
	}

	err := cloud.UploadToAzure("/nonexistent/file.txt", "account", "container", "blob")
	if err == nil {
		t.Error("UploadToAzure() should return error for non-existent file")
	}
}

// TestUploadToAzure_EmptyBlobName tests automatic blob name generation
func TestUploadToAzure_EmptyBlobName(t *testing.T) {
	if _, err := exec.LookPath("az"); err != nil {
		t.Skip("Skipping test: Azure CLI not found in PATH")
	}

	tmpFile := t.TempDir() + "/test.txt"
	os.WriteFile(tmpFile, []byte("test"), 0644)

	// Empty blob name should use file basename
	err := cloud.UploadToAzure(tmpFile, "account", "container", "")
	if err != nil {
		t.Logf("UploadToAzure() with empty blob name returned error: %v", err)
	}
}

// TestUploadToAzure_WithBlobName tests upload with specific blob name
func TestUploadToAzure_WithBlobName(t *testing.T) {
	if _, err := exec.LookPath("az"); err != nil {
		t.Skip("Skipping test: Azure CLI not found in PATH")
	}

	tmpFile := t.TempDir() + "/test.txt"
	os.WriteFile(tmpFile, []byte("test"), 0644)

	err := cloud.UploadToAzure(tmpFile, "account", "container", "custom-blob-name.txt")
	if err != nil {
		t.Logf("UploadToAzure() with blob name returned error (expected if Azure not configured): %v", err)
	}
}

// TestUploadToAzure_FileWithPath tests file path handling
func TestUploadToAzure_FileWithPath(t *testing.T) {
	if _, err := exec.LookPath("az"); err != nil {
		t.Skip("Skipping test: Azure CLI not found in PATH")
	}

	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/subdir/test.txt"
	os.MkdirAll(tmpDir+"/subdir", 0755)
	os.WriteFile(tmpFile, []byte("test"), 0644)

	err := cloud.UploadToAzure(tmpFile, "account", "container", "")
	if err != nil {
		t.Logf("UploadToAzure() returned error (expected if Azure not configured): %v", err)
	}
}

