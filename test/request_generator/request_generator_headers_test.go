package request_generator_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/request_generator"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	"github.com/alexplayer15/parmesan/test_helpers"
	test_data "github.com/alexplayer15/parmesan/test_oas_data"
	"github.com/stretchr/testify/assert"
)

func Test_WhenOASHasHeadersWithExamples_ShouldReturnExampleHeaders(t *testing.T) {
	//Arrange
	oas := test_data.BaseOAS()
	param := test_builder.NewParameterBuilder().
		WithName("X-Session-ID").
		WithIn("header").
		WithExample("1").
		Build()

	method := oas.Paths["/users"]["post"]
	method.Parameters = append(method.Parameters, *param)
	oas.Paths["/users"]["post"] = method

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	headerMap, err := test_helpers.ExtractHeaders(t, result)
	assert.NoError(t, err)
	assert.Equal(t, "1", headerMap["X-Session-ID"])
}

func Test_WhenOASHasHeadersWithoutExamples_ShouldReturnFallbackValues(t *testing.T) {
	//Arrange
	oas := test_data.BaseOAS()
	param := test_builder.NewParameterBuilder().
		WithName("X-Session-ID").
		WithIn("header").
		Build()

	method := oas.Paths["/users"]["post"]
	method.Parameters = append(method.Parameters, *param)
	oas.Paths["/users"]["post"] = method

	//Act
	result, err := request_generator.GenerateHttpRequest(oas)

	//Assert
	assert.NoError(t, err)
	headerMap, err := test_helpers.ExtractHeaders(t, result)
	assert.NoError(t, err)
	assert.Equal(t, "default-value", headerMap["X-Session-ID"])
}
