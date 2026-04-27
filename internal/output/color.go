package output

import "fmt"

type Colorizer interface {
	Red(v ...any) string
	Yellow(v ...any) string
	Green(v ...any) string
	Blue(v ...any) string
	Cyan(v ...any) string
	Gray(v ...any) string
	Bold(v ...any) string
}

type colorizer struct {
	noColor bool
}

func NewColorizer(noColor bool) Colorizer {
	return &colorizer{noColor: noColor}
}

func (c colorizer) color(code string, v ...any) string {
	s := fmt.Sprint(v...)
	if c.noColor {
		return s
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", code, s)
}

func (c *colorizer) Red(v ...any) string {
	return c.color("31", v...)
}

func (c *colorizer) Yellow(v ...any) string {
	return c.color("33", v...)
}

func (c *colorizer) Green(v ...any) string {
	return c.color("32", v...)
}

func (c *colorizer) Blue(v ...any) string {
	return c.color("34", v...)
}

func (c *colorizer) Cyan(v ...any) string {
	return c.color("36", v...)
}

func (c *colorizer) Gray(v ...any) string {
	return c.color("90", v...)
}

func (c *colorizer) Bold(v ...any) string {
	return c.color("1", v...)
}
