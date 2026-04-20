package errors

import (
	"fmt"
)

type ErrorCode string

const (
	ErrConfigLoad     ErrorCode = "CONFIG_LOAD_FAILED"
	ErrConfigNotFound ErrorCode = "CONFIG_NOT_FOUND"

	ErrFlagRead   ErrorCode = "FLAG_READ_FAILED"
	ErrValidation ErrorCode = "VALIDATION_FAILED"
	ErrIO         ErrorCode = "IO_OPERATION_FAILED"
	ErrTemplate   ErrorCode = "TEMPLATE_RENDER_FAILED"
	ErrRuntime    ErrorCode = "RUNTIME_EXECUTION_FAILED"
	ErrDependency ErrorCode = "DEPENDENCY_FAILED"
	ErrInternal   ErrorCode = "INTERNAL_ERROR"

	ErrInvalidConfig ErrorCode = "INVALID_CONFIGURATION"
	ErrExecution     ErrorCode = "EXECUTION_FAILED"
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

func NewFlagReadError(flagName string, err error) error {
	return &CoreError{
		Code:    ErrFlagRead,
		Scope:   flagName,
		Message: "failed to read flag value",
		Err:     err,
	}
}

func NewValidationError(scope, msg string) error {
	return &CoreError{
		Code:    ErrValidation,
		Scope:   scope,
		Message: msg,
	}
}

func NewIOError(scope, msg string, err error) error {
	return &CoreError{
		Code:    ErrIO,
		Scope:   scope,
		Message: msg,
		Err:     err,
	}
}

func NewTemplateError(scope, msg string, err error) error {
	return &CoreError{
		Code:    ErrTemplate,
		Scope:   scope,
		Message: msg,
		Err:     err,
	}
}

func NewRuntimeError(scope, msg string, err error) error {
	return &CoreError{
		Code:    ErrRuntime,
		Scope:   scope,
		Message: msg,
		Err:     err,
	}
}

func NewDependencyError(scope, msg string, err error) error {
	return &CoreError{
		Code:    ErrDependency,
		Scope:   scope,
		Message: msg,
		Err:     err,
	}
}

func NewInternalError(scope, msg string, err error) error {
	return &CoreError{
		Code:    ErrInternal,
		Scope:   scope,
		Message: msg,
		Err:     err,
	}
}

func NewInvalidConfigError(scope, msg string) error {
	return &CoreError{
		Code:    ErrInvalidConfig,
		Scope:   scope,
		Message: msg,
	}
}

func NewExecutionError(scope, msg string, err error) error {
	return &CoreError{
		Code:    ErrExecution,
		Scope:   scope,
		Message: msg,
		Err:     err,
	}
}
