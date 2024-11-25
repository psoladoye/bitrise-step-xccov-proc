package cmd

import (
	"log"
	"os"

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
	log.Println("executing => Execute")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	log.Println("executing => init")
}

func initConfig() {
	log.Println("executing => initConfig")
}
