package flag_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/test_helpers"
	"github.com/stretchr/testify/assert"
)

func Test_WhenWithServerFlagIsUsedWithCorrectArgument_ShouldSelectCorrectServer(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupGenRequestTest(t, "oas.yml", "../testOas.yml", "--with-server", "0")

	// Act
	err := cmd.Execute()

	//Assert
	assert.NoError(t, err)
}

func Test_WhenWithSeverFlagIsGivenAnIndexNotFoundInOAS_ShouldErrorAndInformUser(t *testing.T) {
	//Arrange
	cmd, _ := test_helpers.SetupGenRequestTest(t, "oas.yml", "../testOas.yml", "--with-server", "2")

	// Act
	err := cmd.Execute()

	//Assert
	assert.EqualError(t, err, "invalid server index 2: There are 1 servers available starting from 0")
}
