package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
	"github.com/alexplayer15/parmesan/data"
	"github.com/alexplayer15/parmesan/request_generator"
)

var output string

var RootCmd = &cobra.Command{
	Use:   "parmesan",
	Short: "CLI tool to generate requests based off your OAS",
}

var GenerateRequestCmd = &cobra.Command{
	Use:   "generate-request",
	Short: "Generate a .http request from an OpenAPI file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		if err := verifyIfFileExists(file); err != nil {
			return err
		}
		oas, err := ReadOASFile(args[0])
		if err != nil {
			fmt.Println("Error reading OAS file:", err)
			os.Exit(1)
		}

		httpRequest, err := request_generator.GenerateHttpRequest(oas, "/hello")
		filename := filepath.Join(".", "request.http")
		os.WriteFile(filename, []byte(httpRequest), 0644)
		if err != nil {
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

func ReadOASFile(file string) (oas_struct.OAS, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return oas_struct.OAS{}, fmt.Errorf("failed to read file: %w", err)
	}

	ext := filepath.Ext(file)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

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

	if oas.OpenAPI == "" {
		return oas_struct.OAS{}, fmt.Errorf("missing required OAS field: openapi")
	}
	if oas.Info.Title == "" {
		return oas_struct.OAS{}, fmt.Errorf("missing required OAS field: info")
	}
	if len(oas.Paths) == 0 {
		return oas_struct.OAS{}, fmt.Errorf("missing required OAS field: paths")
	}

	return oas, nil
}

func init() {
	GenerateRequestCmd.Flags().StringVar(&output, "output", "http", "Output format (e.g., http)")
	RootCmd.AddCommand(GenerateRequestCmd)
}
