package request_generator_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/request_generator"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	test_data "github.com/alexplayer15/parmesan/test_oas_data"
	"github.com/stretchr/testify/assert"
)

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
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("university").
			WithType("string").
			WithExample("University of Manchester").
			Build()).
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("degree").
			WithType("string").
			WithExample("Chemical Engineering").
			Build()).
		WithProperty(test_builder.NewPropertyBuilder().
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
