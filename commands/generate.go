package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
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
		if err := validateOasArgument(file); err != nil {
			return err
		}
		fmt.Printf("Generating request from: %s\n", file)
		fmt.Printf("Output format: %s\n", output)
		return nil
	},
}

func validateOasArgument(file string) error {
	fileExistsErr := verifyIfFileExists(file)
	if fileExistsErr != nil {
		return fileExistsErr
	}
	oasContentErr := verifyOasContent(file)
	if oasContentErr != nil {
		return oasContentErr
	}
	return nil
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

func verifyOasContent(file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ext := filepath.Ext(file)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	var data map[string]any

	switch ext {
	case "json":
		if err := json.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("invalid YAML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	requiredFields := []string{"openapi", "info", "paths"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			return fmt.Errorf("missing required OAS field: %s", field)
		}
	}

	return nil
}

func init() {
	GenerateRequestCmd.Flags().StringVar(&output, "output", "http", "Output format (e.g., http)")
	RootCmd.AddCommand(GenerateRequestCmd)
}
