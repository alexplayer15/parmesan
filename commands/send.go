package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/alexplayer15/parmesan/request_generator"
	"github.com/alexplayer15/parmesan/request_sender"
	"github.com/spf13/cobra"
)

func newSendRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-request",
		Short: "Send a HTTP request from an OpenAPI Spec",
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

			method, _ := cmd.Flags().GetString("method")
			paths, _ := cmd.Flags().GetStringSlice("path")

			//add in validation for method and path input

			for _, req := range requests {

				if method != "*" && req.Method != method {
					continue
				}

				if !urlMatchesPaths(req.Url, paths) {
					continue
				}

				response, err := request_sender.SendHTTPRequest(req)
				if err != nil {
					log.Printf("Failed to send request %s %s: %v", req.Method, req.Url, err)
					continue
				}

				// Print the response for now
				fmt.Printf("Response for %s %s: %s\n", req.Method, req.Url, response)
			}

			return nil
		},
	}

	// Define flags
	cmd.Flags().Int("with-server", 0, "Which server url to use from OAS. 0 = First URL.")
	cmd.Flags().String("method", "*", "Choose with requests you want to send from your OAS by method. Default is all methods.")
	cmd.Flags().StringSlice("path", []string{}, "Choose with requests you want to send from your OAS by path. Default is all paths.")

	return cmd
}

func urlMatchesPaths(url string, paths []string) bool {
	if len(paths) == 0 {
		return true
	}
	for _, path := range paths {
		if strings.Contains(url, path) {
			return true
		}
	}
	return false
}
