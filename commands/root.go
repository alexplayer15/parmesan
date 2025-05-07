package commands

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "parmesan",
		Short: "CLI tool to generate requests based off your OAS",
	}

	rootCmd.AddCommand(newGenerateRequestCmd())

	return rootCmd
}

var RootCmd = NewRootCmd()
