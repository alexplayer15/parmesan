package output_tests

import (
	"path/filepath"
	"testing"

	"github.com/alexplayer15/parmesan/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WhenOutputFlagIsUsedIfOutputIsValid_ShouldSuccessfullyOutputFileToIntendedLocation(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupCommandTest(t, "oas.yml", "../testOas.yml", "--output", "http")

	// Act
	err := cmd.Execute()

	//Assert
	assert.NoError(t, err)
	require.FileExists(t, "http/oas.http")
}

func Test_WhenOutputFlagIsNotGivenAndCommandIsSuccessful_ShouldCreateOutputFileInCurrentDirectory(t *testing.T) {
	// Arrange
	cmd, tmpDir := test_helpers.SetupCommandTest(t, "oas.yml", "../testOas.yml")

	// Act
	err := cmd.Execute()

	// Assert
	assert.NoError(t, err)
	require.FileExists(t, filepath.Join(tmpDir, "oas.http"))
}

func Test_WhenOutputFlagArgumentIsAnExistingFile_ShouldErrorAndInformTheUser(t *testing.T) {
	// Arrange
	cmd, _ := test_helpers.SetupCommandTest(t, "oas.yml", "../testOas.yml", "--output", "oas.yml")

	// Act
	err := cmd.Execute()

	//Assert
	assert.EqualError(t, err, "invalid output directory: path exists and is a file")
}

func Test_WhenOutputFlagIsUsedButNoArgumentIsGiven_ShouldReturnErrorAndInformUser(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupCommandTest(t, "oas.yml", "../testOas.yml", "--output")

	//Act
	err := cmd.Execute()

	//Assert
	assert.EqualError(t, err, "flag needs an argument: --output")
}

func Test_WhenOutputFlagIsUsedBAndMoreThanOneArgumentIsGiven_ShouldReturnErrorAndInformUser(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupCommandTest(t, "oas.yml", "../testOas.yml", "--output", "dummyArg", "dummyArg2")

	//Act
	err := cmd.Execute()

	//Assert
	assert.EqualError(t, err, "accepts 1 arg(s), received 2")
}
