package cmd

import (
	"com.github/psoladoye/bitrise-step-xccov-proc/utils"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var xcresultPath string
var outputCoveragePath string

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

func init() {
	rootCmd.AddCommand(extractCovCmd)

	// Add flags specific to the extract command
	extractCovCmd.Flags().StringVar(&xcresultPath, "xcresult-path", "", "Path to the .xcresult file (required)")
	extractCovCmd.Flags().StringVar(&outputCoveragePath, "coverage-output", "", "Path to save the generated coverage.json file (required)")
}
