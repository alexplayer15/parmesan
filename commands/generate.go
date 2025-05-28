package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	oas_struct "github.com/alexplayer15/parmesan/data"
	"github.com/alexplayer15/parmesan/request_generator"
	"github.com/spf13/cobra"
)

func newGenerateRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-request",
		Short: "Generate a HTTP request from an OpenAPI Spec",
		Args:  cobra.ExactArgs(1),
	}

	flags := bindFlags(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		oasFile := args[0]

		oas, err := handleOAS(oasFile)
		if err != nil {
			return err
		}

		outputDir := flags.OutputDir

		if err := validateOutputPath(outputDir); err != nil {
			return err
		}
		if err := ensureDirectory(outputDir); err != nil {
			return err
		}

		outputFile := filepath.Join(outputDir, changeExtension(oasFile, ".http"))

		chosenServerIndex := flags.WithServer

		if err := validateChosenServerUrl(chosenServerIndex, oas); err != nil {
			return err
		}

		httpRequest, err := request_generator.GenerateHttpRequest(oas, chosenServerIndex)
		if err != nil {
			return fmt.Errorf("failed to generate HTTP request: %w", err)
		}

		if err := os.WriteFile(outputFile, []byte(httpRequest), 0644); err != nil {
			return fmt.Errorf("failed to write HTTP file: %w", err)
		}

		return nil
	}
	return cmd
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

func validateChosenServerUrl(chosenServerUrl int, oas oas_struct.OAS) error {
	if chosenServerUrl < 0 || chosenServerUrl >= len(oas.Servers) {
		return fmt.Errorf("invalid server index %d: There are %d servers available starting from 0", chosenServerUrl, len(oas.Servers))
	}

	return nil
}

func changeExtension(filePath, newExt string) string {
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return base + newExt
}
