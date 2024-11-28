package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Root command: the base command for the CLI
var rootCmd = &cobra.Command{
	Use:   "xccov-proc",
	Short: "A tool to process xcode coverage files with customization options",
	Long: `xccov-proc is a CLI tool to process and filter Xcode coverage reports.
	It allows for excluding files, focusing on specific targets, and recalculating
	coverage metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("executing => root command.")
	},
}

// Execute starts the CLI application
func Execute() {
	// Check for Bitrise step context
	log.Println("executing => Execute [root]")
	log.Printf("BITRISE_STEP_ID = %s", os.Getenv("BITRISE_STEP_ID"))
	if os.Getenv("BITRISE_STEP_ID") != "" {
		executeBitriseMode()
	} else {
		log.Println("using else branch")
		if err := rootCmd.Execute(); err != nil {
			os.Exit(1)
		}
	}
}

func init() {
	log.Println("executing => init [root]")
}

func initConfig() {
	log.Println("executing => initConfig")
}

func executeBitriseMode() {
	log.Println("running in Bitrise step mode [root]")
	xcresultPath := os.Getenv("xcresult_path")
	target := os.Getenv("target")
	excludeFiles := os.Getenv("exclude_files")
	excludeConfigPath := os.Getenv("exclude_config_path")
	outputCoveragePath := "./recalculated-coverage.json"

	// Convert space delimited excludeFiles to a slice
	var excludeFileList []string
	if excludeFiles != "" {
		excludeFileList = strings.Split(excludeFiles, " ")
	}
	log.Printf("exclude: %v", excludeFileList)

	// Validate required inputs
	if xcresultPath == "" || target == "" {
		log.Fatalf("Missing required Bitrise inputs: xcresult_path and target")
	}
	log.Printf("xcresult: %s, target: %s", xcresultPath, target)

	// Extract cov from xcresult
	args := []string{
		"extract",
		"--xcresult-path", xcresultPath,
		"--coverage-output", outputCoveragePath,
	}

	// Override os.Args to simulate CLI invocation
	os.Args = append([]string{os.Args[0]}, args...)

	// Execute the root command with the new arguments
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute rootCmd with process: %v", err)
	}
	log.Println("Done extracting...")

	// Construct command-line arguments
	args = []string{
		"process", // Simulates the "process" subcommand
		"--xcresult-path", xcresultPath,
		"--target", target,
		"--output", outputCoveragePath,
	}

	if excludeConfigPath != "" {
		args = append(args, "--exclude-config", excludeConfigPath)
	}
	for _, exclude := range excludeFileList {
		args = append(args, "--exclude", exclude)
	}

	// Override os.Args to simulate CLI invocation
	os.Args = append([]string{os.Args[0]}, args...)

	// Execute the root command with the new arguments
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute rootCmd with process: %v", err)
	}
	log.Println("Done calculating...")
}
