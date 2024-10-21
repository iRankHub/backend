package script

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// FormatGoCode formats all Go files in the project directory and its subdirectories
func FormatGoCode() error {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Navigate to the project root (two levels up from cmd/server)
	projectRoot := filepath.Join(wd, "..", "..")

	// Run gofmt on all Go files in the project root and subdirectories
	cmd := exec.Command("gofmt", "-l", "-w", ".")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to format Go code: %w\n%s", err, output)
	}

	// If there's any output, it means files were modified
	if len(output) > 0 {
		return nil // Formatting successful, changes were made
	}

	// No output means no changes were necessary
	return nil
}
