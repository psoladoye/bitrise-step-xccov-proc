package utils

import (
	"fmt"
	"os/exec"
)

// ExportCoverageAsEnv exports a coverage float value as an environment variable using "envman".
// It formats the value to two decimal places and assigns it to the key "XCCOV_PROC_COVERAGE_VALUE".
// Returns an error if the command execution fails.
func ExportCoverageAsEnv(value float64) error {
	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "XCCOV_PROC_COVERAGE_VALUE", "--value",
		fmt.Sprintf("%.2f", value)).CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		return err
	}
	return nil
}
