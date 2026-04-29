package errormock

import (
	appErr "github.com/sekhudin/trax/internal/errors"
)

func With(code appErr.ErrorCode) *appErr.CoreError {
	return &appErr.CoreError{Code: code}
}
