package clierror

import (
	"encoding/json"
	"fmt"
	"os"

	appErr "trax/internal/errors"

	"github.com/spf13/viper"
)

func printText(err error) {
	switch e := err.(type) {
	case *appErr.CoreError:
		fmt.Printf("\n❌ [%s]", e.Code)

		if e.Scope != "" {
			fmt.Printf(" (%s)", e.Scope)
		}

		fmt.Printf(" %s\n", e.Message)

		if e.Err != nil {
			fmt.Printf("   ↳ %v\n", e.Err)
		}

	default:
		fmt.Printf("❌ %v\n", err)
	}
}

func printJSON(err error) {
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
	fmt.Fprintln(os.Stderr, "\n"+string(b))
}

func Print(err error) {
	if err == nil {
		return
	}

	if viper.GetBool("debug") {
		printJSON(err)
		return
	}

	printText(err)
}
