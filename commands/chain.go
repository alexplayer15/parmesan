package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	hooks_logic "github.com/alexplayer15/parmesan/hooks"
	"github.com/alexplayer15/parmesan/request_generator"
	"github.com/alexplayer15/parmesan/request_sender"
	"github.com/spf13/cobra"
)

type SavedResponse struct {
	Method   string `json:"method"`
	Url      string `json:"url"`
	Status   int    `json:"status"`
	Response any    `json:"response"`
}

type RulesFile []ChainRule

type ChainRule struct {
	Name   string            `yaml:"name,omitempty"`
	Method string            `yaml:"method"`
	Path   string            `yaml:"path"`
	Pick   map[string]string `yaml:"pick,omitempty"`
	Needs  map[string]string `yaml:"needs,omitempty"`
}

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

			type ExtractedValues map[string]map[string]any

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
				if hooks != "" {
					matchingHook := hooks_logic.TryAndFindHookForThisRequest(hooksFile, req)

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

	cmd.Flags().String("hooks", "", "Location of hooks file to modify request values.")
	cmd.Flags().String("rules", "", "Rules your chain of requests should follow.")

	return cmd
}
