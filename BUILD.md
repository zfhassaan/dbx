# Build Instructions

This document provides instructions for building DBX executables for different platforms.

## Prerequisites

- Go 1.23 or higher installed
- Git (for cloning the repository)

## Quick Build

### Windows

**Option 1: Using batch script**
```cmd
build.bat
```

**Option 2: Using Make (if you have Make installed)**
```cmd
make build-windows
```

**Option 3: Manual build**
```cmd
go build -o dist/dbx-windows-amd64.exe main.go
```

### Linux

**Option 1: Using shell script**
```bash
chmod +x build.sh
./build.sh
```

**Option 2: Using Make**
```bash
make build-linux
```

**Option 3: Manual build**
```bash
go build -o dist/dbx-linux-amd64 main.go
```

## Cross-Compilation

Go supports cross-compilation natively. You can build for any platform from any platform.

### Build Windows executable from Linux/Mac

```bash
GOOS=windows GOARCH=amd64 go build -o dist/dbx-windows-amd64.exe main.go
```

### Build Linux executable from Windows

```cmd
set GOOS=linux
set GOARCH=amd64
go build -o dist/dbx-linux-amd64 main.go
```

### Build for Linux ARM64 (for ARM servers)

```bash
GOOS=linux GOARCH=arm64 go build -o dist/dbx-linux-arm64 main.go
```

## Build All Platforms

### Using Makefile

```bash
make build-all
```

This will create:
- `dist/dbx-windows-amd64.exe` - Windows executable
- `dist/dbx-linux-amd64` - Linux executable (amd64)

### Using Build Scripts

**Windows:**
```cmd
build.bat
```

**Linux/Mac:**
```bash
./build.sh
```

## Build Flags

The build scripts use the following flags for optimization:
- `-s`: Omit symbol table and debug information
- `-w`: Omit DWARF symbol table

These flags reduce the executable size significantly.

## Output Location

All built executables are placed in the `dist/` directory.

## Testing the Build

After building, you can test the executable:

**Windows:**
```cmd
dist\dbx-windows-amd64.exe --help
```

**Linux:**
```bash
./dist/dbx-linux-amd64 --help
```

## Distribution

To distribute the application:

1. Build the executable for your target platform
2. Copy the executable to the target machine
3. Ensure the target machine has the required database client tools installed:
   - MySQL: `mysqldump` and `mysql`
   - PostgreSQL: `pg_dump` and `pg_restore`
   - MongoDB: `mongodump` and `mongorestore`
   - SQLite: No additional tools needed (built-in)

4. For cloud storage uploads, ensure the appropriate CLI tools are installed:
   - AWS S3: `aws` CLI
   - Google Cloud Storage: `gsutil`
   - Azure Blob Storage: `az` CLI

## Troubleshooting

### Build fails with "cannot find package"

Run:
```bash
go mod download
go mod tidy
```

### Executable is too large

The build scripts already include size optimization flags. If you need even smaller binaries, you can use UPX compression (optional):

```bash
upx --best dist/dbx-linux-amd64
```

### Permission denied on Linux

Make sure the executable has execute permissions:
```bash
chmod +x dist/dbx-linux-amd64
```

