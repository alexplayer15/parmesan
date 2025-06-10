package chain_logic_tests

import (
	"testing"

	"github.com/alexplayer15/parmesan/chain_logic"
	"github.com/alexplayer15/parmesan/data"
	"github.com/alexplayer15/parmesan/errors"
	test_builder "github.com/alexplayer15/parmesan/test_builders"
	"github.com/alexplayer15/parmesan/test_helpers"
	"github.com/stretchr/testify/assert"
)

func Test_WhenGivenAValidRulesFileAndHttpRequest_ShouldOrderRequestsAccordingToRulesFile(t *testing.T) {
	//Arrange
	requestOne := test_builder.NewParsedHttpRequest().
		WithMethod("GET").
		WithUrl("http://localhost:8081/example-path-one").
		WithHeader("Content", "application/json").
		WithJSONBody(test_helpers.NewTestRequest("Alex", 25)).
		Build()

	requestTwo := test_builder.NewParsedHttpRequest().
		WithMethod("POST").
		WithUrl("http://localhost:9090/example-path-two").
		WithHeader("Content", "application/json").
		WithJSONBody(test_helpers.NewTestRequest("Mia", 27)).
		Build()

	requests := []data.Request{requestOne, requestTwo}

	rules := test_builder.NewRuleSetBuilder().
		WithRule("first", "/example-path-two", "POST").
		EndRule().
		WithRule("second", "/example-path-one", "GET").
		EndRule().
		Build()

	//Act
	orderedRequests, err := chain_logic.OrderRequests(requests, rules)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, "POST", orderedRequests[0].Method)
	assert.Equal(t, "http://localhost:9090/example-path-two", orderedRequests[0].Url)
	assert.Equal(t, "GET", orderedRequests[1].Method)
	assert.Equal(t, "http://localhost:8081/example-path-one", orderedRequests[1].Url)
}

func Test_WhenOrderingRequests_IfNoRequestsMatch_ShouldErrorAndInformUser(t *testing.T) {
	//Arrange
	requestOne := test_builder.NewParsedHttpRequest().
		WithMethod("GET").
		WithUrl("http://localhost:8081/example-path-one").
		WithHeader("Content", "application/json").
		WithJSONBody(test_helpers.NewTestRequest("Alex", 25)).
		Build()

	requests := []data.Request{requestOne}

	rules := test_builder.NewRuleSetBuilder().
		WithRule("second", "/example-path-two", "POST").
		EndRule().
		Build()

	//Act
	_, err := chain_logic.OrderRequests(requests, rules)

	//Assert
	assert.ErrorIs(t, err, errors.ErrNoMatchingRequestsInRulesFile)
}

func Test_WhenApplyingInjectionRules_IfNoInjectionRulesAreDefined_ShouldError(t *testing.T) {
	//Arrange
	request := test_builder.NewParsedHttpRequest().
		WithMethod("GET").
		WithUrl("http://localhost:8081/example-path-one").
		WithHeader("Content", "application/json").
		WithJSONBody(test_helpers.NewTestRequest("Alex", 25)).
		Build()

	rules := test_builder.NewRuleSetBuilder().
		WithRule("first", "/example-path-one", "GET").
		EndRule().
		Build()

	extractedValues := make(map[string]any)

	//Act
	_, err := chain_logic.ApplyInjectionRules(request, rules, extractedValues)

	//Assert
	var ruleError *errors.RuleError
	assert.ErrorAs(t, err, &ruleError)
	assert.Equal(t, errors.ErrMissingInjectionRule, ruleError.ErrorCode)
}

func Test_WhenApplyingInjectionRules_IfMissingHeaderValueFromInjectionRule_ShouldError(t *testing.T) {
	//Arrange
	request := test_builder.NewParsedHttpRequest().
		WithMethod("GET").
		WithUrl("http://localhost:8081/example-path-one").
		WithHeader("Content", "application/json").
		WithJSONBody(test_helpers.NewTestRequest("Alex", 25)).
		Build()

	rules := test_builder.NewRuleSetBuilder().
		WithRule("first", "/example-path-one", "GET").
		AddHeaderInject("example-name", "example-from").
		EndRule().
		Build()

	extractedValues := make(map[string]any)

	//Act
	_, err := chain_logic.ApplyInjectionRules(request, rules, extractedValues)

	//Assert
	var ruleError *errors.RuleError
	assert.ErrorAs(t, err, &ruleError)
	assert.Equal(t, errors.ErrMissingHeaderValue, ruleError.ErrorCode)
}

func Test_WhenApplyingInjectionRules_IfMissingBodyValueFromInjectionRule_ShouldError(t *testing.T) {
	//Arrange
	request := test_builder.NewParsedHttpRequest().
		WithMethod("GET").
		WithUrl("http://localhost:8081/example-path-one").
		WithHeader("Content", "application/json").
		WithJSONBody(test_helpers.NewTestRequest("Alex", 25)).
		Build()

	rules := test_builder.NewRuleSetBuilder().
		WithRule("first", "/example-path-one", "GET").
		AddBodyInject("example-name", "example-from", "string").
		EndRule().
		Build()

	extractedValues := make(map[string]any)

	//Act
	_, err := chain_logic.ApplyInjectionRules(request, rules, extractedValues)

	//Assert
	var ruleError *errors.RuleError
	assert.ErrorAs(t, err, &ruleError)
	assert.Equal(t, errors.ErrMissingBodyValue, ruleError.ErrorCode)
}
