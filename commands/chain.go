package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alexplayer15/parmesan/chain_logic"
	flag_helpers "github.com/alexplayer15/parmesan/commands/flag_hepers"
	"github.com/alexplayer15/parmesan/data"
	hooks_logic "github.com/alexplayer15/parmesan/hooks"
	"github.com/alexplayer15/parmesan/request_generator"
	"github.com/alexplayer15/parmesan/request_sender"
	"github.com/spf13/cobra"
)

type RulesFile []ChainRule

type ChainRule struct {
	Name   string            `yaml:"name,omitempty"`
	Method string            `yaml:"method"`
	Path   string            `yaml:"path"`
	Pick   map[string]string `yaml:"pick,omitempty"`
	Needs  map[string]string `yaml:"needs,omitempty"`
}

var extractedValues []string

func newChainRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-request",
		Short: "Send an ordered chain of HTTP request from an OpenAPI Spec",
		Args:  cobra.ExactArgs(1),
	}

	flags := bindFlags(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		oasFile := args[0]

		oas, err := handleOAS(oasFile)
		if err != nil {
			return err
		}

		chosenServerIndex := flags.WithServer

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

		rulesFile := flags.RulesFile

		rules, err := chain_logic.UnmarshalRulesFile(rulesFile)
		if err != nil {
			return err
		}

		orderedRequests, err := chain_logic.OrderRequests(requests, rules)

		var allResponses []data.SavedResponse
		hooks := flags.HooksFile

		var hooksFile hooks_logic.HooksFile

		if hooks != "" {
			hooksFile, err = hooks_logic.UnmarshalHooksFile(hooks)
			if err != nil {
				return err
			}
		}

		for _, req := range orderedRequests {
			if hooks != "" {
				matchingHook := hooks_logic.TryAndFindHookForThisRequest(hooksFile, req)

				if !matchingHook.IsEmpty() {
					req.Body, err = hooks_logic.ModifyRequestBodyUsingHook(matchingHook, req.Body)
					if err != nil {
						return err
					}
				}
			}

			// if i != 0 {
			// 	chain_logic.ApplyInjectionRules(req, rules, extractedValues)
			// }

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

			savedResp := data.SavedResponse{
				Method:   req.Method,
				Url:      req.Url,
				Status:   statusCode,
				Response: parsedBody,
			}

			// extractedValues, err = chain_logic.ApplyExtractionRules(savedResp.Response, rules, req)
			// if err != nil {
			// 	return err
			// }

			allResponses = append(allResponses, savedResp)
		}

		outputDir := flags.OutputDir

		if err := flag_helpers.ValidateOutput(outputDir); err != nil {
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
	}
	return cmd
}
