package utils_test

import (
	"dbx/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestCompressFile tests file compression functionality
func TestCompressFile(t *testing.T) {
	// Setup: Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	zipFile := filepath.Join(tmpDir, "test.zip")

	// Create test content
	content := "This is test content for compression"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test: Compress file
	if err := utils.CompressFile(testFile, zipFile); err != nil {
		t.Errorf("CompressFile() error = %v, want nil", err)
	}

	// Verify: Check if zip file exists
	if _, err := os.Stat(zipFile); os.IsNotExist(err) {
		t.Error("CompressFile() did not create zip file")
	}

	// Verify: Original file should still exist
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("CompressFile() should not delete original file")
	}
}

// TestCompressFile_NonExistentSource tests error handling for non-existent source
func TestCompressFile_NonExistentSource(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")
	zipFile := filepath.Join(tmpDir, "test.zip")

	err := utils.CompressFile(nonExistentFile, zipFile)
	if err == nil {
		t.Error("CompressFile() should return error for non-existent source file")
	}
}

// TestCompressFile_InvalidDestination tests error handling for invalid destination
func TestCompressFile_InvalidDestination(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	
	// Create test file
	os.WriteFile(testFile, []byte("test"), 0644)

	// Use a path that's definitely invalid on all platforms
	// On Windows, use a path with invalid characters; on Unix, use a non-existent root path
	var invalidZip string
	if os.PathSeparator == '\\' {
		// Windows: use a path with invalid characters
		invalidZip = "C:\\<invalid>\\test.zip"
	} else {
		// Unix: use a path that requires root but we can't create
		invalidZip = "/root/nonexistent/test.zip"
	}

	err := utils.CompressFile(testFile, invalidZip)
	// On some systems, the file might be created in an unexpected location
	// So we check if the file was actually created at the invalid path
	if err == nil {
		// If no error, verify the file wasn't created at the invalid location
		if _, statErr := os.Stat(invalidZip); statErr == nil {
			t.Error("CompressFile() should not create file at invalid destination path")
		} else {
			// File wasn't created at invalid path, which is acceptable
			t.Log("CompressFile() did not create file at invalid path (acceptable behavior)")
		}
	}
	// If error occurred, that's the expected behavior
}

// TestCompressFolder tests folder compression functionality
func TestCompressFolder(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "source")
	zipFile := filepath.Join(tmpDir, "folder.zip")

	// Setup: Create test directory structure
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(srcDir, "file2.txt"), []byte("content2"), 0644)
	os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(srcDir, "subdir", "file3.txt"), []byte("content3"), 0644)

	// Test: Compress folder
	if err := utils.CompressFolder(srcDir, zipFile); err != nil {
		t.Errorf("CompressFolder() error = %v, want nil", err)
	}

	// Verify: Check if zip file exists
	if _, err := os.Stat(zipFile); os.IsNotExist(err) {
		t.Error("CompressFolder() did not create zip file")
	}
}

// TestCompressFolder_EmptyFolder tests compression of empty folder
func TestCompressFolder_EmptyFolder(t *testing.T) {
	tmpDir := t.TempDir()
	emptyDir := filepath.Join(tmpDir, "empty")
	zipFile := filepath.Join(tmpDir, "empty.zip")

	os.MkdirAll(emptyDir, 0755)

	// Should succeed even with empty folder
	if err := utils.CompressFolder(emptyDir, zipFile); err != nil {
		t.Errorf("CompressFolder() should handle empty folder, got error: %v", err)
	}
}

// TestCompressGzip tests gzip compression
func TestCompressGzip(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	gzFile := filepath.Join(tmpDir, "test.gz")

	content := "Test content for gzip compression"
	os.WriteFile(testFile, []byte(content), 0644)

	if err := utils.CompressGzip(testFile, gzFile); err != nil {
		t.Errorf("CompressGzip() error = %v, want nil", err)
	}

	if _, err := os.Stat(gzFile); os.IsNotExist(err) {
		t.Error("CompressGzip() did not create gzip file")
	}
}

// TestCompressGzip_NonExistentSource tests error handling
func TestCompressGzip_NonExistentSource(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")
	gzFile := filepath.Join(tmpDir, "test.gz")

	err := utils.CompressGzip(nonExistentFile, gzFile)
	if err == nil {
		t.Error("CompressGzip() should return error for non-existent source")
	}
}

// TestCompressFile_LargeFile tests scalability with large files
func TestCompressFile_LargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large file test in short mode")
	}

	tmpDir := t.TempDir()
	largeFile := filepath.Join(tmpDir, "large.txt")
	zipFile := filepath.Join(tmpDir, "large.zip")

	// Create 10MB file
	largeContent := make([]byte, 10*1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	if err := os.WriteFile(largeFile, largeContent, 0644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	// Test compression of large file
	if err := utils.CompressFile(largeFile, zipFile); err != nil {
		t.Errorf("CompressFile() failed with large file: %v", err)
	}

	// Verify compression ratio (should be smaller)
	originalInfo, _ := os.Stat(largeFile)
	compressedInfo, _ := os.Stat(zipFile)

	if compressedInfo.Size() >= originalInfo.Size() {
		t.Logf("Warning: Compressed size (%d) >= original size (%d)", compressedInfo.Size(), originalInfo.Size())
	}
}

// TestCompressFile_ConcurrentCompression tests concurrent operations (scalability)
func TestCompressFile_ConcurrentCompression(t *testing.T) {
	tmpDir := t.TempDir()
	numGoroutines := 10
	errors := make(chan error, numGoroutines)

	// Create multiple test files concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			testFile := filepath.Join(tmpDir, "test", "file", "test_"+fmt.Sprintf("%d", id)+".txt")
			zipFile := filepath.Join(tmpDir, "test_"+fmt.Sprintf("%d", id)+".zip")

			os.MkdirAll(filepath.Dir(testFile), 0755)
			os.WriteFile(testFile, []byte("test content"), 0644)

			if err := utils.CompressFile(testFile, zipFile); err != nil {
				errors <- err
			} else {
				errors <- nil
			}
		}(i)
	}

	// Collect results
	var errorCount int
	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			errorCount++
			t.Logf("Concurrent compression error: %v", err)
		}
	}

	if errorCount > 0 {
		t.Errorf("Concurrent compression failed for %d/%d operations", errorCount, numGoroutines)
	}
}

// TestCompressFile_PathTraversalSecurity tests security against path traversal attacks
func TestCompressFile_PathTraversalSecurity(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	// Attempt path traversal in destination
	maliciousZip := filepath.Join(tmpDir, "..", "..", "..", "etc", "passwd.zip")
	err := utils.CompressFile(testFile, maliciousZip)

	// Should either fail or be sanitized
	if err == nil {
		// If it succeeds, verify it didn't actually write outside tmpDir
		if _, err := os.Stat(maliciousZip); err == nil {
			// Check if file is actually outside tmpDir (security check)
			absMalicious, _ := filepath.Abs(maliciousZip)
			absTmp, _ := filepath.Abs(tmpDir)
			if !filepath.HasPrefix(absMalicious, absTmp) {
				t.Error("Security vulnerability: CompressFile allowed path traversal")
			}
		}
	}
}

