package errors

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyHTTPFile = errors.New("HTTP file is empty")
)

type ValidationError struct {
	ParameterName  string
	ParameterValue string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s is not a valid %s", e.ParameterValue, e.ParameterName)
}

func NewValidationError(parameterName string, parameterValue string) error {
	return &ValidationError{ParameterName: parameterName, ParameterValue: parameterValue}
}
