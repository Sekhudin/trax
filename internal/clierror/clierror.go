package clierror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

func printText(w io.Writer, err error) {
	switch e := err.(type) {
	case *appErr.CoreError:
		fmt.Fprintf(w, "\n❌ [%s]", e.Code)

		if e.Scope != "" {
			fmt.Fprintf(w, " (%s)", e.Scope)
		}

		fmt.Fprintf(w, " %s\n", e.Message)

		if e.Err != nil {
			fmt.Fprintf(w, "   ↳ %v\n", e.Err)
		}

	default:
		fmt.Fprintf(w, "❌ %v\n", err)
	}
}

func printJSON(w io.Writer, err error) {
	payload := map[string]any{
		"error": err.Error(),
	}

	if e, ok := err.(*appErr.CoreError); ok {
		payload["code"] = e.Code
		payload["scope"] = e.Scope
		payload["message"] = e.Message

		if e.Err != nil {
			payload["cause"] = e.Err.Error()
		}
	}

	b, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintln(w, "\n"+string(b))
}

func Print(w io.Writer, err error) {
	if err == nil {
		return
	}

	if viper.GetBool("debug") {
		printJSON(w, err)
		return
	}

	printText(w, err)
}

func ExitCode(err error) int {
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
