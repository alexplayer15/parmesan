package command_tests

import (
	"testing"

	"github.com/alexp/parmesan/src/commands"
	"github.com/stretchr/testify/assert"
)

func Test_WhenGenerateRequestIsNotGivenAnArg_ShouldFail(t *testing.T) {
	//Arrange
	commands.RootCmd.SetArgs([]string{"generate-request"})

	//Act
	err := commands.RootCmd.Execute()

	//Assert
	assert.Error(t, err, "Command should fail if no argument is given")
}

func Test_WhenGenerateRequestIsGivenMoreThanOneArg_ShouldFail(t *testing.T) {
	//Arrange
	commands.RootCmd.SetArgs([]string{"generate-request", "file1.yaml", "file2.yaml"})

	//Act
	err := commands.RootCmd.Execute()

	//Assert
	assert.Error(t, err, "Command should fail if no argument is given")
}
func Test_WhenGenerateRequestIsGivenValidArguments_ShouldNotError(t *testing.T) {
	//Arrange 
	commands.RootCmd.SetArgs([]string{"generate-request", "file1.yaml"})

	//Act 
	err := commands.RootCmd.Execute()

	//Assert
	assert.NoError(t, err, "Should be no error when 1 arg and correct command is entered")
}
