# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2024-12-XX

### Added
- **Database Support**: Full support for MySQL, PostgreSQL, MongoDB, and SQLite
- **Backup Types**: Support for full, incremental, and differential backups
- **Compression**: Automatic compression of backup files (gzip/zip)
- **Cloud Storage**: Integration with AWS S3, Google Cloud Storage, and Azure Blob Storage
- **Restore Operations**: Full database restore and selective table/collection restore
- **Scheduling**: Automated backup scheduling using cron syntax
- **Logging**: Comprehensive logging system for all backup and restore operations
- **Notifications**: Slack integration for backup status notifications
- **CLI Interface**: Interactive menu-based interface and command-line commands
- **Cross-Platform**: Support for Windows, Linux, and macOS
- **Test Suite**: Comprehensive test coverage with 110+ tests across all modules
- **Documentation**: Complete README with installation and usage instructions

### Features
- Interactive menu system for easy database operations
- Command-line interface using Cobra framework
- Password hiding for secure credential input
- Automatic detection of required database tools
- Helpful error messages and installation guides
- Test coverage reporting and analysis tools

### Technical
- Go 1.24+ required
- Modular architecture following SOLID principles
- Comprehensive error handling
- Cross-platform compatibility
- Zero external dependencies when built

[0.1.0]: https://github.com/zfhassaan/dbx/releases/tag/v0.1.0

## [0.2.0] - 2024-12-XX

### Security
- **Fixed MySQL Password Security**: Changed from command-line password passing to `MYSQL_PWD` environment variable for secure credential handling
- Passwords are no longer visible in process list

### Added
- **CLI Cloud Upload**: Added `--upload` flag to all backup commands (mysql, postgres, mongo, sqlite)
- **Cloud Provider Selection**: Support for selecting cloud provider via `--cloud` flag (s3, gcs, azure)
- **Scheduled Backup Cloud Upload**: Automatic cloud upload support for scheduled backups
- **Backup Verification**: Added backup file verification with size reporting and checksum calculation utilities
- **Selective Restore Validation**: Enhanced validation for table/collection restore operations
- **Progress Indicators**: Added status messages and progress feedback for long-running operations

### Improved
- **PostgreSQL Incremental Backups**: Added comprehensive documentation and warnings about WAL archiving requirements
- **Error Messages**: Enhanced error messages for restore operations with better validation feedback
- **MySQL Table Restore**: Improved table extraction logic with better pattern matching

### Technical
- Created `internal/utils/verify.go` for backup verification utilities
- Created `internal/utils/progress.go` for progress indication
- Enhanced scheduler with cloud upload integration
- Improved security across all database operations

### Breaking Changes
None. All changes are backward compatible.

[0.2.0]: https://github.com/zfhassaan/dbx/releases/tag/v0.2.0

