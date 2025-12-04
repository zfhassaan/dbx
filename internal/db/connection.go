package db

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// TestConnection checks database connectivity before running backup.
func TestConnection(dbType string, params map[string]string) error {
	switch strings.ToLower(dbType) {
	case "mysql":
		return testMySQL(params)
	case "postgres", "postgresql":
		return testPostgres(params)
	case "mongo", "mongodb":
		return testMongo(params)
	case "sqlite":
		return testSQLite(params)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// ---------------- MySQL ----------------
func testMySQL(params map[string]string) error {
	host := params["host"]
	user := params["user"]
	pass := params["pass"]
	db := params["dbname"]

	if _, err := exec.LookPath("mysql"); err != nil {
		return errors.New("mysql client not found in PATH")
	}

	args := []string{"-h", host, "-u", user, "-e", "SELECT 1", db}
	if pass != "" {
		args = []string{"-h", host, "-u", user, fmt.Sprintf("--password=%s", pass), "-e", "SELECT 1", db}
	}

	cmd := exec.Command("mysql", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Run with timeout to avoid hang if something goes wrong
	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		return errors.New("connection timed out (possibly waiting for password input)")
	}

	return nil
}

// ---------------- PostgreSQL ----------------
// ---------------- PostgreSQL ----------------
func testPostgres(params map[string]string) error {
	host := params["host"]
	port := params["port"]
	user := params["user"]
	pass := params["pass"]
	db := params["dbname"]

	if _, err := exec.LookPath("psql"); err != nil {
		return errors.New("psql not found in PATH")
	}

	// Only set PGPASSWORD if a password is provided
	if pass != "" {
		os.Setenv("PGPASSWORD", pass)
	} else {
		os.Unsetenv("PGPASSWORD")
	}

	args := []string{"-h", host, "-p", port, "-U", user, "-d", db, "-c", "\\q"}

	cmd := exec.Command("psql", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		return errors.New("connection timed out (possibly waiting for password input)")
	}
	return nil
}

// ---------------- MongoDB ----------------
// ---------------- MongoDB ----------------
func testMongo(params map[string]string) error {
	uri := params["uri"]
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	clientCmd := ""
	if _, err := exec.LookPath("mongosh"); err == nil {
		clientCmd = "mongosh"
	} else if _, err := exec.LookPath("mongo"); err == nil {
		clientCmd = "mongo"
	} else {
		return errors.New("neither mongo nor mongosh client found in PATH")
	}

	cmd := exec.Command(clientCmd, uri, "--quiet", "--eval", "db.runCommand({ping:1})")
	cmd.Stdout = nil
	cmd.Stderr = nil

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		return errors.New("connection timed out (possibly waiting for password input)")
	}
	return nil
}

// ---------------- SQLite ----------------
func testSQLite(params map[string]string) error {
	dbPath := params["path"]
	if dbPath == "" {
		return errors.New("sqlite database path required")
	}
	if _, err := os.Stat(dbPath); err != nil {
		return fmt.Errorf("failed to access sqlite file: %w", err)
	}
	return nil
}
