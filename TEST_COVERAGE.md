# Test Coverage Documentation

## Quick Start

### Check Test Coverage

**Windows:**
```cmd
coverage.bat
```

**Linux/Mac:**
```bash
./coverage
```

**Or using Make:**
```bash
make coverage
```

**Or directly:**
```bash
go run scripts/coverage.go
```

### Why Coverage Shows 0% with `go test -cover ./...`?

When tests are in a separate `tests/` directory with `*_test` package names, running `go test -cover ./...` measures coverage of the **test packages themselves** (which have no statements), not the source packages.

**Solution**: Use the `-coverpkg` flag to specify which source packages to measure:

```bash
go test -cover ./tests/internal/utils -coverpkg=./internal/utils
```

This tells Go: "Run tests in `./tests/internal/utils` but measure coverage of `./internal/utils`"

### Current Coverage Status

- ✅ **notify**: 100.0% coverage
- ✅ **utils**: 83.1% coverage  
- ✅ **logs**: 76.2% coverage
- ✅ **scheduler**: 50.0% coverage
- ⚠️ **db**: 16.1% coverage (needs more tests)
- ⚠️ **cloud**: 10.0% coverage (needs more tests)

**Average Coverage: 55.9%**

### Generate HTML Coverage Report

```bash
make coverage-html
```

This will:
1. Generate coverage profiles for each module
2. Open an interactive HTML report in your browser
3. Show which lines are covered (green) and which are not (red)

---

## Overview

This test suite provides comprehensive coverage following industry best practices including:
- **Security**: Input validation, path traversal protection, password security
- **Scalability**: Concurrent operations, large file handling
- **Utility**: Core functionality testing
- **OOP Standards**: Encapsulation, abstraction testing
- **DRY Principle**: Reusable test helpers
- **SOLID Principles**: Single responsibility, dependency injection patterns

## Test Files

### 1. `internal/utils/compress_test.go`
**Coverage**: Compression utilities
- ✅ File compression (zip)
- ✅ Folder compression
- ✅ Gzip compression
- ✅ Error handling (non-existent files, invalid paths)
- ✅ Large file handling (scalability)
- ✅ Concurrent compression (scalability)
- ✅ Path traversal security

**Test Cases**: 10

### 2. `internal/db/backup_types_test.go`
**Coverage**: Backup metadata management
- ✅ Backup type constants
- ✅ Metadata path generation
- ✅ Save/Load metadata cycle
- ✅ Corrupted file handling
- ✅ JSON serialization/deserialization
- ✅ Time field preservation

**Test Cases**: 7

### 3. `internal/db/connection_test.go`
**Coverage**: Database connection testing
- ✅ Unsupported database types
- ✅ Empty/invalid parameters
- ✅ SQLite file existence checks
- ✅ Case-insensitive type matching
- ✅ Error handling

**Test Cases**: 6

### 4. `internal/db/sqlite_test.go`
**Coverage**: SQLite backup operations
- ✅ Valid database backup
- ✅ Non-existent database handling
- ✅ Empty path validation
- ✅ Invalid output directory handling

**Test Cases**: 4

### 5. `internal/scheduler/scheduler_test.go`
**Coverage**: Job scheduling system
- ✅ Valid job addition
- ✅ Invalid cron schedule handling
- ✅ Unsupported database type validation
- ✅ Job listing
- ✅ Corrupted config file handling
- ✅ Concurrent job addition (scalability)
- ✅ Closure capture fix verification (critical bug fix)
- ✅ JSON serialization

**Test Cases**: 10

### 6. `internal/cloud/storage_test.go`
**Coverage**: Cloud storage operations
- ✅ Missing AWS CLI detection
- ✅ Empty bucket validation
- ✅ Valid input handling
- ✅ Path traversal security
- ✅ Empty prefix handling

**Test Cases**: 5

### 7. `internal/logs/logger_test.go`
**Coverage**: Logging functionality
- ✅ Successful log entry creation
- ✅ Error logging
- ✅ Multiple log entries
- ✅ Concurrent log writes (scalability)
- ✅ Log content verification

**Test Cases**: 4

### 8. `internal/notify/slack_test.go`
**Coverage**: Slack notifications
- ✅ Valid webhook notification
- ✅ Invalid webhook handling
- ✅ Empty message handling
- ✅ Long message handling (scalability)
- ✅ Special character handling
- ✅ Server error handling
- ✅ Concurrent notifications (scalability)

**Test Cases**: 7

## Test Helpers

### `test_helpers.go`
Reusable test utilities following DRY principle:
- `TestHelper`: Centralized test helper struct
- `CreateTestFile()`: Create test files with content
- `CreateTestDir()`: Create test directories
- `FileExists()`: File existence checking
- `Cleanup()`: Automatic cleanup

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/utils -v

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

## Coverage Goals

- **Unit Tests**: 80%+ coverage for utility functions
- **Integration Tests**: Core database operations
- **Security Tests**: All input validation paths
- **Scalability Tests**: Concurrent operations, large files
- **Error Handling**: All error paths covered

## SOLID Principles in Tests

1. **Single Responsibility**: Each test function tests one specific behavior
2. **Open/Closed**: Tests are extensible without modifying existing tests
3. **Liskov Substitution**: Test helpers can be substituted
4. **Interface Segregation**: Tests use minimal interfaces
5. **Dependency Inversion**: Tests use dependency injection patterns

## Security Test Coverage

See `SECURITY_TEST.md` for detailed security test documentation.

## Future Enhancements

- [ ] Integration tests with real database instances
- [ ] Performance benchmarks
- [ ] Load testing for concurrent operations
- [ ] Fuzzing tests for input validation
- [ ] Mock frameworks for external dependencies

