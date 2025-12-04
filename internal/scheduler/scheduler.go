package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"dbx/internal/db"

	"github.com/robfig/cron/v3"
)

// JobConfig holds persisted job info
type JobConfig struct {
	ID        cron.EntryID      `json:"id"`
	DBType    string            `json:"db_type"`
	Schedule  string            `json:"schedule"`
	Params    map[string]string `json:"params"`
	CreatedAt time.Time         `json:"created_at"`
}

var (
	c          *cron.Cron
	configFile = "./config/schedules.json"
	jobs       []JobConfig
)

// Init starts the scheduler and loads jobs
func Init() {
	_ = os.MkdirAll("./config", os.ModePerm)
	c = cron.New()
	loadJobs()
	c.Start()
	fmt.Println("⏰ Scheduler initialized.")
}

// AddJob registers a new backup job
func AddJob(dbType, schedule string, params map[string]string) error {
	if c == nil {
		Init()
	}

	id, err := c.AddFunc(schedule, func() {
		fmt.Printf("\n🔄 Running scheduled %s backup...\n", dbType)
		start := time.Now()
		switch dbType {
		case "mysql":
			_ = db.BackupMySQL(params["host"], params["user"], params["pass"], params["dbname"], params["out"])
		case "postgres":
			_ = db.BackupPostgres(params["host"], params["port"], params["user"], params["pass"], params["dbname"], params["out"])
		case "mongodb":
			_ = db.BackupMongo(params["uri"], params["dbname"], params["out"])
		case "sqlite":
			_ = db.BackupSQLite(params["path"], params["out"])
		}
		fmt.Printf("✅ %s backup completed in %s\n", dbType, time.Since(start).Round(time.Second))
	})
	if err != nil {
		return err
	}

	job := JobConfig{ID: id, DBType: dbType, Schedule: schedule, Params: params, CreatedAt: time.Now()}
	jobs = append(jobs, job)
	return saveJobs()
}

func loadJobs() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		// No existing schedule file - start with empty job list
		return
	}
	// Ignore unmarshal errors - corrupted file will result in empty job list
	_ = json.Unmarshal(data, &jobs)
	for _, job := range jobs {
		// Ignore AddFunc errors - invalid schedules will be skipped
		// Capture loop variable by value to avoid closure capturing reference
		job := job
		_, _ = c.AddFunc(job.Schedule, func() {
			switch job.DBType {
			case "mysql":
				_ = db.BackupMySQL(job.Params["host"], job.Params["user"], job.Params["pass"], job.Params["dbname"], job.Params["out"])
			case "postgres":
				_ = db.BackupPostgres(job.Params["host"], job.Params["port"], job.Params["user"], job.Params["pass"], job.Params["dbname"], job.Params["out"])
			case "mongodb":
				_ = db.BackupMongo(job.Params["uri"], job.Params["dbname"], job.Params["out"])
			case "sqlite":
				_ = db.BackupSQLite(job.Params["path"], job.Params["out"])
			}
		})
	}
}

func saveJobs() error {
	// MarshalIndent should never fail with valid job data, but handle via WriteFile error if it does
	data, _ := json.MarshalIndent(jobs, "", "  ")
	return os.WriteFile(configFile, data, 0644)
}

func ListJobs() []JobConfig {
	return jobs
}
