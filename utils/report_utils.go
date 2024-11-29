package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// File represents a file in the coverage report
type File struct {
	Name            string  `json:"name"`
	Path            string  `json:"path"`
	CoveredLines    int     `json:"coveredLines"`
	ExecutableLines int     `json:"executableLines"`
	LineCoverage    float64 `json:"lineCoverage"`
}

// Target represents a target in the coverage report
type Target struct {
	Name            string  `json:"name"`
	CoveredLines    int     `json:"coveredLines"`
	ExecutableLines int     `json:"executableLines"`
	LineCoverage    float64 `json:"lineCoverage"`
	Files           []File  `json:"files"`
}

// CoverageReport represents the structure of the coverage report JSON
type CoverageReport struct {
	Targets []Target `json:"targets"`
}

// ProcessCoverage processes a coverage report file, filters out specified files, and recalculates the line coverage
// for a given target. It returns the updated coverage percentage and an error if the processing fails.
// Parameters:
// - filePath: The path to the coverage report file.
// - target: The target within the coverage report to process.
// - excludeFiles: A map of file names to be excluded from the coverage calculation.
func ProcessCoverage(filePath, target string, excludeFiles map[string]struct{}) (float64, error) {
	var coverage = 0.0

	data, err := os.ReadFile(filePath)
	if err != nil {
		return coverage, fmt.Errorf("failed to read coverage file: %w", err)
	}

	var report CoverageReport
	if err := json.Unmarshal(data, &report); err != nil {
		return coverage, fmt.Errorf("failed to parse coverage JSON: %w", err)
	}

	// Find the target
	var targetData *Target
	for i := range report.Targets {
		if report.Targets[i].Name == target {
			targetData = &report.Targets[i]
			break
		}
	}
	if targetData == nil {
		return coverage, fmt.Errorf("target %s not found in coverage report", target)
	}

	// Filter files
	var filteredFiles []File
	for _, file := range targetData.Files {
		if _, excluded := excludeFiles[file.Name]; !excluded {
			filteredFiles = append(filteredFiles, file)
		}
	}
	targetData.Files = filteredFiles

	// Recalculate coverage
	totalCoveredLines := 0
	totalExecutableLines := 0
	for _, file := range targetData.Files {
		totalCoveredLines += file.CoveredLines
		totalExecutableLines += file.ExecutableLines
	}

	if totalExecutableLines > 0 {
		targetData.LineCoverage = float64(totalCoveredLines) / float64(totalExecutableLines)
	}

	// Output results
	coverage = targetData.LineCoverage * 100
	fmt.Printf("Updated Coverage for Target '%s':\n", target)
	fmt.Printf("  Covered Lines: %d\n", totalCoveredLines)
	fmt.Printf("  Executable Lines: %d\n", totalExecutableLines)
	fmt.Printf("  Line Coverage: %.2f%%\n", coverage)
	return coverage, nil
}

// GenerateCoverageFile generates a coverage report from an Xcode result bundle and writes it to the specified output file.
// Takes the path to the xcresult bundle and the desired output file path as arguments. Returns an error if the process fails.
func GenerateCoverageFile(xcresultPath, outputPath string) error {
	fmt.Println("Generating coverage report...")
	cmd := exec.Command("xcrun", "xccov", "view", "--report", "--json", xcresultPath)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create coverage file: %w", err)
	}
	defer outputFile.Close()

	cmd.Stdout = outputFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run xcrun xccov: %w", err)
	}

	fmt.Printf("Coverage report generated at %s\n", outputPath)
	return nil
}
