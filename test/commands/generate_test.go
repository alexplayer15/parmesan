package command_tests

import (
	"testing"
	"os"
	"github.com/alexplayer15/parmesan/commands"

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
	testOas := "oas.yml"
	os.WriteFile(testOas, []byte("dummy content"), 0644)
	defer os.Remove(testOas)
	commands.RootCmd.SetArgs([]string{"generate-request", testOas})

	//Act
	err := commands.RootCmd.Execute()

	//Assert
	assert.Error(t, err, "Command should fail if no argument is given")
}
func Test_WhenGenerateRequestIsGivenValidArguments_ShouldNotError(t *testing.T) {
	//Arrange 
	testOas := "oas.yml"
	os.WriteFile(testOas, []byte("dummy content"), 0644)
	defer os.Remove(testOas)
	commands.RootCmd.SetArgs([]string{"generate-request", testOas})

	//Act 
	err := commands.RootCmd.Execute()

	//Assert
	assert.NoError(t, err, "Should be no error when 1 arg and correct command is entered")
}
func Test_WhenOASDoesNotExist_ShouldReturnNoFileFoundError(t *testing.T){
	//Arrange 
	commands.RootCmd.SetArgs([]string{"generate-request", "oasDoesNotExist.yml"})

	//Act 
	err := commands.RootCmd.Execute()

	//Assert
	assert.EqualError(t, err, "file does not exist")
}

func Test_WhenFileDoesNotHaveAValidExtension_ShouldReturnError(t *testing.T){
	//Arrange 
	dummyFilename := "dummy.txt"
	os.WriteFile(dummyFilename, []byte("dummy content"), 0644)
	defer os.Remove(dummyFilename)
	commands.RootCmd.SetArgs([]string{"generate-request", dummyFilename})

	//Act 
	err := commands.RootCmd.Execute()

	//Assert
	assert.EqualError(t, err, "OAS must be a JSON or YAML file")
}

func Test_WhenOASExistsAndIsValid_ShouldReturnNoError(t *testing.T){
	//Arrange 
	commands.RootCmd.SetArgs([]string{"generate-request", "oasDoesNotExist.yml"})

	//Act 
	err := commands.RootCmd.Execute()

	//Assert
	assert.EqualError(t, err, "file does not exist")
}
