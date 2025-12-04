package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Colors for terminal output
const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorCyan  = "\033[36m"
)

func main() {
	fmt.Println("==========================================")
	fmt.Println("  DBX Test Coverage Report")
	fmt.Println("==========================================")
	fmt.Println()

	// Map test directories to source packages
	testToSource := map[string]string{
		"./tests/internal/cloud":     "./internal/cloud",
		"./tests/internal/db":        "./internal/db",
		"./tests/internal/logs":      "./internal/logs",
		"./tests/internal/notify":    "./internal/notify",
		"./tests/internal/scheduler": "./internal/scheduler",
		"./tests/internal/utils":     "./internal/utils",
	}

	totalCoverage := 0.0
	moduleCount := 0

	fmt.Println("Module Coverage:")
	fmt.Println("-----------------------------------")

	for testDir, sourcePkg := range testToSource {
		moduleName := filepath.Base(testDir)
		
		// Run test with coverage for the source package
		cmd := exec.Command("go", "test", "-cover", testDir, "-coverpkg", sourcePkg)
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			fmt.Printf("%s❌ %s: Error running tests%s\n", ColorRed, moduleName, ColorReset)
			fmt.Println(string(output))
			continue
		}

		// Parse coverage percentage from output
		outputStr := string(output)
		coverage := parseCoverage(outputStr)
		
		if coverage >= 0 {
			totalCoverage += coverage
			moduleCount++
			color := ColorGreen
			if coverage < 50 {
				color = ColorRed
			} else if coverage < 80 {
				color = ColorCyan
			}
			fmt.Printf("%s✅ %s: %.1f%% coverage%s\n", color, moduleName, coverage, ColorReset)
		} else {
			fmt.Printf("%s⚠️  %s: Could not determine coverage%s\n", ColorCyan, moduleName, ColorReset)
		}
	}

	fmt.Println()
	fmt.Println("==========================================")
	fmt.Println("  Summary")
	fmt.Println("==========================================")
	
	if moduleCount > 0 {
		avgCoverage := totalCoverage / float64(moduleCount)
		color := ColorGreen
		if avgCoverage < 50 {
			color = ColorRed
		} else if avgCoverage < 80 {
			color = ColorCyan
		}
		fmt.Printf("Average Coverage: %s%.1f%%%s\n", color, avgCoverage, ColorReset)
		fmt.Printf("Modules Tested: %d\n", moduleCount)
	} else {
		fmt.Println("No coverage data available")
	}
	fmt.Println()

	// Generate detailed coverage report
	fmt.Println("Generating detailed coverage report...")
	fmt.Println("Run: go tool cover -html=coverage.out")
	
	// Generate combined coverage profile
	generateCombinedCoverage(testToSource)
}

func parseCoverage(output string) float64 {
	// Look for "coverage: XX.X% of statements" pattern
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "coverage:") {
			// Extract percentage
			parts := strings.Split(line, "coverage:")
			if len(parts) > 1 {
				coverageStr := strings.TrimSpace(parts[1])
				coverageStr = strings.TrimSuffix(coverageStr, "% of statements")
				coverageStr = strings.TrimSpace(coverageStr)
				
				var coverage float64
				if _, err := fmt.Sscanf(coverageStr, "%f", &coverage); err == nil {
					return coverage
				}
			}
		}
	}
	return -1
}

func generateCombinedCoverage(testToSource map[string]string) {
	// Create a combined coverage profile
	profiles := []string{}
	
	for testDir, sourcePkg := range testToSource {
		profileFile := fmt.Sprintf("coverage_%s.out", filepath.Base(testDir))
		
		cmd := exec.Command("go", "test", "-coverprofile", profileFile, testDir, "-coverpkg", sourcePkg)
		cmd.Run() // Ignore errors - some packages might not have coverage
		
		if _, err := os.Stat(profileFile); err == nil {
			profiles = append(profiles, profileFile)
		}
	}

	if len(profiles) > 0 {
		// Combine profiles using go tool cover
		fmt.Printf("Coverage profiles generated: %s\n", strings.Join(profiles, ", "))
		fmt.Println("To view HTML report: go tool cover -html=" + profiles[0])
	}
}

