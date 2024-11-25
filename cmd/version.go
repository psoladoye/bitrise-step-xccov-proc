package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of xccov-proc",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Executing version command")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
