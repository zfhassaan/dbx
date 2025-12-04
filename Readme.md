# DBX - Database Backup Utility

A cross-platform CLI tool to backup and restore multiple database systems (MySQL, PostgreSQL, MongoDB, SQLite). Supports automatic scheduling, compression, cloud storage uploads (AWS S3, GCS, Azure Blob Storage), and detailed logging.

[![Version](https://img.shields.io/badge/version-0.2.0-blue.svg)](https://github.com/zfhassaan/dbx/releases/tag/v0.2.0)
[![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

---

## Features

### Database Support
- **MySQL** - Full, incremental, and differential backups
- **PostgreSQL** - Full, incremental, and differential backups  
- **MongoDB** - Full database backups with compression
- **SQLite** - File-based backups with compression

### Backup & Restore
- **Multiple Backup Types**: Full, incremental, and differential backups
- **Selective Restore**: Restore specific tables (MySQL/PostgreSQL) or collections (MongoDB)
- **Compression**: Automatic compression using gzip/zip
- **Connection Testing**: Verify database connectivity before operations

### Cloud Storage
- **AWS S3** - Upload backups to Amazon S3
- **Google Cloud Storage** - Upload to GCS buckets
- **Azure Blob Storage** - Upload to Azure containers

### Automation & Monitoring
- **Scheduling**: Automated backups using cron syntax
- **Logging**: Comprehensive logging system with configurable log directory
- **Notifications**: Slack webhook integration for backup status
- **Interactive Menu**: User-friendly menu-based interface
- **CLI Commands**: Full command-line interface using Cobra

### Developer Experience
- **Cross-Platform**: Windows, Linux, macOS support
- **Test Suite**: 110+ comprehensive tests with coverage reporting
- **Zero Dependencies**: No external dependencies when built
- **Password Security**: Hidden password input for secure credential entry

---

## Project Structure

```
dbx/
├── cmd/                          # CLI commands (Cobra framework)
│   ├── root.go                   # Root command and banner
│   ├── backup.go                 # Backup command parent
│   ├── mysql.go                  # MySQL backup subcommand
│   ├── postgres.go               # PostgreSQL backup subcommand
│   ├── mongodb.go                # MongoDB backup subcommand
│   ├── sqlite.go                 # SQLite backup subcommand
│   ├── restore.go                # Restore command with subcommands
│   └── schedule.go               # Schedule command (add/list)
├── internal/
│   ├── db/                       # Database operations
│   │   ├── mysql.go              # MySQL backup implementation
│   │   ├── mysql_restore.go     # MySQL restore implementation
│   │   ├── postgres.go           # PostgreSQL backup implementation
│   │   ├── postgres_restore.go  # PostgreSQL restore implementation
│   │   ├── mongodb.go            # MongoDB backup implementation
│   │   ├── mongodb_restore.go   # MongoDB restore implementation
│   │   ├── sqlite.go             # SQLite backup implementation
│   │   ├── sqlite_restore.go    # SQLite restore implementation
│   │   ├── connection.go         # Database connection testing
│   │   └── backup_types.go      # Backup type definitions
│   ├── cloud/                    # Cloud storage handlers
│   │   ├── storage.go            # AWS S3 upload
│   │   ├── gcs.go                # Google Cloud Storage upload
│   │   └── azure.go              # Azure Blob Storage upload
│   ├── scheduler/                # Backup scheduling
│   │   └── scheduler.go          # Cron-based scheduler
│   ├── logs/                     # Logging utility
│   │   └── logger.go             # Custom logger
│   ├── notify/                   # Notifications
│   │   └── slack.go              # Slack webhook integration
│   └── utils/                    # Utilities
│       └── compress.go           # Compression utilities
├── tests/                        # Test suite
│   ├── internal/                 # Tests mirroring internal structure
│   │   ├── cloud/                # Cloud storage tests
│   │   ├── db/                   # Database operation tests
│   │   ├── logs/                 # Logger tests
│   │   ├── notify/               # Notification tests
│   │   ├── scheduler/            # Scheduler tests
│   │   └── utils/                # Utility tests
│   └── test_helpers.go           # Common test utilities
├── scripts/                      # Build and test scripts
│   ├── runtests/                 # Test runner
│   │   └── main.go
│   └── coverage/                 # Coverage analysis
│       └── main.go
├── config/                       # Configuration files
│   └── schedules.json            # Scheduled backup jobs
├── main.go                       # Interactive menu entrypoint
├── Makefile                      # Build automation
├── build.sh                      # Cross-compilation script (Linux/Mac)
├── build.bat                     # Cross-compilation script (Windows)
├── CHANGELOG.md                  # Version history
└── README.md                     # This file
```

---

## Installation

### Prerequisites

- **Go 1.24+** - [Download Go](https://golang.org/dl/)
- **Database Tools** (for the databases you want to backup):
  - MySQL: `mysqldump`, `mysql` client
  - PostgreSQL: `pg_dump`, `pg_restore`, `psql` client
  - MongoDB: `mongodump`, `mongorestore` (MongoDB Database Tools)
  - SQLite: Built-in (no external tools needed)

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
make build-all          # Build for Windows, Linux, and Linux ARM64
make build-windows      # Build Windows executable only
make build-linux        # Build Linux executable only
make build-linux-arm64 # Build Linux ARM64 executable
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

## Usage

### Interactive Menu

Run without arguments to launch the interactive menu:

```bash
./dbx
# or
./dbx.exe
```

The menu provides:
- Database backup options (MySQL, PostgreSQL, MongoDB, SQLite)
- Restore operations
- Connection testing
- Backup scheduling
- Log viewing
- Cloud storage help

### Command-Line Interface

#### Backup Commands

**MySQL Backup:**
```bash
dbx backup mysql --host localhost --user root --password secret --database mydb --out ./backups --type full
```

**PostgreSQL Backup:**
```bash
dbx backup postgres --host localhost --port 5432 --user postgres --password secret --database mydb --out ./backups --type incremental
```

**MongoDB Backup:**
```bash
dbx backup mongo --uri mongodb://localhost:27017 --database mydb --out ./backups
```

**SQLite Backup:**
```bash
dbx backup sqlite --path /path/to/database.db --out ./backups
```

#### Restore Commands

**MySQL Restore:**
```bash
dbx restore mysql --host localhost --user root --password secret --database mydb --file ./backups/backup.sql

# Restore specific table
dbx restore mysql --host localhost --user root --password secret --database mydb --file ./backups/backup.sql --table users
```

**PostgreSQL Restore:**
```bash
dbx restore postgres --host localhost --port 5432 --user postgres --password secret --database mydb --file ./backups/backup.dump

# Restore specific table
dbx restore postgres --host localhost --port 5432 --user postgres --password secret --database mydb --file ./backups/backup.dump --table users
```

**MongoDB Restore:**
```bash
dbx restore mongo --uri mongodb://localhost:27017 --database mydb --file ./backups/mydb_backup

# Restore specific collection
dbx restore mongo --uri mongodb://localhost:27017 --database mydb --file ./backups/mydb_backup --collection users
```

**SQLite Restore:**
```bash
dbx restore sqlite --path /path/to/restored.db --file ./backups/backup.db
```

#### Scheduling Commands

**Add Scheduled Backup:**
```bash
dbx schedule add --db mysql --host localhost --user root --password secret --database mydb --out ./backups --cron "0 2 * * *"
```

**List Scheduled Backups:**
```bash
dbx schedule list
```

**Cron Examples:**
- `0 2 * * *` - Daily at 2 AM
- `0 */6 * * *` - Every 6 hours
- `0 0 * * 0` - Weekly on Sunday at midnight
- `0 0 1 * *` - Monthly on the 1st at midnight

#### Help

Get help for any command:
```bash
dbx --help
dbx backup --help
dbx backup mysql --help
dbx restore --help
dbx schedule --help
```

---

## Cloud Storage Setup

### AWS S3

1. Install AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
2. Configure credentials:
   ```bash
   aws configure
   ```
3. Upload backups:
   - Use the interactive menu and select "Upload to AWS S3" after backup
   - Or set environment variables:
     ```bash
     export DBX_S3_BUCKET=my-backup-bucket
     export DBX_S3_PREFIX=dbx/
     ```

### Google Cloud Storage

1. Install Google Cloud SDK: https://cloud.google.com/sdk/docs/install
2. Authenticate:
   ```bash
   gcloud auth login
   ```
3. Upload backups using the interactive menu

### Azure Blob Storage

1. Install Azure CLI: https://docs.microsoft.com/cli/azure/install-azure-cli
2. Authenticate:
   ```bash
   az login
   ```
3. Upload backups using the interactive menu

---

## Notifications

### Slack Integration

Add your Slack webhook URL to environment variables:

```bash
export SLACK_WEBHOOK=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

Or create a `.env` file:
```
SLACK_WEBHOOK=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

Backup operations will automatically send notifications to Slack with:
- Backup status (SUCCESS/FAILED)
- Database name
- Duration
- Host and user information
- Error details (if any)

---

## Testing

### Run All Tests

```bash
make test
# or
go run scripts/runtests/main.go
```

### Run Tests for Specific Module

```bash
go test -v ./tests/internal/db
go test -v ./tests/internal/cloud
go test -v ./tests/internal/scheduler
```

### Test Coverage

```bash
make coverage
# or
go run scripts/coverage/main.go
```

### Current Test Coverage

- **cloud**: 20.0% coverage
- **db**: 42.4% coverage
- **logs**: 95.2% coverage
- **notify**: 100.0% coverage
- **scheduler**: 61.1% coverage
- **utils**: 83.1% coverage

**Total**: 110+ tests passing across all modules

---

## Security & Performance

### Security
- **Password Hiding**: Secure password input using `golang.org/x/term`
- **Environment Variables**: Credentials stored in environment variables, not in code
- **No Plaintext Storage**: Passwords never stored in plaintext
- **Secure File Permissions**: Log files and configs use appropriate permissions

### Performance
- **Native Tools**: Uses native database tools (`mysqldump`, `pg_dump`, etc.) for optimal performance
- **Streaming Compression**: Low memory footprint with streaming compression
- **Efficient Scheduling**: Lightweight cron-based scheduler
- **Parallel Operations**: Supports concurrent backup operations

---

## Configuration

### Log Directory

Set custom log directory:
```bash
export DBX_LOG_DIR=/var/log/dbx
```

Default: `./logs`

### Scheduled Backups

Scheduled backups are stored in `config/schedules.json`. This file is automatically created and managed by the scheduler.

---

## Troubleshooting

### Database Tools Not Found

If you see errors about missing database tools:

**MySQL:**
```bash
# Ubuntu/Debian
sudo apt install mysql-client

# macOS
brew install mysql-client

# Windows
# Download MySQL Installer from mysql.com
```

**PostgreSQL:**
```bash
# Ubuntu/Debian
sudo apt install postgresql-client

# macOS
brew install postgresql

# Windows
# Download from postgresql.org
```

**MongoDB:**
```bash
# Ubuntu/Debian
sudo apt install mongodb-database-tools

# macOS
brew install mongodb-database-tools

# Windows
# Download MongoDB Tools from mongodb.com
```

### Permission Errors

Ensure you have:
- Read access to source databases
- Write access to backup directories
- Execute permissions for database tools

### Cloud Upload Failures

- Verify CLI tools are installed and configured
- Check credentials and permissions
- Ensure bucket/container names are correct
- Verify network connectivity

---

## TODO

- [ ] Add web-based dashboard for backup logs
- [ ] Add encryption for backup files
- [ ] Support for differential backup in MongoDB
- [ ] Retry failed backups automatically
- [ ] Add backup verification/checksum
- [ ] Support for backup retention policies
- [ ] Add email notifications

---

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Setup

```bash
# Clone repository
git clone https://github.com/zfhassaan/dbx.git
cd dbx

# Install dependencies
go mod download

# Run tests
make test

# Check coverage
make coverage

# Build
make build
```

---

## License

MIT License — feel free to fork and contribute!

See [LICENSE](LICENSE) file for details.

---

## Documentation

- [CHANGELOG.md](CHANGELOG.md) - Version history
- [TEST_COVERAGE.md](TEST_COVERAGE.md) - Test coverage details
- [TEST_RUNNER.md](TEST_RUNNER.md) - Test runner documentation
- [TEST_SUMMARY.md](TEST_SUMMARY.md) - Test suite summary

---

## Contributors

Built by [zfhassaan](https://github.com/zfhassaan) — PRs welcome!

---

## Links

- **Project URL**: https://roadmap.sh/projects/database-backup-utility
- **GitHub Repository**: https://github.com/zfhassaan/dbx
- **Releases**: https://github.com/zfhassaan/dbx/releases

---

## Version

Current version: **v0.2.0**

For version history, see [CHANGELOG.md](CHANGELOG.md).
