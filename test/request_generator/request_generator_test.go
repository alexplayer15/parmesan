package request_generator_tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexplayer15/parmesan/commands"
	"github.com/alexplayer15/parmesan/request_generator"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	test_data "github.com/alexplayer15/parmesan/test_oas_data"
	"github.com/stretchr/testify/assert"
)

var testOasDir = "../../test_oas_data/"

func Test_WhenOASHasStringPropertyWithAnExample_ShouldReturnExampleString(t *testing.T) {
	//Arrange
	testOasPath := filepath.Join(testOasDir, "simpleOAS.yml")
	outputFilePath := "simpleOAS.http"

	commands.RootCmd.SetArgs([]string{"generate-request", testOasPath})

	//Act
	err := commands.RootCmd.Execute()

	//Assert
	_, statErr := os.Stat(outputFilePath)
	assert.NoError(t, statErr, "expected output file to be created")

	outputContent, readErr := os.ReadFile(outputFilePath)
	assert.NoError(t, readErr, "expected output file to be readable")

	assert.Contains(t, string(outputContent), `"name": "Alex"`)

	assert.NoError(t, err)
	os.Remove(outputFilePath)
}

func Test_WhenOASHasIntegerPropertyWithAnExample_ShouldReturnExampleInteger(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("age").
		WithType("integer").
		WithExample(25).
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"age": 25`)
}

func Test_WhenOASHasObjectProperty_ShouldReturnCorrectObjectInRequest(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithType("object").
		AddProperty(test_builder.NewPropertyBuilder().
			WithName("university").
			WithType("string").
			WithExample("University of Manchester").
			Build(),
		).
		AddProperty(test_builder.NewPropertyBuilder().
			WithName("degree").
			WithType("string").
			WithExample("Chemical Engineering").
			Build(),
		).
		AddProperty(test_builder.NewPropertyBuilder().
			WithName("grade").
			WithType("string").
			WithExample("2:1").
			Build(),
		).
		Build()

	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	result, err := request_generator.GenerateHttpRequest(oas)

	assert.NoError(t, err)
	assert.Contains(t, result, `"education"`)
	assert.Contains(t, result, `"university": "University of Manchester"`)
	assert.Contains(t, result, `"degree": "Chemical Engineering"`)
	assert.Contains(t, result, `"grade": "2:1"`)
}

func Test_WhenOASHasArrayPropertyContainingStringsWithAnExample_ShouldReturnExampleArray(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("favouriteColours").
		WithType("array").
		WithExample([]any{"blue", "purple", "red"}).
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"favouriteColours": ["blue", "purple", "red"]`)
}

func Test_WhenOASHasArrayPropertyContainingIntsWithAnExample_ShouldReturnExampleArray(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("favouriteNumbers").
		WithType("array").
		WithExample([]any{15, 15, 15}).
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"favouriteNumbers": [15, 15, 15]`)
}

func Test_WhenOASHasArrayPropertyContainingAnObjectWithAnExample_ShouldReturnExampleArray(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	educationItemSchema := test_builder.NewSchemaBuilder().
		AddProperty(test_builder.NewPropertyBuilder().
			WithName("university").
			WithType("string").
			WithExample("University of Manchester").
			Build()).
		AddProperty(test_builder.NewPropertyBuilder().
			WithName("degree").
			WithType("string").
			WithExample("Chemical Engineering").
			Build()).
		AddProperty(test_builder.NewPropertyBuilder().
			WithName("grade").
			WithType("string").
			WithExample("2:1").
			Build()).
		Build()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithType("array").
		WithItems(educationItemSchema).
		Build()

	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"education":`)
	assert.Contains(t, result, `"university": "University of Manchester"`)
	assert.Contains(t, result, `"degree": "Chemical Engineering"`)
	assert.Contains(t, result, `"grade": "2:1"`)
}
