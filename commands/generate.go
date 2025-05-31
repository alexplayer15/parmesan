package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag_helpers "github.com/alexplayer15/parmesan/commands/flag_hepers"
	oas_struct "github.com/alexplayer15/parmesan/data"
	"github.com/alexplayer15/parmesan/request_generator"
	"github.com/spf13/cobra"
)

func newGenerateRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-request",
		Short: "Generate a HTTP request from an OpenAPI Spec",
		Args:  cobra.ExactArgs(1),
	}

	flags := bindGenerateRequestFlags(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		oasFile := args[0]

		oas, err := handleOAS(oasFile)
		if err != nil {
			return err
		}

		outputDir := flags.OutputDir

		if err := flag_helpers.ValidateOutput(outputDir); err != nil {
			return err
		}

		outputFile := filepath.Join(outputDir, changeExtension(oasFile, ".http"))

		chosenServerIndex := flags.WithServer

		if err := validateChosenServerUrl(chosenServerIndex, oas); err != nil {
			return err
		}

		httpRequest, err := request_generator.GenerateHttpRequest(oas, chosenServerIndex)
		if err != nil {
			return fmt.Errorf("failed to generate HTTP request: %w", err)
		}

		if err := os.WriteFile(outputFile, []byte(httpRequest), 0644); err != nil {
			return fmt.Errorf("failed to write HTTP file: %w", err)
		}

		return nil
	}
	return cmd
}

func validateChosenServerUrl(chosenServerUrl int, oas oas_struct.OAS) error {
	if chosenServerUrl < 0 || chosenServerUrl >= len(oas.Servers) {
		return fmt.Errorf("invalid server index %d: There are %d servers available starting from 0", chosenServerUrl, len(oas.Servers))
	}

	return nil
}

func changeExtension(filePath, newExt string) string {
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return base + newExt
}
