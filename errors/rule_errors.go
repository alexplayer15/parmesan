package errors

import (
	"errors"
	"fmt"
)

const (
	ErrMissingInjectionRule ErrorCode = "MissingInjectionRule"
	ErrMissingHeaderValue   ErrorCode = "MissingHeaderValue"
	ErrMissingBodyValue     ErrorCode = "MissingBodyValue"
)

var (
	ErrNoMatchingRequestsInRulesFile = errors.New("no request associated with the requests defined in the rules")
)

type RuleError struct {
	RuleField string
	ErrorCode ErrorCode
	Msg       string
}

func (e *RuleError) Error() string {
	return fmt.Sprintf("rule error: %s: %s", e.RuleField, e.Msg)
}

func NewRuleError(ruleField string, errorCode ErrorCode, msg string) error {
	return &RuleError{
		RuleField: ruleField,
		ErrorCode: errorCode,
		Msg:       msg,
	}
}

func MissingInjectionRuleError(ruleField string) error {
	return NewRuleError(ruleField, ErrMissingInjectionRule, fmt.Sprintf("you have not defined any injection rules for %s", ruleField))
}

func MissingHeaderValue(ruleField string) error {
	return NewRuleError(ruleField, ErrMissingHeaderValue, fmt.Sprintf("injection failed: missing value for header %s", ruleField))
}

func MissingBodyValue(ruleField string) error {
	return NewRuleError(ruleField, ErrMissingBodyValue, fmt.Sprintf("injection failed: missing value for body path %s", ruleField))
}
