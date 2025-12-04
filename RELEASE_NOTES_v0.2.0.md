# Release v0.2.0

## 🎉 What's New

This release focuses on **security improvements**, **new features**, and **enhanced reliability**.

## 🔒 Security Fixes

### MySQL Password Security
- **Fixed**: MySQL password is now passed via `MYSQL_PWD` environment variable instead of command line
- **Impact**: Passwords are no longer visible in process list, significantly improving security
- **Affected**: All MySQL backup and restore operations

## ✨ New Features

### CLI Cloud Upload Support
- Added `--upload` flag to all backup commands (mysql, postgres, mongo, sqlite)
- Support for all three cloud providers: AWS S3, Google Cloud Storage, Azure Blob Storage
- Cloud provider selection via `--cloud` flag
- Environment variable support for automatic configuration

**Example:**
```bash
dbx backup mysql --database mydb --upload --cloud s3 --s3-bucket my-backups
```

### Scheduled Backup Cloud Upload
- Automatic cloud upload for scheduled backups
- Configure via CLI flags or environment variables
- Supports all cloud providers

**Example:**
```bash
dbx schedule add --db mysql --database mydb --cron "0 2 * * *" --upload --s3-bucket my-backups
```

### Backup Verification
- Automatic backup file verification after creation
- File size reporting
- Checksum calculation utilities (`internal/utils/verify.go`)
- Validates backup integrity before completion

### Enhanced Selective Restore Validation
- **MySQL**: Improved table extraction with better pattern matching
- **PostgreSQL**: Added `pg_restore --list` check to verify table exists
- **MongoDB**: Enhanced collection validation
- Better error messages when table/collection not found

### Progress Indicators
- Status messages during backup operations
- Elapsed time tracking for long operations
- Better user feedback during backups

## 📝 Improvements

### PostgreSQL Incremental Backups
- Added comprehensive documentation about WAL archiving requirements
- Clear warnings when incremental/differential backup type is selected
- Instructions on how to configure true incremental backups

### Error Handling
- Enhanced error messages for restore operations
- Better validation feedback
- More descriptive error messages

## 🔧 Technical Changes

### New Utilities
- `internal/utils/verify.go` - Backup verification utilities
- `internal/utils/progress.go` - Progress indicator utility

### Code Quality
- Improved security across all database operations
- Enhanced validation and error handling
- Better code organization

## 📦 Installation

### From Source
```bash
git clone https://github.com/zfhassaan/dbx.git
cd dbx
git checkout v0.2.0
go build -o dbx main.go
```

### Download Pre-built Binaries
Visit the [Releases page](https://github.com/zfhassaan/dbx/releases/tag/v0.2.0) to download pre-built binaries for:
- Windows (amd64)
- Linux (amd64, arm64)
- macOS (amd64)

## 🔄 Migration from v0.1.0

**No breaking changes!** This release is fully backward compatible.

### Environment Variables (New)
- `DBX_S3_BUCKET` - S3 bucket for automatic uploads
- `DBX_S3_PREFIX` - S3 prefix/folder path (default: `dbx/`)
- `DBX_CLOUD_PROVIDER` - Default cloud provider (s3, gcs, azure)
- `DBX_AUTO_UPLOAD` - Enable automatic cloud upload for scheduled backups

## 📚 Documentation

- Updated README with new features
- Enhanced CHANGELOG with detailed changes
- Improved inline code documentation

## 🐛 Bug Fixes

- Fixed MySQL password security vulnerability
- Fixed PostgreSQL incremental backup documentation
- Improved error handling in restore operations

## 🙏 Contributors

Thanks to all contributors and users who provided feedback!

## 📄 Full Changelog

See [CHANGELOG.md](CHANGELOG.md) for complete list of changes.

---

**Release Date**: December 2024  
**Git Tag**: [v0.2.0](https://github.com/zfhassaan/dbx/releases/tag/v0.2.0)

