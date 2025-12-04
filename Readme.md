Project URL:
https://roadmap.sh/projects/database-backup-utility
---

```md
# 🗄️ DBBackupCLI - Database Backup Utility

A cross-platform CLI tool to backup and restore multiple database systems (MySQL, PostgreSQL, MongoDB, SQLite, etc.). Supports automatic scheduling, compression, cloud storage uploads (AWS/GCS/Azure), and detailed logging.

---

## 🚀 Features

- ✅ Supports MySQL, PostgreSQL, MongoDB, SQLite
- ✅ Full, Incremental, and Differential backups
- ✅ Compressed output (gzip/zip)
- ✅ Local and cloud storage (AWS S3, GCS, Azure Blob)
- ✅ Restore from backups (with table/collection-level granularity)
- ✅ Slack notifications for backup status
- ✅ Simple CLI interface with help/docs
- ✅ Cross-platform (Windows/Linux/macOS)
- ✅ Zero external dependencies when built

---

```

## 🏗 Folder Structure
```
dbx/
├── cmd/                    # CLI commands (cobra/urfave structure)
│   ├── root.go             # Entry command
│   ├── backup.go           # `backup` command
│   ├── restore.go          # `restore` command
│   └── schedule.go         # `schedule` command
├── internal/
│   ├── db/                 # Database connectors and dumpers
│   │   ├── mysql.go
│   │   ├── postgres.go
│   │   ├── mongo.go
│   │   └── sqlite.go
│   ├── storage/            # Local + cloud storage handlers
│   │   ├── local.go
│   │   ├── s3.go
│   │   ├── gcs.go
│   │   └── azure.go
│   ├── notifier/           # Slack notifications, etc.
│   │   └── slack.go
│   ├── logger/             # Custom logging utility
│   │   └── logger.go
│   ├── config/             # Config loader (yaml/env)
│   │   └── config.go
│   └── utils/              # Misc utilities (compression, validation)
│       └── compress.go
├── schedules/              # Optional cron job configs
│   └── sample.cron
├── .env                    # Environment variables
├── config.yaml             # User config
├── Dockerfile
├── Makefile
├── go.mod
└── main.go                 # CLI entrypoint

````

---

## ⚙️ Installation

### Option 1: Build from Source

**On Windows:**
```bash
git clone https://github.com/zfhassaan/dbx.git
cd dbx
go build -o dbx.exe main.go
```

**On Linux/Mac:**
```bash
git clone https://github.com/zfhassaan/dbx.git
cd dbx
go build -o dbx main.go
```

### Option 2: Cross-Compile for Multiple Platforms

**Using Makefile (Linux/Mac):**
```bash
make build-all          # Build for Windows and Linux
make build-windows      # Build Windows executable only
make build-linux        # Build Linux executable only
```

**Using Build Scripts:**

**Windows:**
```cmd
build.bat
```

**Linux/Mac:**
```bash
chmod +x build.sh
./build.sh
```

**Manual Cross-Compilation:**
```bash
# Windows executable
GOOS=windows GOARCH=amd64 go build -o dist/dbx-windows-amd64.exe main.go

# Linux executable
GOOS=linux GOARCH=amd64 go build -o dist/dbx-linux-amd64 main.go

# Linux ARM64 (for ARM servers)
GOOS=linux GOARCH=arm64 go build -o dist/dbx-linux-arm64 main.go
```

All executables will be created in the `dist/` directory.

---

## 🧪 Usage

```bash
./dbbackup backup --db mysql --host localhost --user root --password secret --output /backups
./dbbackup restore --file /backups/backup-2025-05-02.sql.gz
./dbbackup schedule --config ./config.yaml
```

Use `--help` with any command to view full options.

---

## ☁️ Cloud Storage Setup

Set your `.env` or `config.yaml` with your preferred cloud credentials:

```yaml
cloud:
  provider: s3
  aws_access_key: YOUR_KEY
  aws_secret_key: YOUR_SECRET
  bucket: my-db-backups
```

---

## 🛡️ Security & Performance Notes

* Credentials never stored in plaintext; always use `.env` or config files with proper permissions
* Backup operations use native tools (`mysqldump`, `pg_dump`, etc.) where appropriate for speed
* Compression handled via Go's streaming writers (low memory footprint)

---

## 🔔 Notifications

Add your Slack webhook in `.env`:

```
SLACK_WEBHOOK=https://hooks.slack.com/services/...
```

---

## 📦 TODO

* [ ] Add web-based dashboard for backup logs
* [ ] Add encryption for backup files
* [ ] Support for differential backup in Mongo
* [ ] Retry failed backups automatically

---

## 📄 License

MIT License — feel free to fork and contribute!

---

## ✨ Contributors
```aiexclude
Built by zfhassaan — PRs welcome!
```



