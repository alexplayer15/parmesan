package request_sender_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/errors"
	"github.com/alexplayer15/parmesan/request_sender"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	"github.com/alexplayer15/parmesan/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WhenHttpFileIsEmpty_ShouldError(t *testing.T) {
	//Arrange
	httpRequest := ""

	//Act
	_, err := request_sender.ParseHttpRequestFile(httpRequest)

	//Assert
	assert.ErrorIs(t, err, errors.ErrEmptyHTTPFile)
}

func Test_WhenInvalidMethodIsProvidedInHttpRequestFile_ShouldError(t *testing.T) {
	//Arrange
	httpRequest := test_builder.NewHTTPRequestBuilder().WithMethod("PASTA").Build()

	httpRequestFile := test_builder.NewHTTPFileBuilder().
		WithHTTPRequests([]*test_builder.HTTPRequestBuilder{httpRequest}).
		Build()

	//Act
	_, err := request_sender.ParseHttpRequestFile(httpRequestFile)

	//Assert
	var ve *errors.ValidationError
	require.ErrorAs(t, err, &ve)
	assert.Equal(t, "method", ve.Param)
	assert.Equal(t, "PASTA", ve.Value)
	assert.Equal(t, errors.ErrCodeInvalidMethod, ve.Code)
}

func Test_WhenInvalidURLPrefixIsProvided_ShouldReturnValidationError(t *testing.T) {
	//Arrange
	httpRequest := test_builder.NewHTTPRequestBuilder().WithURL("bananattp://localhost:8081").Build()
	httpRequestFile := test_builder.NewHTTPFileBuilder().
		WithHTTPRequests([]*test_builder.HTTPRequestBuilder{httpRequest}).
		Build()

	//Act
	_, err := request_sender.ParseHttpRequestFile(httpRequestFile)

	//Assert
	var ve *errors.ValidationError
	require.ErrorAs(t, err, &ve)
	assert.Equal(t, "url", ve.Param)
	assert.Equal(t, "bananattp://localhost:8081", ve.Value)
	assert.Equal(t, errors.ErrCodeInvalidPrefix, ve.Code)
}

func Test_WhenURLIsMissingHost_ShouldReturnValidationError(t *testing.T) {
	//Arrange
	httpRequest := test_builder.NewHTTPRequestBuilder().WithURL("http://:8081").Build()

	httpRequestFile := test_builder.NewHTTPFileBuilder().
		WithHTTPRequests([]*test_builder.HTTPRequestBuilder{httpRequest}).
		Build()

	//Act
	_, err := request_sender.ParseHttpRequestFile(httpRequestFile)

	//Assert
	var ve *errors.ValidationError
	require.ErrorAs(t, err, &ve)
	assert.Equal(t, "url", ve.Param)
	assert.Equal(t, "http://:8081", ve.Value)
	assert.Equal(t, errors.ErrCodeMissingHost, ve.Code)
}

func Test_WhenValidHTTPFileIsGenerated_ShouldReturnValidRequest(t *testing.T) {
	//Arrange
	testRequestBody := test_helpers.NewTestRequest("Alex", 25)
	httpRequest := test_builder.NewHTTPRequestBuilder().
		WithSummary("Test Summary").
		WithMethod("GET").
		WithURL("http://localhost:8081").
		WithHeader("Content", "application/json").
		WithJSONBody(testRequestBody).
		Build()

	httpRequestFile := test_builder.NewHTTPFileBuilder().
		WithHTTPRequests([]*test_builder.HTTPRequestBuilder{httpRequest}).
		Build()

	//Act
	httpRequests, err := request_sender.ParseHttpRequestFile(httpRequestFile)

	//Assert
	assert.NoError(t, err)
	expectedJSON := `{
		"Name": "Alex",
		"Age": 25
	}`
	assert.JSONEq(t, expectedJSON, httpRequests[0].Body)
	assert.Equal(t, "application/json", httpRequests[0].Headers["Content"])
	assert.Equal(t, "GET", httpRequests[0].Method)
	assert.Equal(t, "http://localhost:8081", httpRequests[0].Url)
}

func Test_WhenMultipleRequestsAreParsed_ShouldReturnDetailsCorrectly(t *testing.T) {
	//Arrange
	testRequestBodyOne := test_helpers.NewTestRequest("Alex", 25)
	httpRequestOne := test_builder.NewHTTPRequestBuilder().
		WithSummary("Test Summary").
		WithMethod("GET").
		WithURL("http://localhost:8081").
		WithHeader("Content", "application/json").
		WithJSONBody(testRequestBodyOne).
		Build()

	testRequestBodyTwo := test_helpers.NewTestRequest("Mia", 27)
	httpRequestTwo := test_builder.NewHTTPRequestBuilder().
		WithSummary("Test Summary").
		WithMethod("GET").
		WithURL("http://localhost:8081").
		WithHeader("Content", "application/json").
		WithJSONBody(testRequestBodyTwo).
		Build()

	httpRequestFile := test_builder.NewHTTPFileBuilder().
		WithHTTPRequests([]*test_builder.HTTPRequestBuilder{httpRequestOne, httpRequestTwo}).
		Build()

	//Act
	httpRequests, err := request_sender.ParseHttpRequestFile(httpRequestFile)

	//Assert
	assert.NoError(t, err)
	expectedJSONBodyOne := `{
		"Name": "Alex",
		"Age": 25
	}`
	assert.JSONEq(t, expectedJSONBodyOne, httpRequests[0].Body)
	assert.Equal(t, "application/json", httpRequests[0].Headers["Content"])
	assert.Equal(t, "GET", httpRequests[0].Method)
	assert.Equal(t, "http://localhost:8081", httpRequests[0].Url)

	expectedJSONBodyTwo := `{
		"Name": "Mia",
		"Age": 27
	}`

	assert.JSONEq(t, expectedJSONBodyTwo, httpRequests[1].Body)
	assert.Equal(t, "application/json", httpRequests[1].Headers["Content"])
	assert.Equal(t, "GET", httpRequests[1].Method)
	assert.Equal(t, "http://localhost:8081", httpRequests[1].Url)
}
