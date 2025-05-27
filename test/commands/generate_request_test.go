package command_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/commands"
	"github.com/alexplayer15/parmesan/test_helpers"

	"github.com/stretchr/testify/assert"
)

func Test_WhenGenerateRequestIsNotGivenAnArg_ShouldFail(t *testing.T) {
	//Arrange
	cmd := commands.NewRootCmd()
	cmd.SetArgs([]string{
		"generate-request",
	})

	// Act
	err := cmd.Execute()

	//Assert
	assert.Error(t, err, "Command should fail if no argument is given")
}

func Test_WhenGenerateRequestIsGivenMoreThanOneArg_ShouldFail(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupGenRequestTest(t, "oas.yml", "../testOas.yml", "dummyArg")

	// Act
	err := cmd.Execute()

	//Assert
	assert.Error(t, err, "Command should fail if no argument is given")
}

func Test_WhenOASHasNoServerURL_ShouldReturnError(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupGenRequestTest(t, "oas.yml", "../testOasNoServerUrl.yml")

	// Act
	err := cmd.Execute()

	//Assert
	assert.Error(t, err, "Should be an error alerting user there is no server URL in the OAS")
}
func Test_WhenOASDoesNotExist_ShouldReturnError(t *testing.T) {
	//Arrange
	cmd := commands.NewRootCmd()
	cmd.SetArgs([]string{
		"generate-request",
		"oasDoesNotExist.yml",
	})

	// Act
	err := cmd.Execute()

	//Assert
	assert.Error(t, err)
}

func Test_WhenFileDoesNotHaveAValidExtension_ShouldReturnError(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupGenRequestTest(t, "oas.txt", "../testOas.yml")

	// Act
	err := cmd.Execute()

	//Assert
	assert.Error(t, err, "Should return error informing user the file extension is invalid")
}

func Test_WhenGenerateRequestIsGivenValidArguments_ShouldNotError(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupGenRequestTest(t, "oas.yml", "../testOas.yml")

	// Act
	err := cmd.Execute()

	//Assert
	assert.NoError(t, err, "Should be no error when 1 arg and correct command is entered")
}
