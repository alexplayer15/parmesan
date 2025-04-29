package command_tests

import (
	"testing"

	"github.com/alexp/parmesan/src/commands"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRequestCmd(t *testing.T) {
	//Arrange
	rootCmd := &cobra.Command{
		Use: "parmesan",
	}
	rootCmd.AddCommand(commands.GenerateRequestCmd)

	//Assert
	cmd, _, err := rootCmd.Find([]string{"generate-request"})
	assert.NoError(t, err, "Error should be nil when finding the 'generate-request' command")
	assert.NotNil(t, cmd, "'generate-request' command should be registered")
	assert.Equal(t, "generate-request", cmd.Use, "Command name should be 'generate-request'")
}

func Test_WhenGenerateRequestIsNotGivenAnArg_ShouldFail(t *testing.T){
	//Arrange
	rootCmd := &cobra.Command{
		Use: "parmesan",
	}
	rootCmd.AddCommand(commands.GenerateRequestCmd)
	rootCmd.SetArgs([]string{"generate-request"})

	//Act
	err := rootCmd.Execute()

	//Assert
    assert.Error(t, err, "Command should fail if no argument is given")
}

func Test_WhenGenerateRequestIsGivenMoreThanOneArg_ShouldFail(t *testing.T){
	//Arrange
	rootCmd := &cobra.Command{
		Use: "parmesan",
	}
	rootCmd.AddCommand(commands.GenerateRequestCmd)
	rootCmd.SetArgs([]string{"generate-request", "file1.yaml", "file2.yaml"})

	//Act
	err := rootCmd.Execute()

	//Assert
    assert.Error(t, err, "Command should fail if no argument is given")
}

func Test_WhenGenerateRequestIsInitialised_ShouldHaveAnOutputFlag(t *testing.T){
	//Arrange
	rootCmd := &cobra.Command{
		Use: "parmesan",
	}
	rootCmd.AddCommand(commands.GenerateRequestCmd)

	//Act
	cmd, _, err := rootCmd.Find([]string{"generate-request"})

	//Assert
	assert.NoError(t, err, "Command should be found")
	flag := cmd.Flags().Lookup("output")
	assert.NotNil(t, flag, "The 'output' flag should be registered")
}