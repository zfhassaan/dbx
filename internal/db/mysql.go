package db

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func BackupMySQL(host, user, password, database, outDir string) error {
	ts := time.Now().Format("2006-01-02_15-04")
	outFile := filepath.Join(outDir, fmt.Sprintf("%s-%s.sql", database, ts))

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	args := []string{"-h", host, "-u", user}
	if password != "" {
		args = append(args, "-p"+password)
	}
	args = append(args, database)

	cmd := exec.Command("mysqldump", args...)
	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	var stderrBuf bytes.Buffer
	go func() {
		io.Copy(&stderrBuf, stderr)
	}()

	// Start mysqldump process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("mysqldump failed to start: %v", err)
	}

	// Buffer the dump output in memory
	var outputBuf bytes.Buffer
	if _, err := io.Copy(&outputBuf, stdout); err != nil {
		return err
	}

	// Wait for mysqldump to finish and check for error
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("mysqldump failed: %v\n%s", err, stderrBuf.String())
	}

	// Write to .sql file only after successful dump
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := outputBuf.WriteTo(file); err != nil {
		return err
	}

	return nil
}
