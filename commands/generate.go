package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	fileExtensionErr := verifyFileExtension(file)
	if fileExtensionErr != nil {
		return fileExtensionErr
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

func verifyFileExtension(file string) error {
	allowedExtensions := []string{".yml", ".yaml", ".json"}

	ext := strings.ToLower(filepath.Ext(file))

	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return nil
		}
	}

	return fmt.Errorf("OAS must be a JSON or YAML file")
}

func init() {
	GenerateRequestCmd.Flags().StringVar(&output, "output", "http", "Output format (e.g., http)")
	RootCmd.AddCommand(GenerateRequestCmd)
}
