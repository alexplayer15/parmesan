package commands

import (
	"fmt"
	"os"

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
		if err := verifyIfFileExists(file); err != nil {
			return err
		}
		fmt.Printf("Generating request from: %s\n", file)
		fmt.Printf("Output format: %s\n", output)
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

func init() {
	GenerateRequestCmd.Flags().StringVar(&output, "output", "http", "Output format (e.g., http)")
    RootCmd.AddCommand(GenerateRequestCmd)
}
