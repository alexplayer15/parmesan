package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	hooks_logic "github.com/alexplayer15/parmesan/hooks"
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

			if err := request_sender.ValidateHTTPMethod(method); err != nil {
				return err
			}
			for _, path := range paths {
				if err := validatePathInput(path); err != nil {
					return err //to do: use typed validation error in the above method and return it here
				}
			}

			var allResponses []SavedResponse
			hooks, _ := cmd.Flags().GetString("hooks")

			var hooksFile hooks_logic.HooksFile

			if hooks != "" {
				hooksFile, err = hooks_logic.UnmarshalHooksFile(hooks)
				if err != nil {
					return err
				}
			}

			for _, req := range requests {
				if method != "*" && req.Method != method {
					continue
				}

				if !urlMatchesPaths(req.Url, paths) {
					continue
				}

				if hooks != "" {
					matchingHook := hooks_logic.TryAndFindHookForThisRequest(hooksFile, req)

					//to do: think of a way to reduce this nesting
					if !matchingHook.IsEmpty() {
						req.Body, err = hooks_logic.ModifyRequestBodyUsingHook(matchingHook, req.Body)
						if err != nil {
							return err
						}
					}
				}

				responseBody, statusCode, err := request_sender.SendHTTPRequest(req)
				if err != nil {
					log.Printf("Failed to send request %s %s: %v", req.Method, req.Url, err)
					continue
				}

				var parsedBody any
				if err := json.Unmarshal([]byte(responseBody), &parsedBody); err != nil {
					log.Printf("Failed to parse JSON body for %s %s: %v. Saving as string.", req.Method, req.Url, err)
					parsedBody = responseBody
				}

				savedResp := SavedResponse{
					Method:   req.Method,
					Url:      req.Url,
					Status:   statusCode,
					Response: parsedBody,
				}

				allResponses = append(allResponses, savedResp)
			}

			outputDir, _ := cmd.Flags().GetString("output")

			if err := validateOutputPath(outputDir); err != nil {
				return err
			}
			if err := ensureDirectory(outputDir); err != nil {
				return err
			}

			filePath := filepath.Join(outputDir, changeExtension(oasFile, ".json"))

			jsonData, err := json.MarshalIndent(allResponses, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal all responses: %w", err)
			}

			if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
				return fmt.Errorf("failed to write final output file: %w", err)
			}

			fmt.Printf("Saved all responses to %s\n", filePath)

			return nil
		},
	}

	cmd.Flags().Int("with-server", 0, "Which server url to use from OAS. 0 = First URL.")
	cmd.Flags().String("method", "*", "Choose with requests you want to send from your OAS by method. Default is all methods.")
	cmd.Flags().StringSlice("path", []string{}, "Choose with requests you want to send from your OAS by path. Default is all paths.")
	cmd.Flags().String("output", ".", "Directory of output for HTTP responses.")
	cmd.Flags().String("hooks", "", "Location of hooks file to modify request values.")

	return cmd
}

func urlMatchesPaths(url string, paths []string) bool {
	if len(paths) == 0 {
		return true
	}
	for _, path := range paths {
		validatePathInput(path)
		if strings.Contains(url, path) {
			return true
		}
	}
	return false
}

func validatePathInput(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path input cannot be empty or only spaces")
	}

	if strings.Contains(path, " ") {
		return fmt.Errorf("path input cannot contain spaces")
	}

	if strings.Contains(path, "http://") || strings.Contains(path, "https://") {
		return fmt.Errorf("path input must not contain full URLs")
	}

	if strings.ContainsAny(path, "\\?#") {
		return fmt.Errorf("path input contains illegal characters (\\, ?, #)")
	}

	return nil
}
