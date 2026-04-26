package clierror

import (
	"errors"

	appErr "github.com/sekhudin/trax/internal/errors"
	"github.com/sekhudin/trax/internal/output"
)

type Handler struct {
	out *output.Context
}

func New(out *output.Context) *Handler {
	return &Handler{out: out}
}

func (h *Handler) Print(err error) {
	if err == nil {
		return
	}

	var ce *appErr.CoreError
	if errors.As(err, &ce) {
		h.printCoreError(ce)
		return
	}

	h.out.Error("runtime", err.Error())
}

func (h *Handler) printCoreError(e *appErr.CoreError) {
	scope := e.Scope
	if scope == "" {
		scope = "core"
	}

	h.out.Error(scope, e.Message)

	if e.Err != nil {
		h.out.Cause("cause", e.Err.Error())
	}
}

func (h *Handler) ExitCode(err error) int {
	var ce *appErr.CoreError
	if errors.As(err, &ce) {
		switch ce.Code {
		case appErr.ErrValidation:
			return 2
		case appErr.ErrConfigNotFound:
			return 3
		case appErr.ErrConfigLoad:
			return 4
		default:
			return 1
		}
	}
	return 1
}
