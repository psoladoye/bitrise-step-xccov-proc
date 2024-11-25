package cmd

import (
	"com.github/psoladoye/bitrise-step-xccov-proc/utils"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

var target string
var excludeFiles []string
var excludeConfigPath string

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process and filter an Xcode coverage report",
	Long: `Processes an Xcode coverage file, excludes specified files,
and recalculates coverage metrics for the specified target.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("executing => process command")

		// Validate flags
		if xcresultPath == "" || target == "" {
			log.Println("error: --xcresult-path and --target are required.")
			cmd.Usage()
			os.Exit(1)
		}

		// Parse exclusions from YAML config
		excludePaths := make(map[string]struct{})
		if excludeConfigPath != "" {
			log.Printf("Loading exclusion config from: %s\n", excludeConfigPath)
			configExclusions, err := utils.ParseExclusionConfig(excludeConfigPath)
			if err != nil {
				log.Fatalf("Failed to load exclusion config: %v\n", err)
			}
			for path := range configExclusions {
				excludePaths[path] = struct{}{}
			}
		}

		// Add exclusions from command-line arguments
		for _, excludeFile := range excludeFiles {
			excludePaths[excludeFile] = struct{}{}
		}

		// Generate coverage.json file
		outputPath := outputCoveragePath
		if outputPath == "" {
			outputPath = "./coverage.json" // Default output path
		}

		log.Printf("Extracting coverage from: %s\n", xcresultPath)
		if err := utils.GenerateCoverageFile(xcresultPath, outputPath); err != nil {
			log.Fatalf("Failed to generate coverage file: %v\n", err)
		}

		// Process coverage file
		log.Printf("Processing coverage file for target: %s\n", target)
		coverage, err := utils.ProcessCoverage(outputPath, target, excludePaths)
		if err != nil {
			log.Fatalf("Failed to process coverage file: %v\n", err)
		}

		cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "CODE_COVERAGE", "--value", fmt.Sprintf("%f", coverage)).CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
			os.Exit(1)
		}

		fmt.Printf("Coverage processing completed successfully. Processed file: %s\n", outputPath)
	},
}

func init() {
	// Add the process command to the root command
	rootCmd.AddCommand(processCmd)

	// Define flags for the process command
	processCmd.Flags().StringVar(&target, "target", "", "Target name to process (required)")
	processCmd.Flags().StringSliceVar(&excludeFiles, "exclude", nil, "Files to exclude from coverage")
	processCmd.Flags().StringVar(&excludeConfigPath, "exclude-config", "", "Path to YAML config for exclusions")
	processCmd.Flags().StringVar(&xcresultPath, "xcresult-path", "", "Path to the .xcresult file (required)")
	processCmd.Flags().StringVar(&outputCoveragePath, "output", "", "Path to save the generated coverage.json file (optional)")
}
