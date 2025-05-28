package commands

import (
	"fmt"

	"github.com/alexplayer15/parmesan/request_generator"
	"github.com/alexplayer15/parmesan/request_sender"
	"github.com/spf13/cobra"
)

func newChainRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-request",
		Short: "Send an ordered chain of HTTP request from an OpenAPI Spec",
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

			chosenServerIndex, _ := cmd.Flags().GetInt("with-server")

			if err := validateChosenServerUrl(chosenServerIndex, oas); err != nil {
				return err
			}

			httpRequestFile, err := request_generator.GenerateHttpRequest(oas, chosenServerIndex)
			if err != nil {
				return fmt.Errorf("failed to generate HTTP request: %w", err)
			}

			requests, err := request_sender.ParseHttpRequestFile(httpRequestFile)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("hooks", "", "Location of hooks file to modify request values.")
	cmd.Flags().String("rules", "", "Rules your chain of requests should follow.")

	return cmd
}
