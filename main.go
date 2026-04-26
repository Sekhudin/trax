package main

import (
	"fmt"
	"os"

	"github.com/sekhudin/trax/cmd"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "[BUG]", r)
			os.Exit(1)
		}
	}()

	cmd.Execute()
}
