package request_generator_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/request_generator"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	"github.com/alexplayer15/parmesan/test_helpers"
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
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"favouriteColours": ["blue", "purple", "red"]`)
}

func Test_WhenOASHasArrayPropertyContainingStringsWithoutAnExample_ShouldReturnFallbackArray(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("favouriteColours").
		WithType("array").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"favouriteColours": []`)
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
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"favouriteNumbers": [15, 15, 15]`)
}

func Test_WhenOASHasArrayPropertyContainingIntsWithoutAnExample_ShouldReturnFallbackArray(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("favouriteNumbers").
		WithType("array").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"favouriteNumbers": []`)
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
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.NoError(t, err)
	body, err := test_helpers.ExtractBody(result)
	assert.NoError(t, err)
	test_helpers.AssertJSONHasArrayWithObject(t, body, "education", []string{"university", "degree", "grade"})
	test_helpers.AssertJSONExamplesForObjectsInAnArray(t, body, "education", map[string]any{
		"university": "University of Manchester",
		"degree":     "Chemical Engineering",
		"grade":      "2:1",
	})
}

func Test_WhenOASHasArrayPropertyContainingAnObjectWithoutAnExample_ShouldReturnFallbackValuesBasedOnPropertyType(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	educationItemSchema := test_builder.NewSchemaBuilder().
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("university").
			WithType("string").
			Build()).
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("degree").
			WithType("string").
			Build()).
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("grade").
			WithType("integer").
			Build()).
		Build()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithType("array").
		WithItems(educationItemSchema).
		Build()

	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.NoError(t, err)
	body, err := test_helpers.ExtractBody(result)
	assert.NoError(t, err)
	test_helpers.AssertJSONHasArrayWithObject(t, body, "education", []string{"university", "degree", "grade"})
	test_helpers.AssertJSONExamplesForObjectsInAnArray(t, body, "education", map[string]any{
		"university": "example value",
		"degree":     "example value",
		"grade":      float64(0),
	})
}

func Test_WhenOASHasObjectPropertyWithExamples_ShouldReturnObject(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithType("object").
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("university").
			WithType("string").
			WithExample("University of Manchester").
			Build(),
		).
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("degree").
			WithType("string").
			WithExample("Chemical Engineering").
			Build(),
		).
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("grade").
			WithType("string").
			WithExample("2:1").
			Build(),
		).
		Build()

	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	body, err := test_helpers.ExtractBody(result)
	assert.NoError(t, err)
	test_helpers.AssertJSONHasObject(t, body, "education", []string{"university", "degree", "grade"})
	test_helpers.AssertJSONExamplesForObject(t, body, "education", map[string]any{
		"university": "University of Manchester",
		"degree":     "Chemical Engineering",
		"grade":      "2:1",
	})
}

func Test_WhenOASHasObjectPropertyWithoutExamples_ShouldReturnFallbackValues(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithType("object").
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("university").
			WithType("string").
			Build(),
		).
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("degree").
			WithType("string").
			Build(),
		).
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("grade").
			WithType("string").
			Build(),
		).
		Build()

	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	body, err := test_helpers.ExtractBody(result)
	assert.NoError(t, err)
	test_helpers.AssertJSONHasObject(t, body, "education", []string{"university", "degree", "grade"})
	test_helpers.AssertJSONExamplesForObject(t, body, "education", map[string]any{
		"university": "example value",
		"degree":     "example value",
		"grade":      "example value",
	})
}
