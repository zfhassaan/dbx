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

