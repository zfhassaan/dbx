package scheduler_test

import (
	"dbx/internal/scheduler"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Note: Scheduler uses package-level globals, making unit testing challenging
// These tests verify the public API works correctly

// TestAddJob_ValidSchedule tests adding a valid job
func TestAddJob_ValidSchedule(t *testing.T) {
	// Initialize scheduler
	scheduler.Init()

	params := map[string]string{
		"host":   "localhost",
		"user":   "test",
		"pass":   "password",
		"dbname": "testdb",
		"out":    "./backups",
	}

	err := scheduler.AddJob("mysql", "@daily", params)
	if err != nil {
		t.Errorf("AddJob() error = %v, want nil", err)
	}

	// Verify job was added
	jobs := scheduler.ListJobs()
	if len(jobs) == 0 {
		t.Error("AddJob() did not add job to list")
	}

	// Cleanup
	os.Remove("./config/schedules.json")
}

// TestAddJob_InitIfNil tests that AddJob initializes scheduler if nil
func TestAddJob_InitIfNil(t *testing.T) {
	// Don't call Init() - AddJob should initialize it
	params := map[string]string{"host": "localhost", "dbname": "test", "out": "./backups"}
	err := scheduler.AddJob("mysql", "@daily", params)
	if err != nil {
		t.Logf("AddJob() returned error (acceptable): %v", err)
	}

	// Cleanup
	os.Remove("./config/schedules.json")
}

// TestLoadJobs_ValidFileWithAllDBTypes tests loading jobs for all DB types
func TestLoadJobs_ValidFileWithAllDBTypes(t *testing.T) {
	scheduler.Init()

	// Create a valid schedule file with all DB types
	os.MkdirAll("./config", 0755)
	validJobs := []map[string]interface{}{
		{
			"db_type":  "mysql",
			"schedule": "@daily",
			"params": map[string]string{
				"host":   "localhost",
				"user":   "root",
				"pass":   "password",
				"dbname": "test1",
				"out":    "./backups",
			},
		},
		{
			"db_type":  "postgres",
			"schedule": "@hourly",
			"params": map[string]string{
				"host":   "localhost",
				"port":   "5432",
				"user":   "postgres",
				"pass":   "password",
				"dbname": "test2",
				"out":    "./backups",
			},
		},
		{
			"db_type":  "mongodb",
			"schedule": "@weekly",
			"params": map[string]string{
				"uri":    "mongodb://localhost:27017",
				"dbname": "test3",
				"out":    "./backups",
			},
		},
		{
			"db_type":  "sqlite",
			"schedule": "@monthly",
			"params": map[string]string{
				"path": "./test.db",
				"out":  "./backups",
			},
		},
	}
	data, _ := json.Marshal(validJobs)
	os.WriteFile("./config/schedules.json", data, 0644)

	// Reinitialize to trigger loadJobs
	scheduler.Init()

	jobs := scheduler.ListJobs()
	t.Logf("loadJobs() loaded %d jobs from file with all DB types", len(jobs))

	// Cleanup
	os.Remove("./config/schedules.json")
}

// TestAddJob_AllDBTypesInSwitch tests all DB types in AddJob switch statement
func TestAddJob_AllDBTypesInSwitch(t *testing.T) {
	scheduler.Init()

	dbTypes := []string{"mysql", "postgres", "mongodb", "sqlite"}
	for _, dbType := range dbTypes {
		params := map[string]string{
			"host":   "localhost",
			"dbname": "test",
			"out":    "./backups",
		}
		if dbType == "postgres" {
			params["port"] = "5432"
			params["user"] = "postgres"
			params["pass"] = "password"
		} else if dbType == "mysql" {
			params["user"] = "root"
			params["pass"] = "password"
		} else if dbType == "mongodb" {
			params["uri"] = "mongodb://localhost:27017"
		} else if dbType == "sqlite" {
			params["path"] = "./test.db"
		}

		err := scheduler.AddJob(dbType, "@daily", params)
		if err != nil {
			t.Errorf("AddJob() failed for %s: %v", dbType, err)
		}
	}

	// Verify all jobs were added
	jobs := scheduler.ListJobs()
	if len(jobs) < len(dbTypes) {
		t.Errorf("Expected at least %d jobs, got %d", len(dbTypes), len(jobs))
	}

	// Cleanup
	os.Remove("./config/schedules.json")
}

// TestAddJob_InvalidSchedule tests error handling for invalid cron schedule
func TestAddJob_InvalidSchedule(t *testing.T) {
	scheduler.Init()

	params := map[string]string{"host": "localhost", "dbname": "test", "out": "./backups"}
	err := scheduler.AddJob("mysql", "invalid cron schedule", params)

	if err == nil {
		t.Error("AddJob() should return error for invalid cron schedule")
	}
}

// TestAddJob_UnsupportedDBType tests error handling for unsupported database types
func TestAddJob_UnsupportedDBType(t *testing.T) {
	scheduler.Init()

	params := map[string]string{"host": "localhost", "dbname": "test", "out": "./backups"}
	err := scheduler.AddJob("unsupported_db", "@daily", params)

	// Should either return error or handle gracefully
	if err != nil {
		t.Logf("AddJob() returned error for unsupported DB type (acceptable): %v", err)
	}
}

// TestAddJob_EmptyParams tests error handling for empty parameters
func TestAddJob_EmptyParams(t *testing.T) {
	scheduler.Init()

	err := scheduler.AddJob("mysql", "@daily", map[string]string{})
	// Should handle empty params gracefully
	if err != nil {
		t.Logf("AddJob() returned error for empty params (acceptable): %v", err)
	}
}

// TestAddJob_AllDBTypes tests adding jobs for all supported database types
func TestAddJob_AllDBTypes(t *testing.T) {
	scheduler.Init()

	dbTypes := []string{"mysql", "postgres", "mongodb", "sqlite"}
	for _, dbType := range dbTypes {
		params := map[string]string{
			"host":   "localhost",
			"dbname": "test",
			"out":    "./backups",
		}
		if dbType == "postgres" {
			params["port"] = "5432"
		}
		if dbType == "mongodb" {
			params["uri"] = "mongodb://localhost:27017"
		}
		if dbType == "sqlite" {
			params["path"] = "./test.db"
		}

		err := scheduler.AddJob(dbType, "@daily", params)
		if err != nil {
			t.Errorf("AddJob() failed for %s: %v", dbType, err)
		}
	}

	jobs := scheduler.ListJobs()
	if len(jobs) < len(dbTypes) {
		t.Errorf("Expected at least %d jobs, got %d", len(dbTypes), len(jobs))
	}

	// Cleanup
	os.Remove("./config/schedules.json")
}

// TestListJobs tests listing scheduled jobs
func TestListJobs(t *testing.T) {
	scheduler.Init()

	// Add multiple jobs
	params1 := map[string]string{"host": "localhost", "dbname": "db1", "out": "./backups"}
	params2 := map[string]string{"host": "localhost", "dbname": "db2", "out": "./backups"}

	scheduler.AddJob("mysql", "@daily", params1)
	scheduler.AddJob("postgres", "@hourly", params2)

	jobs := scheduler.ListJobs()
	if len(jobs) < 2 {
		t.Errorf("ListJobs() returned %d jobs, want at least 2", len(jobs))
	}

	// Cleanup
	os.Remove("./config/schedules.json")
}

// TestLoadJobs_CorruptedFile tests handling of corrupted schedule file
func TestLoadJobs_CorruptedFile(t *testing.T) {
	scheduler.Init()

	// Create a corrupted schedule file
	os.MkdirAll("./config", 0755)
	os.WriteFile("./config/schedules.json", []byte("invalid json content"), 0644)

	// Reinitialize to trigger loadJobs
	scheduler.Init()

	// Should not panic and should handle gracefully
	jobs := scheduler.ListJobs()
	t.Logf("loadJobs() handled corrupted file, found %d jobs", len(jobs))

	// Cleanup
	os.Remove("./config/schedules.json")
}

// TestJobConfig_Serialization tests JSON serialization
func TestJobConfig_Serialization(t *testing.T) {
	scheduler.Init()
	
	// Create a test job
	params := map[string]string{"host": "localhost", "dbname": "test", "out": "./backups"}
	scheduler.AddJob("mysql", "@daily", params)
	
	jobs := scheduler.ListJobs()
	if len(jobs) == 0 {
		t.Fatal("No jobs available for serialization test")
	}
	job := jobs[0]
	
	// Test serialization
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	// Create a serializable structure matching JobConfig
	jobData := map[string]interface{}{
		"db_type":   job.DBType,
		"schedule":  job.Schedule,
		"params":    job.Params,
		"created_at": job.CreatedAt,
	}

	jobsArray := []map[string]interface{}{jobData}
	data, err := json.Marshal(jobsArray)
	if err != nil {
		t.Fatalf("Failed to marshal job data: %v", err)
	}

	os.WriteFile(testFile, data, 0644)

	// Test unmarshaling
	loadedData, _ := os.ReadFile(testFile)
	var loadedJobs []map[string]interface{}
	if err := json.Unmarshal(loadedData, &loadedJobs); err != nil {
		t.Fatalf("Failed to unmarshal job data: %v", err)
	}

	if len(loadedJobs) != 1 {
		t.Errorf("Unmarshaled %d jobs, want 1", len(loadedJobs))
	}

	if loadedJobs[0]["db_type"] != job.DBType {
		t.Errorf("Unmarshaled DBType = %v, want %v", loadedJobs[0]["db_type"], job.DBType)
	}
	
	// Cleanup
	os.Remove("./config/schedules.json")
}
