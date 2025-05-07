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

var GenerateRequestCmd = &cobra.Command{
	Use:   "generate-request",
	Short: "Generate a .http request from an OpenAPI file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		oasFile := args[0]
		if err := verifyIfFileExists(oasFile); err != nil {
			return err
		}
		oas, err := parseOASFile(oasFile)
		if err != nil {
			return fmt.Errorf("error reading OAS file: %w", err)
		}
		if err := checkIfOASFileIsValid(oas); err != nil {
			return fmt.Errorf("invalid OAS structure: %w", err)
		}

		httpRequest, err := request_generator.GenerateHttpRequest(oas)
		if err != nil {
			return fmt.Errorf("failed to generate HTTP request: %w", err)
		}

		outputFile := changeExtension(oasFile, ".http")

		if err := os.WriteFile(outputFile, []byte(httpRequest), 0644); err != nil {
			return fmt.Errorf("failed to write HTTP file: %w", err)
		}

		return nil
	},
}

func verifyIfFileExists(file string) error {
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

func changeExtension(filePath, newExt string) string {
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return base + newExt
}
