package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	oas_struct "github.com/alexplayer15/parmesan/data"
	"github.com/alexplayer15/parmesan/request_generator"
	"github.com/spf13/cobra"
)

func newGenerateRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-request",
		Short: "Generate a HTTP request from an OpenAPI Spec",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			oasFile := args[0]
			if err := checkIfFileExists(oasFile); err != nil {
				return err
			}
			oas, err := parseOASFile(oasFile)
			if err != nil {
				return fmt.Errorf("error reading OAS file: %w", err)
			}
			if err := checkIfOASFileIsValid(oas); err != nil {
				return fmt.Errorf("invalid OAS structure: %w", err)
			}

			outputDir, _ := cmd.Flags().GetString("output")

			if err := validateOutputPath(outputDir); err != nil {
				return err
			}
			if err := ensureDirectory(outputDir); err != nil {
				return err
			}

			outputFile := filepath.Join(outputDir, changeExtension(oasFile, ".http"))

			chosenServerIndex, _ := cmd.Flags().GetInt("with-server")

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
		},
	}

	// Define flags
	cmd.Flags().String("output", ".", "Directory of output")
	cmd.Flags().Int("with-server", 0, "Which server url to use from OAS. 0 = First URL.")

	return cmd
}

func checkIfFileExists(file string) error {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}
	if err != nil {
		return fmt.Errorf("error checking file: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("provided 'file' is a directory")
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

func parseOASFile(file string) (oas_struct.OAS, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return oas_struct.OAS{}, fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.TrimPrefix(filepath.Ext(file), ".")

	var oas oas_struct.OAS

	switch ext {
	case "json":
		if err := json.Unmarshal(content, &oas); err != nil {
			return oas_struct.OAS{}, fmt.Errorf("invalid JSON: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal(content, &oas); err != nil {
			return oas_struct.OAS{}, fmt.Errorf("invalid YAML: %w", err)
		}
	default:
		return oas_struct.OAS{}, fmt.Errorf("unsupported file extension: %s", ext)
	}

	return oas, nil
}

func checkIfOASFileIsValid(oas oas_struct.OAS) error {
	if oas.OpenAPI == "" {
		return fmt.Errorf("missing required OAS field: openapi")
	}
	if oas.Info.Title == "" {
		return fmt.Errorf("missing required OAS field: info")
	}
	if len(oas.Servers) == 0 {
		return fmt.Errorf("no server URL found in OAS")
	}
	if len(oas.Paths) == 0 {
		return fmt.Errorf("missing required OAS field: paths")
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
