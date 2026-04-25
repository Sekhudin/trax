package output

import "fmt"

type Colorizer struct {
	NoColor bool
}

func NewColorizer(noColor bool) Colorizer {
	return Colorizer{NoColor: noColor}
}

func (c Colorizer) color(code string, v ...any) string {
	s := fmt.Sprint(v...)
	if c.NoColor {
		return s
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", code, s)
}

func (c Colorizer) Red(v ...any) string {
	return c.color("31", v...)
}

func (c Colorizer) Yellow(v ...any) string {
	return c.color("33", v...)
}

func (c Colorizer) Green(v ...any) string {
	return c.color("32", v...)
}

func (c Colorizer) Blue(v ...any) string {
	return c.color("34", v...)
}

func (c Colorizer) Cyan(v ...any) string {
	return c.color("36", v...)
}

func (c Colorizer) Gray(v ...any) string {
	return c.color("90", v...)
}

func (c Colorizer) Bold(v ...any) string {
	return c.color("1", v...)
}
