#!/bin/bash

# Comprehensive Test Runner for DBX Project
# This script:
# 1. Checks if code compiles (stops on error)
# 2. Runs all tests
# 3. Reports pass/fail counts per module
# 4. Stops on first compilation error

set -e  # Exit on error

echo "=========================================="
echo "  DBX Project Test Suite Runner"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Check if code compiles
echo "📦 Step 1: Checking code compilation..."
if ! go build ./... 2>&1; then
    echo ""
    echo -e "${RED}❌ COMPILATION ERROR DETECTED${NC}"
    echo ""
    echo "Please fix compilation errors before running tests."
    echo "Run 'go build ./...' to see detailed errors."
    exit 1
fi
echo -e "${GREEN}✅ Code compiles successfully${NC}"
echo ""

# Step 2: Run tests and collect results
echo "🧪 Step 2: Running all tests..."
echo ""

# Track results
TOTAL_PASSED=0
TOTAL_FAILED=0
FAILED_MODULES=()

# Function to run tests for a module and parse results
run_module_tests() {
    local module=$1
    local module_path=$2
    
    echo "Testing module: $module"
    echo "-----------------------------------"
    
    # Run tests and capture output
    local output=$(go test -v $module_path 2>&1)
    local exit_code=$?
    
    # Parse test results
    local passed=$(echo "$output" | grep -c "--- PASS:" || echo "0")
    local failed=$(echo "$output" | grep -c "--- FAIL:" || echo "0")
    
    TOTAL_PASSED=$((TOTAL_PASSED + passed))
    TOTAL_FAILED=$((TOTAL_FAILED + failed))
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ $module: $passed passed, $failed failed${NC}"
    else
        echo -e "${RED}❌ $module: $passed passed, $failed failed${NC}"
        FAILED_MODULES+=("$module")
        echo "$output" | grep -A 5 "FAIL:"
    fi
    echo ""
}

# Test each module
run_module_tests "utils" "./internal/utils"
run_module_tests "db/backup_types" "./internal/db"
run_module_tests "db/connection" "./internal/db"
run_module_tests "db/sqlite" "./internal/db"
run_module_tests "scheduler" "./internal/scheduler"
run_module_tests "cloud/storage" "./internal/cloud"
run_module_tests "logs" "./internal/logs"
run_module_tests "notify" "./internal/notify"

# Also test from tests directory if it exists
if [ -d "tests" ]; then
    echo "Testing from tests/ directory..."
    echo "-----------------------------------"
    for test_dir in tests/internal/*/; do
        if [ -d "$test_dir" ]; then
            module_name=$(basename "$test_dir")
            run_module_tests "tests/$module_name" "./$test_dir"
        fi
    done
fi

# Step 3: Summary
echo "=========================================="
echo "  Test Summary"
echo "=========================================="
echo -e "Total Passed: ${GREEN}$TOTAL_PASSED${NC}"
echo -e "Total Failed: ${RED}$TOTAL_FAILED${NC}"
echo ""

if [ ${#FAILED_MODULES[@]} -gt 0 ]; then
    echo -e "${RED}❌ Failed Modules:${NC}"
    for module in "${FAILED_MODULES[@]}"; do
        echo -e "  ${RED}  - $module${NC}"
    done
    echo ""
    echo "Please fix the failing tests before proceeding."
    exit 1
else
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
fi

