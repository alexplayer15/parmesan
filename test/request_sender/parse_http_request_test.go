package request_sender_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/errors"
	"github.com/alexplayer15/parmesan/request_sender"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	"github.com/stretchr/testify/assert"
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
	assert.ErrorAs(t, err, &ve)
	assert.Equal(t, "method", ve.ParameterName)
	assert.Equal(t, "PASTA", ve.ParameterValue)
}
