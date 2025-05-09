package errors

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyHTTPFile = errors.New("HTTP file is empty")
)

type InvalidHTTPMethodError struct {
	Method string
}

func (e *InvalidHTTPMethodError) Error() string {
	return fmt.Sprintf("invalid HTTP method '%s', please use GET, POST, PUT, or DELETE", e.Method)
}

func NewInvalidHTTPMethodError(method string) error {
	return &InvalidHTTPMethodError{Method: method}
}

type InvalidHTTPUrlError struct {
	Url string
}

func (e *InvalidHTTPUrlError) Error() string {
	return fmt.Sprintf("invalid HTTP url '%s'", e.Url)
}

func NewInvalidHTTPUrlError(url string) error {
	return &InvalidHTTPUrlError{Url: url}
}
