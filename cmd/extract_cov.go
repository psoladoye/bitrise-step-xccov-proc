package cmd

import (
	"com.github/psoladoye/bitrise-step-xccov-proc/utils"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// xcresultPath is a string variable that holds the path to the .xcresult file used for extracting coverage data.
var xcresultPath string

// outputCoveragePath defines the path where the generated coverage.json file will be saved.
var outputCoveragePath string

// extractCovCmd is a Cobra command that extracts coverage data from an .xcresult file and outputs it as a coverage.json file.
var extractCovCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract coverage data from an Xcode result bundle",
	Long: `Extracts coverage information from a given .xcresult file
and generates a coverage.json file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("extracting coverage report...")

		if xcresultPath == "" || outputCoveragePath == "" {
			log.Println("error: both --xcresult-path and --output are required.")
			cmd.Usage()
			os.Exit(1)
		}

		if err := utils.GenerateCoverageFile(xcresultPath, outputCoveragePath); err != nil {
			log.Printf("error generating coverage resport: %v/n", err)
			os.Exit(1)
		}
		fmt.Printf("coverage report generated succesfully at: %s\n", outputCoveragePath)
	},
}

// init sets up the extractCovCmd by adding it to rootCmd and defines required flags for extracting coverage data.
func init() {
	rootCmd.AddCommand(extractCovCmd)

	// Add flags specific to the extract command
	extractCovCmd.Flags().StringVar(&xcresultPath, "xcresult-path", "", "Path to the .xcresult file (required)")
	extractCovCmd.Flags().StringVar(&outputCoveragePath, "coverage-output", "", "Path to save the generated coverage.json file (required)")
}
