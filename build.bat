@echo off
REM Build script for Windows
REM This script builds executables for Windows and Linux

set APP_NAME=dbx
set BUILD_DIR=dist

echo ========================================
echo DBX Build Script for Windows
echo ========================================
echo.

REM Create build directory if it doesn't exist
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

echo [1/3] Building Windows executable (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o %BUILD_DIR%\%APP_NAME%-windows-amd64.exe main.go
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Windows build failed!
    exit /b 1
)
echo ✓ Windows build complete: %BUILD_DIR%\%APP_NAME%-windows-amd64.exe
echo.

echo [2/3] Building Linux executable (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o %BUILD_DIR%\%APP_NAME%-linux-amd64 main.go
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Linux build failed!
    exit /b 1
)
echo ✓ Linux build complete: %BUILD_DIR%\%APP_NAME%-linux-amd64
echo.

echo [3/3] Build Summary:
echo   - Windows: %BUILD_DIR%\%APP_NAME%-windows-amd64.exe
echo   - Linux:   %BUILD_DIR%\%APP_NAME%-linux-amd64
echo.
echo ========================================
echo All builds completed successfully!
echo ========================================
pause

