package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogEntry writes a log line to logs/dbx.log
func LogEntry(dbType, operation, status string, startTime time.Time, err error) {
	endTime := time.Now()
	duration := endTime.Sub(startTime).Round(time.Second)

	// Check for custom log directory from environment variable
	logDir := os.Getenv("DBX_LOG_DIR")
	if logDir == "" {
		logDir = "./logs"
	}
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Println("⚠️ Failed to create logs directory:", err)
		return
	}

	logFile := filepath.Join(logDir, "dbx.log")
	f, ferr := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if ferr != nil {
		fmt.Println("⚠️ Failed to open log file:", ferr)
		return
	}
	// Ignore Close error - log file operations are best-effort
	defer func() { _ = f.Close() }()

	statusMsg := fmt.Sprintf("[%s] %s %s %s (%s)",
		time.Now().Format("2006-01-02 15:04:05"),
		dbType,
		operation,
		status,
		duration,
	)
	if err != nil {
		statusMsg += fmt.Sprintf(" — Error: %v", err)
	}
	statusMsg += "\n"

	if _, werr := f.WriteString(statusMsg); werr != nil {
		fmt.Println("⚠️ Failed to write to log file:", werr)
	}
}
