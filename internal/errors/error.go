package errors

import (
	"fmt"
)

type ErrorCode string

const (
	ErrConfigLoad     ErrorCode = "CONFIG_LOAD_FAILED"
	ErrConfigNotFound ErrorCode = "CONFIG_NOT_FOUND"
)

type CoreError struct {
	Code    ErrorCode
	Scope   string
	Message string
	Err     error
}

func (e *CoreError) Error() string {
	base := fmt.Sprintf("[%s]", e.Code)

	if e.Scope != "" {
		base += fmt.Sprintf(" (%s)", e.Scope)
	}

	if e.Message != "" {
		base += " " + e.Message
	}

	if e.Err != nil {
		base += fmt.Sprintf(": %v", e.Err)
	}

	return base
}

func (e *CoreError) Unwrap() error {
	return e.Err
}

func NewConfigLoadError(scope string, msg string, err error) error {
	return &CoreError{
		Code:    ErrConfigLoad,
		Scope:   scope,
		Message: msg,
		Err:     err,
	}
}

func NewConfigNotFoundError(scope string, msg string) error {
	return &CoreError{
		Code:    ErrConfigNotFound,
		Scope:   scope,
		Message: msg,
	}
}
