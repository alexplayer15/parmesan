package commands

import (
    "fmt"

    "github.com/spf13/cobra"
)

var output string

var GenerateRequestCmd = &cobra.Command{
    Use:   "generate-request",
    Short: "Generate a .http request from an OpenAPI file",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        file := args[0]
        fmt.Printf("Generating request from: %s\n", file)
        fmt.Printf("Output format: %s\n", output)
    },
}

func init() {
    GenerateRequestCmd.Flags().StringVar(&output, "output", "http", "Output format (e.g., http)")
}
