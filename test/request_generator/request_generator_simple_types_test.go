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
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"name": "Alex"`)
}

func Test_WhenOASHasStringPropertyWithoutAnExample_ShouldReturnFallbackString(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("name").
		WithType("string").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"name": "example value"`)
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
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"age": 25`)
}

func Test_WhenOASHasIntegerPropertyWithoutAnExample_ShouldReturnExampleZero(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("age").
		WithType("integer").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	assert.Contains(t, result, `"age": 0`)
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
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	//have a think about how much spacing matters. Test is fragile as stands
	assert.Contains(t, result, `"birthday": "2022-11-23"`)
}

func TestWhenOASHasPropertyOfTypeStringWithFormatDateAndNoExample_ShouldReturnFallbackDate(t *testing.T) {
	// Arrange
	oas := test_data.BaseOAS()

	propName, propValue := test_builder.NewPropertyBuilder().
		WithName("birthday").
		WithType("string").
		WithFormat("date").
		Build()
	oas.Paths["/users"]["post"].RequestBody.Content["application/json"].Schema.Properties[propName] = propValue

	//Act
	result, err := request_generator.GenerateHttpRequest(oas, 0)

	//Assert
	assert.NoError(t, err)
	//have a think about how much spacing matters. Test is fragile as stands
	assert.Contains(t, result, `"birthday": "2022-01-01"`)
}

//TO DO: Write tests to ensure http file is outputted as expected, including formatting.
