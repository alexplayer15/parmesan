package request_generator_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/request_generator"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	"github.com/alexplayer15/parmesan/test_helpers"
	test_data "github.com/alexplayer15/parmesan/test_oas_data"
	"github.com/stretchr/testify/assert"
)

func Test_WhenOASPropertyUsesOneOf_ShouldReturnTheObjectInTheFirstReference(t *testing.T) {
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
		WithOneOfRefs([]string{
			"#/components/schemas/Education",
			"#/components/schemas/Education",
		}).
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
	test_helpers.AssertJSONHasXAmountOfObjects(t, body, 1)
}

func Test_WhenOASPropertyUsesAllOfAndReferencesTheSameSchemaTwice_ShouldReturnOneObject(t *testing.T) {
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
		WithAllOfRefs([]string{
			"#/components/schemas/Education",
			"#/components/schemas/Education",
		}).
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
	test_helpers.AssertJSONHasXAmountOfObjects(t, body, 1)
}

func Test_WhenOASPropertyUsesAllOfReferencingMultipleDifferentObjects_ShouldReturnAllReferencedObjects(t *testing.T) {
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
	jobItemSchema := test_builder.NewSchemaBuilder().
		WithType("object").
		WithProperty(test_builder.NewPropertyBuilder().
			WithName("job").
			WithType("string").
			WithExample("Developer").
			Build()).
		Build()

	oas.Components.Schemas["Education"] = *educationItemSchema
	oas.Components.Schemas["Job"] = *jobItemSchema

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("career").
		WithAllOfRefs([]string{
			"#/components/schemas/Education",
			"#/components/schemas/Job",
		}).
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	body, err := test_helpers.ExtractBody(result)
	assert.NoError(t, err)
	test_helpers.AssertJSONHasObject(t, body, "career", []string{"university", "degree", "grade", "job"})
	test_helpers.AssertJSONExamplesForObject(t, body, "career", map[string]any{
		"university": "University of Manchester",
		"degree":     "Chemical Engineering",
		"grade":      "2:1",
		"job":        "Developer",
	})
	test_helpers.AssertJSONHasXAmountOfObjects(t, body, 1)
}
