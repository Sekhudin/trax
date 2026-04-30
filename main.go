package main

import (
	"fmt"
	"os"

	"github.com/sekhudin/trax/cmd"
)

var (
	Exit = func(code int) {
		os.Exit(code)
	}

	Recover = func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "[INTERNAL ERROR]", r)
			Exit(1)
		}
	}
)

func main() {
	defer Recover()

	cmd.Execute()
}
