package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// versionCmd is a Cobra command that outputs the version number of the xccov-proc CLI tool.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of xccov-proc",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Executing version command")
	},
}

// init adds the versionCmd to the rootCmd, ensuring the version information is included in the command-line interface.
func init() {
	rootCmd.AddCommand(versionCmd)
}
