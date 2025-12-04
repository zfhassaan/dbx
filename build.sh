#!/bin/bash

# Build script for Linux/Mac
# This script builds executables for Windows and Linux

APP_NAME="dbx"
BUILD_DIR="dist"
VERSION="${VERSION:-1.0.0}"
LDFLAGS="-s -w -X main.version=${VERSION}"

echo "========================================"
echo "DBX Build Script"
echo "========================================"
echo ""

# Create build directory if it doesn't exist
mkdir -p "${BUILD_DIR}"

# Build Windows executable
echo "[1/3] Building Windows executable (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "${BUILD_DIR}/${APP_NAME}-windows-amd64.exe" main.go
if [ $? -ne 0 ]; then
    echo "ERROR: Windows build failed!"
    exit 1
fi
echo "✓ Windows build complete: ${BUILD_DIR}/${APP_NAME}-windows-amd64.exe"
echo ""

# Build Linux executable (amd64)
echo "[2/3] Building Linux executable (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "${BUILD_DIR}/${APP_NAME}-linux-amd64" main.go
if [ $? -ne 0 ]; then
    echo "ERROR: Linux build failed!"
    exit 1
fi
echo "✓ Linux build complete: ${BUILD_DIR}/${APP_NAME}-linux-amd64"
echo ""

# Build Linux ARM64 executable (optional, for ARM servers)
echo "[3/4] Building Linux ARM64 executable..."
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o "${BUILD_DIR}/${APP_NAME}-linux-arm64" main.go
if [ $? -ne 0 ]; then
    echo "WARNING: Linux ARM64 build failed (this is optional)"
else
    echo "✓ Linux ARM64 build complete: ${BUILD_DIR}/${APP_NAME}-linux-arm64"
fi
echo ""

# Set executable permissions for Linux builds
chmod +x "${BUILD_DIR}/${APP_NAME}-linux-amd64" 2>/dev/null
chmod +x "${BUILD_DIR}/${APP_NAME}-linux-arm64" 2>/dev/null

echo "========================================"
echo "Build Summary:"
echo "  - Windows: ${BUILD_DIR}/${APP_NAME}-windows-amd64.exe"
echo "  - Linux:   ${BUILD_DIR}/${APP_NAME}-linux-amd64"
echo "  - Linux ARM64: ${BUILD_DIR}/${APP_NAME}-linux-arm64"
echo "========================================"
echo ""
echo "All builds completed successfully!"

