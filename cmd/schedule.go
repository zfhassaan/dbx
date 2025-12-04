package cmd

import (
	"dbx/internal/scheduler"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	scheduleCron string
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule automated backups",
	Long:  "Schedule automated backups using cron syntax. Example: '0 2 * * *' for daily at 2 AM",
}

var scheduleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new scheduled backup",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)
		
		switch dbType {
		case "mysql":
			params["host"] = host
			params["user"] = user
			params["pass"] = password
			params["dbname"] = database
			params["out"] = out
		case "postgres":
			params["host"] = host
			params["port"] = port
			params["user"] = user
			params["pass"] = password
			params["dbname"] = database
			params["out"] = out
		case "mongodb":
			params["uri"] = uri
			params["dbname"] = database
			params["out"] = out
		case "sqlite":
			params["path"] = sqlitePath
			params["out"] = out
		default:
			return fmt.Errorf("unsupported database type: %s", dbType)
		}

		if err := scheduler.AddJob(dbType, scheduleCron, params); err != nil {
			fmt.Println("Failed to schedule backup:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Backup scheduled successfully")
		return nil
	},
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		jobs := scheduler.ListJobs()
		if len(jobs) == 0 {
			fmt.Println("No scheduled backups found")
			return nil
		}
		
		fmt.Println("Scheduled Backups:")
		for i, job := range jobs {
			dbName := job.Params["dbname"]
			if dbName == "" {
				dbName = job.Params["database"]
			}
			if dbName == "" {
				dbName = "N/A"
			}
			fmt.Printf("%d. %s - %s @ %s\n", i+1, job.DBType, dbName, job.Schedule)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.AddCommand(scheduleAddCmd, scheduleListCmd)

	scheduleAddCmd.Flags().StringVar(&dbType, "db", "", "Database type (mysql, postgres, mongodb, sqlite)")
	scheduleAddCmd.Flags().StringVar(&host, "host", "localhost", "Database host")
	scheduleAddCmd.Flags().StringVar(&port, "port", "5432", "Database port (PostgreSQL)")
	scheduleAddCmd.Flags().StringVar(&user, "user", "", "Database user")
	scheduleAddCmd.Flags().StringVar(&password, "password", "", "Database password")
	scheduleAddCmd.Flags().StringVar(&database, "database", "", "Database name")
	scheduleAddCmd.Flags().StringVar(&uri, "uri", "mongodb://localhost:27017", "MongoDB URI")
	scheduleAddCmd.Flags().StringVar(&sqlitePath, "path", "", "SQLite database path")
	scheduleAddCmd.Flags().StringVar(&out, "out", "./backups", "Output directory")
	scheduleAddCmd.Flags().StringVar(&scheduleCron, "cron", "", "Cron schedule (e.g., '0 2 * * *' for daily at 2 AM)")

	scheduleAddCmd.MarkFlagRequired("db")
	scheduleAddCmd.MarkFlagRequired("cron")
}

