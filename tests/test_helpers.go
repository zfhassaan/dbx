package tests

import (
	"os"
	"path/filepath"
	"testing"
)

// TestHelper provides common testing utilities following DRY principle
type TestHelper struct {
	TempDir string
}

// NewTestHelper creates a new test helper instance
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		TempDir: t.TempDir(),
	}
}

// CreateTestFile creates a test file with given content
func (h *TestHelper) CreateTestFile(name, content string) (string, error) {
	filePath := filepath.Join(h.TempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	return filePath, err
}

// CreateTestDir creates a test directory
func (h *TestHelper) CreateTestDir(name string) (string, error) {
	dirPath := filepath.Join(h.TempDir, name)
	err := os.MkdirAll(dirPath, 0755)
	return dirPath, err
}

// FileExists checks if a file exists
func (h *TestHelper) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Cleanup removes all test files (called automatically by testing.T)
func (h *TestHelper) Cleanup() {
	os.RemoveAll(h.TempDir)
}
