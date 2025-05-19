package request_sender_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/errors"
	"github.com/alexplayer15/parmesan/request_sender"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WhenHttpFileIsEmpty_ShouldError(t *testing.T) {
	//Arrange
	httpRequestFile := ""

	//Act
	_, err := request_sender.ParseHttpRequestFile(httpRequestFile)

	//Assert
	assert.ErrorIs(t, err, errors.ErrEmptyHTTPFile)
}

func Test_WhenInvalidMethodIsProvidedInHttpRequestFile_ShouldError(t *testing.T) {
	//Arrange
	httpRequestFile := test_builder.NewHTTPRequestBuilder().WithMethod("PASTA").Build()

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
	httpRequestFile := test_builder.NewHTTPRequestBuilder().WithURL("bananattp://localhost:8081").Build()

	//Act
	_, err := request_sender.ParseHttpRequestFile(httpRequestFile)

	//Assert
	var ve *errors.ValidationError
	require.ErrorAs(t, err, &ve)
	assert.Equal(t, "url", ve.Param)
	assert.Equal(t, "bananattp://localhost:8081", ve.Value)
	assert.Equal(t, errors.ErrCodeInvalidPrefix, ve.Code)
}
