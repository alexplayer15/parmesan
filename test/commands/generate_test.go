package command_tests

import (
	"testing"

	"github.com/alexp/parmesan/src/commands"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRequestCmd(t *testing.T) {
	rootCmd := &cobra.Command{
		Use: "parmesan",
	}

	rootCmd.AddCommand(commands.GenerateRequestCmd)

	cmd, _, err := rootCmd.Find([]string{"generate-request"})
	assert.NoError(t, err, "Error should be nil when finding the 'generate-request' command")
	assert.NotNil(t, cmd, "'generate-request' command should be registered")
	assert.Equal(t, "generate-request", cmd.Use, "Command name should be 'generate-request'")
}
