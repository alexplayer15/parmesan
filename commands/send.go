package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag_helpers "github.com/alexplayer15/parmesan/commands/flag_hepers"
	"github.com/alexplayer15/parmesan/data"
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

		method := flags.Method
		paths := flags.Paths

		if err := request_sender.ValidateHTTPMethod(method); err != nil {
			return err
		}
		for _, path := range paths {
			if err := validatePathInput(path); err != nil {
				return err //to do: use typed validation error in the above method and return it here
			}
		}

		var allResponses []data.SavedResponse
		hooks := flags.HooksFile

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

			responseBody, statusCode, headers, err := request_sender.SendHTTPRequest(req)
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
				Headers:  headers,
			}

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
