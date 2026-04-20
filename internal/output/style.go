package output

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	IconSuccess = Green("✔")
	IconInfo    = Blue("ℹ")
	IconWarn    = Yellow("⚠")
	IconError   = Red("✖")
	IconDetail  = "↳"
)

func color(code string, v ...any) string {
	s := fmt.Sprint(v...)

	noColor := viper.GetBool("no-color")
	debug := viper.GetBool("debug")

	if noColor || debug {
		return s
	}

	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", code, s)
}

func Red(v ...any) string {
	return color("31", v...)
}

func Yellow(v ...any) string {
	return color("33", v...)
}

func Green(v ...any) string {
	return color("32", v...)
}

func Blue(v ...any) string {
	return color("34", v...)
}

func Cyan(v ...any) string {
	return color("36", v...)
}

func Gray(v ...any) string {
	return color("90", v...)
}

func Bold(v ...any) string {
	return color("1", v...)
}
