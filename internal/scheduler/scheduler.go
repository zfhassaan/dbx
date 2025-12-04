package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dbx/internal/cloud"
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
		var backupErr error
		var dbName string
		
		switch dbType {
		case "mysql":
			dbName = params["dbname"]
			backupErr = db.BackupMySQL(params["host"], params["user"], params["pass"], dbName, params["out"])
		case "postgres":
			dbName = params["dbname"]
			backupErr = db.BackupPostgres(params["host"], params["port"], params["user"], params["pass"], dbName, params["out"])
		case "mongodb":
			dbName = params["dbname"]
			backupErr = db.BackupMongo(params["uri"], dbName, params["out"])
		case "sqlite":
			dbName = filepath.Base(params["path"])
			backupErr = db.BackupSQLite(params["path"], params["out"])
		}
		
		if backupErr != nil {
			fmt.Printf("❌ %s backup failed: %v\n", dbType, backupErr)
		} else {
			fmt.Printf("✅ %s backup completed in %s\n", dbType, time.Since(start).Round(time.Second))
			
			// Handle cloud upload if configured
			if params["upload_cloud"] == "true" || os.Getenv("DBX_AUTO_UPLOAD") == "true" {
				if err := handleScheduledCloudUpload(dbName, params); err != nil {
					fmt.Printf("⚠️  Cloud upload failed: %v\n", err)
				} else {
					fmt.Printf("☁️  Backup uploaded to cloud storage\n")
				}
			}
		}
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
			var backupErr error
			var dbName string
			
			switch job.DBType {
			case "mysql":
				dbName = job.Params["dbname"]
				backupErr = db.BackupMySQL(job.Params["host"], job.Params["user"], job.Params["pass"], dbName, job.Params["out"])
			case "postgres":
				dbName = job.Params["dbname"]
				backupErr = db.BackupPostgres(job.Params["host"], job.Params["port"], job.Params["user"], job.Params["pass"], dbName, job.Params["out"])
			case "mongodb":
				dbName = job.Params["dbname"]
				backupErr = db.BackupMongo(job.Params["uri"], dbName, job.Params["out"])
			case "sqlite":
				dbName = filepath.Base(job.Params["path"])
				backupErr = db.BackupSQLite(job.Params["path"], job.Params["out"])
			}
			
			if backupErr == nil {
				// Handle cloud upload if configured
				if job.Params["upload_cloud"] == "true" || os.Getenv("DBX_AUTO_UPLOAD") == "true" {
					_ = handleScheduledCloudUpload(dbName, job.Params)
				}
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

// handleScheduledCloudUpload handles cloud upload for scheduled backups
func handleScheduledCloudUpload(dbName string, params map[string]string) error {
	outDir := params["out"]
	if outDir == "" {
		outDir = "./backups"
	}
	
	// Find the most recent backup file
	pattern := filepath.Join(outDir, dbName+"*")
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return fmt.Errorf("no backup file found in %s", outDir)
	}
	
	backupFile := matches[len(matches)-1]
	cloudProvider := params["cloud_provider"]
	if cloudProvider == "" {
		cloudProvider = os.Getenv("DBX_CLOUD_PROVIDER")
		if cloudProvider == "" {
			cloudProvider = "s3"
		}
	}
	
	switch strings.ToLower(cloudProvider) {
	case "s3":
		bucket := params["s3_bucket"]
		if bucket == "" {
			bucket = os.Getenv("DBX_S3_BUCKET")
		}
		if bucket == "" {
			return fmt.Errorf("S3 bucket name required (set s3_bucket in schedule params or DBX_S3_BUCKET env var)")
		}
		prefix := params["s3_prefix"]
		if prefix == "" {
			prefix = os.Getenv("DBX_S3_PREFIX")
			if prefix == "" {
				prefix = "dbx/"
			}
		}
		return cloud.UploadToS3(backupFile, bucket, prefix)
	case "gcs":
		bucket := params["gcs_bucket"]
		if bucket == "" {
			return fmt.Errorf("GCS bucket name required (set gcs_bucket in schedule params)")
		}
		prefix := params["gcs_prefix"]
		if prefix == "" {
			prefix = "dbx/"
		}
		return cloud.UploadToGCS(backupFile, bucket, prefix)
	case "azure":
		account := params["azure_account"]
		container := params["azure_container"]
		if account == "" || container == "" {
			return fmt.Errorf("Azure account and container required (set azure_account and azure_container in schedule params)")
		}
		return cloud.UploadToAzure(backupFile, account, container, params["azure_blob"])
	default:
		return fmt.Errorf("unsupported cloud provider: %s", cloudProvider)
	}
}
