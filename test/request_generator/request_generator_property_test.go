package request_generator_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/request_generator"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	test_data "github.com/alexplayer15/parmesan/test_oas_data"
	"github.com/stretchr/testify/assert"
)

func Test_WhenOASHasStringPropertyWithAnExample_ShouldReturnExampleString(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("name").
		WithType("string").
		WithExample("Alex").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"name": "Alex"`)
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

	result, err := request_generator.GenerateHttpRequest(oas)

	//to do - change to test for structure also
	assert.NoError(t, err)
	assert.Contains(t, result, `"education"`)
	assert.Contains(t, result, `"university": "University of Manchester"`)
	assert.Contains(t, result, `"degree": "Chemical Engineering"`)
	assert.Contains(t, result, `"grade": "2:1"`)
}

func TestWhenOASHasPropertyOfTypeStringWithFormatDateAndAnExample_ShouldReturnDate(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("birthday").
		WithType("string").
		WithFormat("date").
		WithExample("2022-11-23").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	//have a think about how much spacing matters. Test is fragile as stands
	assert.Contains(t, result, `"birthday": "2022-11-23"`)
}
