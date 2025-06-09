package request_sender_tests

// import (
// 	"testing"

// 	"github.com/alexplayer15/parmesan/request_sender"
// 	test_builder "github.com/alexplayer15/parmesan/test_builders"
// )

// func Test_WhenIHaveAValidHTTPRequest_ShouldSuccessfullySendRequest(t *testing.T) {
// 	//Arrange
// 	httpRequest := test_builder.NewParsedHttpRequest().
// 		WithMethod("GET").
// 		WithUrl("http://localhost:8081").
// 		WithHeader("Content", "application/json").
// 		WithJSONBody(NewTestRequest("Alex", 25)).
// 		Build()

// 	//Act
// 	responseBody, statusCode, responseHeaders, err := request_sender.SendHTTPRequest(httpRequest)

// 	//Assert

// 	//Think about what behaviour we need to test here. We can check we get the response information as expected
// 	//but be careful not to test the library.

// }
