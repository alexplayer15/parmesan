package commands

import (
	"github.com/spf13/cobra"
)

var output string

var RootCmd = &cobra.Command{
	Use:   "parmesan",
	Short: "CLI tool to generate requests based off your OAS",
}

func init() {
	GenerateRequestCmd.Flags().StringVar(&output, "output", "http", "Directory of output")
	RootCmd.AddCommand(GenerateRequestCmd)
}
