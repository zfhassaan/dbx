package cmd

import (
	"github.com/spf13/cobra"
)

// Shared variables for all commands
var (
	dbType, host, user, password, database, out, port, uri string
	backupType, sqlitePath string
	outDir string // Alias for out, used in some commands
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a database",
	Long:  "Backup a database. Use subcommands (mysql, postgres, mongo, sqlite) for specific database types.",
}

func init() {
	rootCmd.AddCommand(backupCmd)
	// Subcommands are added in their respective files
}
