package errors

import "fmt"

type HookError struct {
	HookField string
	HookValue any
	Code      ErrorCode
	Msg       string
}

const (
	ErrCodeHookBodyFieldDoesNotExist ErrorCode = "HookBodyFieldNotFound"
)

func (e *HookError) Error() string {
	return fmt.Sprintf("hook error: %s %s: %s", e.HookField, e.HookValue, e.Msg)
}

func NewHookError(hookField string, hookValue any, code ErrorCode, msg string) error {
	return &HookError{
		HookField: hookField,
		HookValue: hookValue,
		Code:      code,
		Msg:       msg,
	}
}

func NewMissingHookFieldError(hookField string) error {
	return NewHookError(hookField, nil, ErrCodeHookBodyFieldDoesNotExist, fmt.Sprintf("the field %s in the hooks file does not exist in the request body you are trying to modify", hookField))
}
