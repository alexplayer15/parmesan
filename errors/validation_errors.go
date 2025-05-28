package errors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrCodeInvalidMethod ErrorCode = "InvalidMethod"
	ErrCodeInvalidPrefix ErrorCode = "InvalidPrefix"
	ErrCodeURLParsing    ErrorCode = "URLParsing"
	ErrCodeMissingHost   ErrorCode = "MissingHost"
)

var (
	ErrEmptyHTTPFile = errors.New("HTTP file is empty")
)

type ValidationError struct {
	Param string
	Value string
	Code  ErrorCode
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s (%s): %s", e.Param, e.Value, e.Msg)
}

func NewValidationError(param, value string, code ErrorCode, msg string) error {
	return &ValidationError{
		Param: param,
		Value: value,
		Code:  code,
		Msg:   msg,
	}
}

func NewInvalidMethodError(param, value string) error {
	return NewValidationError(param, value, ErrCodeInvalidMethod, fmt.Sprintf("%s is an invalid HTTP method", value))
}

func NewInvalidPrefixError(param, value string) error {
	return NewValidationError(param, value, ErrCodeInvalidPrefix, "missing http:// or https:// prefix")
}

func NewURLParsingError(param, value string) error {
	return NewValidationError(param, value, ErrCodeURLParsing, "error parsing URL")
}

func NewMissingHostError(param, value string) error {
	return NewValidationError(param, value, ErrCodeMissingHost, "host is missing from URL")
}
