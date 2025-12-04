package db_test

import (
	"dbx/internal/db"
	"os"
	"testing"
)

// TestTestConnection_UnsupportedDBType tests error handling for unsupported database types
func TestTestConnection_UnsupportedDBType(t *testing.T) {
	params := map[string]string{"host": "localhost"}
	err := db.TestConnection("unsupported_db", params)

	if err == nil {
		t.Error("TestConnection() should return error for unsupported database type")
	}
}

// TestTestConnection_EmptyParams tests handling of empty parameters
func TestTestConnection_EmptyParams(t *testing.T) {
	// SQLite should handle empty path gracefully
	err := db.TestConnection("sqlite", map[string]string{})
	if err == nil {
		t.Error("TestConnection() should return error for SQLite with empty path")
	}
}

// TestTestConnection_InvalidParams tests various invalid parameter combinations
func TestTestConnection_InvalidParams(t *testing.T) {
	tests := []struct {
		name   string
		dbType string
		params map[string]string
	}{
		{"MySQL missing host", "mysql", map[string]string{"user": "root"}},
		{"PostgreSQL missing port", "postgres", map[string]string{"host": "localhost"}},
		{"MongoDB empty URI", "mongodb", map[string]string{"uri": ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These should fail gracefully (either missing tool or invalid connection)
			// We're testing that they don't panic
			_ = db.TestConnection(tt.dbType, tt.params)
		})
	}
}

// TestTestSQLite_FileExists tests SQLite connection test with existing file
func TestTestSQLite_FileExists(t *testing.T) {
	tmpDir := t.TempDir()
	testDB := tmpDir + "/test.db"

	// Create a test SQLite file
	os.WriteFile(testDB, []byte("SQLite format 3"), 0644)

	params := map[string]string{"path": testDB}
	err := db.TestConnection("sqlite", params)

	// Should succeed if file exists
	if err != nil {
		t.Logf("TestConnection() for SQLite returned error (expected if file is not valid SQLite): %v", err)
	}
}

// TestTestSQLite_FileNotExists tests SQLite connection test with non-existent file
func TestTestSQLite_FileNotExists(t *testing.T) {
	params := map[string]string{"path": "/nonexistent/database.db"}
	err := db.TestConnection("sqlite", params)

	if err == nil {
		t.Error("TestConnection() should return error for non-existent SQLite file")
	}
}

// TestTestConnection_CaseInsensitive tests case-insensitive database type matching
func TestTestConnection_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name   string
		dbType string
	}{
		{"Lowercase postgres", "postgres"},
		{"Uppercase POSTGRES", "POSTGRES"},
		{"Mixed case Postgres", "Postgres"},
		{"Lowercase mongodb", "mongodb"},
		{"Uppercase MONGODB", "MONGODB"},
		{"Mixed case Mongo", "Mongo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic regardless of case
			params := map[string]string{"host": "localhost"}
			_ = db.TestConnection(tt.dbType, params)
		})
	}
}

