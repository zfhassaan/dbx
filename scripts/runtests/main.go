package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// TestResult holds test results for a module
type TestResult struct {
	Module    string
	Passed    int
	Failed    int
	HasErrors bool
	Output    string
}

// Colors for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
)

func main() {
	fmt.Println("==========================================")
	fmt.Println("  DBX Project Test Suite Runner")
	fmt.Println("==========================================")
	fmt.Println()

	// Step 1: Check if code compiles
	fmt.Println("📦 Step 1: Checking code compilation...")
	if err := checkCompilation(); err != nil {
		fmt.Println()
		fmt.Println(ColorRed + "❌ COMPILATION ERROR DETECTED" + ColorReset)
		fmt.Println()
		fmt.Println("Please fix compilation errors before running tests.")
		fmt.Println("Run 'go build ./...' to see detailed errors.")
		fmt.Println()
		fmt.Println("Press ENTER to exit...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		os.Exit(1)
	}
	fmt.Println(ColorGreen + "✅ Code compiles successfully" + ColorReset)
	fmt.Println()

	// Step 2: Run tests
	fmt.Println("🧪 Step 2: Running all tests...")
	fmt.Println()

	// Test modules from tests/ directory
	var modules []struct {
		name string
		path string
	}

	if dirExists("tests") {
		testDirs, _ := filepath.Glob("tests/internal/*")
		for _, dir := range testDirs {
			if info, err := os.Stat(dir); err == nil && info.IsDir() {
				moduleName := filepath.Base(dir)
				modules = append(modules, struct {
					name string
					path string
				}{moduleName, "./" + dir})
			}
		}
	} else {
		fmt.Println(ColorYellow + "⚠️  Warning: tests/ directory not found" + ColorReset)
		fmt.Println()
		os.Exit(1)
	}

	var results []TestResult
	totalPassed := 0
	totalFailed := 0
	var failedModules []string

	for _, module := range modules {
		result := runTests(module.name, module.path)
		results = append(results, result)
		totalPassed += result.Passed
		totalFailed += result.Failed

		if result.HasErrors {
			failedModules = append(failedModules, module.name)
		}
	}

	// Step 3: Print summary
	fmt.Println("==========================================")
	fmt.Println("  Test Summary")
	fmt.Println("==========================================")
	fmt.Printf("Total Passed: %s%d%s\n", ColorGreen, totalPassed, ColorReset)
	fmt.Printf("Total Failed: %s%d%s\n", ColorRed, totalFailed, ColorReset)
	fmt.Println()

	// Print detailed results
	fmt.Println("Detailed Results:")
	fmt.Println("-----------------------------------")
	for _, result := range results {
		status := "✅"
		if result.HasErrors {
			status = "❌"
		}
		fmt.Printf("%s %s: %d passed, %d failed\n", status, result.Module, result.Passed, result.Failed)
	}
	fmt.Println()

	if len(failedModules) > 0 {
		fmt.Println(ColorRed + "❌ Failed Modules:" + ColorReset)
		for _, module := range failedModules {
			fmt.Printf("  %s  - %s%s\n", ColorRed, module, ColorReset)
		}
		fmt.Println()
		fmt.Println("Please fix the failing tests or compilation errors before proceeding.")
		fmt.Println()
		fmt.Println("Press ENTER to exit...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		os.Exit(1)
	} else if totalPassed == 0 && totalFailed == 0 {
		fmt.Println(ColorYellow + "⚠️  No tests were executed. Check if test files exist in tests/ directory." + ColorReset)
		fmt.Println()
	} else {
		fmt.Println(ColorGreen + "✅ All tests passed!" + ColorReset)
		fmt.Println()
	}
}

func checkCompilation() error {
	cmd := exec.Command("go", "build", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runTests(moduleName, modulePath string) TestResult {
	result := TestResult{
		Module: moduleName,
	}

	fmt.Printf("Testing module: %s\n", moduleName)
	fmt.Println("-----------------------------------")

	// Check if test files exist
	testFiles, _ := filepath.Glob(modulePath + "/*_test.go")
	if len(testFiles) == 0 {
		fmt.Println(ColorYellow + fmt.Sprintf("⚠️  %s: No test files found", moduleName) + ColorReset)
		fmt.Println()
		return result
	}

	cmd := exec.Command("go", "test", "-v", modulePath)
	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	// Check for compilation errors
	if strings.Contains(result.Output, "build failed") || strings.Contains(result.Output, "compilation error") {
		fmt.Println(ColorRed + "❌ COMPILATION ERROR in " + moduleName + ColorReset)
		fmt.Println(result.Output)
		result.HasErrors = true
		fmt.Println()
		return result
	}

	if err != nil {
		result.HasErrors = true
	}

	// Parse output for pass/fail counts
	passRegex := regexp.MustCompile(`--- PASS:`)
	failRegex := regexp.MustCompile(`--- FAIL:`)

	passMatches := passRegex.FindAllString(result.Output, -1)
	failMatches := failRegex.FindAllString(result.Output, -1)

	result.Passed = len(passMatches)
	result.Failed = len(failMatches)

	// Check if no tests ran
	if result.Passed == 0 && result.Failed == 0 && !strings.Contains(result.Output, "no test files") {
		// Tests might have compilation issues
		if strings.Contains(result.Output, "undefined") || strings.Contains(result.Output, "cannot") {
			fmt.Println(ColorRed + "❌ " + moduleName + ": Compilation errors detected" + ColorReset)
			// Show first few error lines
			lines := strings.Split(result.Output, "\n")
			errorCount := 0
			for _, line := range lines {
				if (strings.Contains(line, "error:") || strings.Contains(line, "undefined") || strings.Contains(line, "cannot")) && errorCount < 3 {
					fmt.Println(ColorRed + line + ColorReset)
					errorCount++
				}
			}
			result.HasErrors = true
		} else {
			fmt.Println(ColorYellow + fmt.Sprintf("⚠️  %s: No tests executed", moduleName) + ColorReset)
		}
	} else if result.HasErrors || result.Failed > 0 {
		result.HasErrors = true
		// Print failure details
		lines := strings.Split(result.Output, "\n")
		for i, line := range lines {
			if strings.Contains(line, "FAIL:") && i+1 < len(lines) {
				fmt.Println(ColorRed + line + ColorReset)
				// Print next few lines for context
				for j := 1; j <= 3 && i+j < len(lines); j++ {
					fmt.Println(lines[i+j])
				}
			}
		}
		fmt.Println(ColorRed + fmt.Sprintf("❌ %s: %d passed, %d failed", moduleName, result.Passed, result.Failed) + ColorReset)
	} else {
		fmt.Println(ColorGreen + fmt.Sprintf("✅ %s: %d passed, %d failed", moduleName, result.Passed, result.Failed) + ColorReset)
	}
	fmt.Println()

	return result
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

