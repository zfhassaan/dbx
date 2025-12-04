@echo off
REM Comprehensive Test Runner for DBX Project (Windows)
REM This script:
REM 1. Checks if code compiles (stops on error)
REM 2. Runs all tests
REM 3. Reports pass/fail counts per module
REM 4. Stops on first compilation error

echo ==========================================
echo   DBX Project Test Suite Runner
echo ==========================================
echo.

REM Step 1: Check if code compiles
echo [Step 1] Checking code compilation...
go build ./... 2>nul
if errorlevel 1 (
    echo.
    echo [ERROR] COMPILATION ERROR DETECTED
    echo.
    echo Please fix compilation errors before running tests.
    echo Run 'go build ./...' to see detailed errors.
    exit /b 1
)
echo [OK] Code compiles successfully
echo.

REM Step 2: Run tests
echo [Step 2] Running all tests...
echo.

set TOTAL_PASSED=0
set TOTAL_FAILED=0
set FAILED_COUNT=0

REM Function to test a module
:test_module
setlocal
set MODULE_NAME=%~1
set MODULE_PATH=%~2

echo Testing module: %MODULE_NAME%
echo -----------------------------------

go test -v %MODULE_PATH% > test_output.tmp 2>&1
set TEST_EXIT=%ERRORLEVEL%

REM Count passed and failed
findstr /C:"--- PASS:" test_output.tmp >nul 2>&1
if errorlevel 1 (set MODULE_PASSED=0) else (
    for /f %%i in ('findstr /C:"--- PASS:" test_output.tmp ^| find /c /v ""') do set MODULE_PASSED=%%i
)

findstr /C:"--- FAIL:" test_output.tmp >nul 2>&1
if errorlevel 1 (set MODULE_FAILED=0) else (
    for /f %%i in ('findstr /C:"--- FAIL:" test_output.tmp ^| find /c /v ""') do set MODULE_FAILED=%%i
)

set /a TOTAL_PASSED+=%MODULE_PASSED%
set /a TOTAL_FAILED+=%MODULE_FAILED%

if %TEST_EXIT% equ 0 (
    echo [OK] %MODULE_NAME%: %MODULE_PASSED% passed, %MODULE_FAILED% failed
) else (
    echo [FAIL] %MODULE_NAME%: %MODULE_PASSED% passed, %MODULE_FAILED% failed
    set /a FAILED_COUNT+=1
    findstr /C:"FAIL:" test_output.tmp
)

del test_output.tmp 2>nul
echo.
endlocal
goto :eof

REM Test each module
call :test_module "utils" "./internal/utils"
call :test_module "db" "./internal/db"
call :test_module "scheduler" "./internal/scheduler"
call :test_module "cloud" "./internal/cloud"
call :test_module "logs" "./internal/logs"
call :test_module "notify" "./internal/notify"

REM Step 3: Summary
echo ==========================================
echo   Test Summary
echo ==========================================
echo Total Passed: %TOTAL_PASSED%
echo Total Failed: %TOTAL_FAILED%
echo.

if %FAILED_COUNT% gtr 0 (
    echo [ERROR] Some tests failed. Please fix them before proceeding.
    exit /b 1
) else (
    echo [OK] All tests passed!
    exit /b 0
)

