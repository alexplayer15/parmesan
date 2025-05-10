package errors

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyHTTPFile = errors.New("HTTP file is empty")
)

type ValidationError struct {
	InvalidParameter string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error caused by: %s", e.InvalidParameter)
}

func NewValidationError(invalidParameter string) error {
	return &ValidationError{InvalidParameter: invalidParameter}
}
