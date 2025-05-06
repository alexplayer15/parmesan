package request_generator_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/request_generator"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	"github.com/alexplayer15/parmesan/test_helpers"
	test_data "github.com/alexplayer15/parmesan/test_oas_data"
	"github.com/stretchr/testify/assert"
)

func Test_WhenOASHasAnArrayPropertyReferencingASchemaOfTypeObject_ShouldReturnArrayWithResolvedObject(t *testing.T) {
	//Arrange
	oas := test_data.BaseOAS()

	educationItemSchema := test_builder.NewSchemaBuilder().
		WithType("object").
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

	oas.Components.Schemas["Education"] = *educationItemSchema

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithType("array").
		WithItemsRef("#/components/schemas/Education").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
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

func Test_WhenOASHasArrayPropertyReferencingASchemaWhichDoesNotSpecifyObjectTypeButContainsPropertiesWithExamples_ShouldReturnExampleSchema(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	educationItemSchema := test_builder.NewSchemaBuilder().
		WithType("").
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

	oas.Components.Schemas["Education"] = *educationItemSchema

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithType("array").
		WithItemsRef("#/components/schemas/Education").
		Build()

	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

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

func Test_WhenOASHasAPropertyReferencingASchemaOfTypeObject_ShouldReturnObject(t *testing.T) {
	//Arrange
	oas := test_data.BaseOAS()

	educationItemSchema := test_builder.NewSchemaBuilder().
		WithType("object").
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

	oas.Components.Schemas["Education"] = *educationItemSchema

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithRef("#/components/schemas/Education").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

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

func Test_WhenOASHasAPropertyReferencingASchemaOfTypeStringWithAnExample_ShouldReturnExampleString(t *testing.T) {
	//Arrange
	oas := test_data.BaseOAS()

	educationItemSchema := test_builder.NewSchemaBuilder().
		WithType("string").
		WithExample("University of Manchester").
		Build()

	oas.Components.Schemas["Education"] = *educationItemSchema

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithRef("#/components/schemas/Education").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"education": "University of Manchester"`)
}

func Test_WhenOASHasAPropertyReferencingASchemaOfTypeStringWithADefault_ShouldReturnDefaultString(t *testing.T) {
	//Arrange
	oas := test_data.BaseOAS()

	educationItemSchema := test_builder.NewSchemaBuilder().
		WithType("string").
		WithDefault("University of Manchester").
		Build()

	oas.Components.Schemas["Education"] = *educationItemSchema

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("education").
		WithRef("#/components/schemas/Education").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"education": "University of Manchester"`)
}

func Test_WhenOASHasASchemaWhichContainsAPropertyUsingANonExistentReference_ShouldErrorAndInformTheUser(t *testing.T) {
	//Arrange
	oas := test_data.BaseOAS()

	educationItemSchema := test_builder.NewSchemaBuilder().
		WithType("object").
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("university").
			WithType("string").
			WithExample("University of Manchester").
			Build()).
		Build()

	requestBodySchema := test_builder.NewSchemaBuilder().
		WithType("object").
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("education").
			WithType("array").
			WithItemsRef("#/components/schemas/NonExistentSchemaName").
			Build()).
		Build()

	mediaType := oas.Paths["/users"]["post"].RequestBody.Content["application/json"]
	mediaType.Schema = *requestBodySchema
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"] = mediaType

	oas.Components.Schemas["Education"] = *educationItemSchema

	//Act
	_, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.Error(t, err)
}

func Test_WhenAnArrayPropertyReferencesASchemaWhichHasPropertiesWithDefaultValues_ShouldReturnDefaultValues(t *testing.T) {
	//Arrange
	oas := test_data.BaseOAS()

	jobItemSchema := test_builder.NewSchemaBuilder().
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("years of experience").
			WithType("integer").
			WithDefault(3).
			Build()).
		Build()

	oas.Components.Schemas["Job"] = *jobItemSchema

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("job").
		WithType("array").
		WithItemsRef("#/components/schemas/Job").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	body, err := test_helpers.ExtractBody(result)
	assert.NoError(t, err)
	test_helpers.AssertJSONHasArrayWithObject(t, body, "job", []string{"years of experience"})
	test_helpers.AssertJSONExamplesForObjectsInAnArray(t, body, "job", map[string]any{
		"years of experience": float64(3),
	})
	test_helpers.AssertJSONHasXAmountOfArrays(t, body, 1)
}
