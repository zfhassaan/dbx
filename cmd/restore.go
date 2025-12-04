package cmd

import (
	"dbx/internal/db"

	"github.com/spf13/cobra"
)

var (
	restoreFile        string
	restoreTable       string
	restoreCollection  string
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a database from backup",
	Long:  "Restore a database from a backup file. Supports full database restore or selective table/collection restore.",
}

var restoreMySQLCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Restore a MySQL database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if restoreTable != "" {
			return db.RestoreMySQLTable(host, user, password, database, restoreFile, restoreTable)
		}
		return db.RestoreMySQL(host, user, password, database, restoreFile)
	},
}

var restorePostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Restore a PostgreSQL database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if restoreTable != "" {
			return db.RestorePostgresTable(host, port, user, password, database, restoreFile, restoreTable)
		}
		return db.RestorePostgres(host, port, user, password, database, restoreFile)
	},
}

var restoreMongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "Restore a MongoDB database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if restoreCollection != "" {
			return db.RestoreMongoCollection(uri, database, restoreFile, restoreCollection)
		}
		return db.RestoreMongo(uri, database, restoreFile)
	},
}

var restoreSQLiteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "Restore a SQLite database",
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := cmd.Flag("target").Value.String()
		return db.RestoreSQLite(restoreFile, targetPath)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.AddCommand(restoreMySQLCmd, restorePostgresCmd, restoreMongoCmd, restoreSQLiteCmd)

	// MySQL restore flags
	restoreMySQLCmd.Flags().StringVar(&host, "host", "localhost", "MySQL host")
	restoreMySQLCmd.Flags().StringVar(&user, "user", "root", "MySQL user")
	restoreMySQLCmd.Flags().StringVar(&password, "password", "", "MySQL password")
	restoreMySQLCmd.Flags().StringVar(&database, "database", "", "Database name")
	restoreMySQLCmd.Flags().StringVar(&restoreFile, "file", "", "Path to backup file")
	restoreMySQLCmd.Flags().StringVar(&restoreTable, "table", "", "Restore specific table only (optional)")
	restoreMySQLCmd.MarkFlagRequired("database")
	restoreMySQLCmd.MarkFlagRequired("file")

	// PostgreSQL restore flags
	restorePostgresCmd.Flags().StringVar(&host, "host", "localhost", "PostgreSQL host")
	restorePostgresCmd.Flags().StringVar(&port, "port", "5432", "PostgreSQL port")
	restorePostgresCmd.Flags().StringVar(&user, "user", "postgres", "PostgreSQL user")
	restorePostgresCmd.Flags().StringVar(&password, "password", "", "PostgreSQL password")
	restorePostgresCmd.Flags().StringVar(&database, "database", "", "Database name")
	restorePostgresCmd.Flags().StringVar(&restoreFile, "file", "", "Path to backup file")
	restorePostgresCmd.Flags().StringVar(&restoreTable, "table", "", "Restore specific table only (optional)")
	restorePostgresCmd.MarkFlagRequired("database")
	restorePostgresCmd.MarkFlagRequired("file")

	// MongoDB restore flags
	restoreMongoCmd.Flags().StringVar(&uri, "uri", "mongodb://localhost:27017", "MongoDB URI")
	restoreMongoCmd.Flags().StringVar(&database, "database", "", "Database name")
	restoreMongoCmd.Flags().StringVar(&restoreFile, "file", "", "Path to backup directory")
	restoreMongoCmd.Flags().StringVar(&restoreCollection, "collection", "", "Restore specific collection only (optional)")
	restoreMongoCmd.MarkFlagRequired("database")
	restoreMongoCmd.MarkFlagRequired("file")

	// SQLite restore flags
	restoreSQLiteCmd.Flags().StringVar(&restoreFile, "file", "", "Path to backup file")
	restoreSQLiteCmd.Flags().String("target", "", "Target database path (optional, defaults to restored_<backup_name>)")
	restoreSQLiteCmd.MarkFlagRequired("file")
}

