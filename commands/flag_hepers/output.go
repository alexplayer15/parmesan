package flag_helpers

import (
	"fmt"
	"os"
)

func ValidateOutput(outputDir string) error {

	if err := validateOutputPath(outputDir); err != nil {
		return err
	}
	if err := ensureDirectory(outputDir); err != nil {
		return err
	}
	return nil
}

func validateOutputPath(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {

		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to check output path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("invalid output directory: path exists and is a file")
	}
	return nil
}

func ensureDirectory(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	return nil
}
