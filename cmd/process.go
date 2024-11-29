package cmd

import (
	"com.github/psoladoye/bitrise-step-xccov-proc/utils"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// target specifies the name of the target to process for coverage metrics recalculation.
var target string

// excludeFiles is a slice of strings representing file paths to be excluded from coverage processing.
var excludeFiles []string

// excludeConfigPath specifies the path to a YAML configuration file containing paths that should be excluded from processing.
var excludeConfigPath string

// processCmd is a Cobra command used to process and filter an Xcode coverage report, recalculating coverage metrics.
// The command requires the --xcresult-path and --target flags to specify the path of the Xcode coverage file and the target.
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

		err = utils.ExportCoverageAsEnv(coverage)
		if err != nil {
			log.Fatalf("Failed to export coverage as env vars: %v\n", err)
		}

		fmt.Printf("Coverage processing completed successfully. Processed file: %s\n", outputPath)
	},
}

// init initializes the command-line interface by adding the 'process' command to the root command.
// It sets up flags for the process command allowing users to specify target, exclude files,
// exclude configuration path, path to the .xcresult file, and an optional output path for the coverage file.
func init() {
	log.Println("executing ==>  init [process]")
	// Add the process command to the root command
	rootCmd.AddCommand(processCmd)

	// Define flags for the process command
	processCmd.Flags().StringVar(&target, "target", "", "Target name to process (required)")
	processCmd.Flags().StringSliceVar(&excludeFiles, "exclude", nil, "Files to exclude from coverage")
	processCmd.Flags().StringVar(&excludeConfigPath, "exclude-config", "", "Path to YAML config for exclusions")
	processCmd.Flags().StringVar(&xcresultPath, "xcresult-path", "", "Path to the .xcresult file (required)")
	processCmd.Flags().StringVar(&outputCoveragePath, "output", "", "Path to save the generated coverage.json file (optional)")
}
