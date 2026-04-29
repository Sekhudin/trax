package configmock

import "github.com/sekhudin/trax/internal/config"

type Config struct {
	RoutesCalled bool
	RoutesFn     func() *config.RoutesConfig
}

func (c *Config) Reset() {
	c.RoutesCalled = false

	c.RoutesFn = func() *config.RoutesConfig {
		return &config.RoutesConfig{
			Symbols: &config.RoutesSymbols{},
		}
	}
}

func (c *Config) Routes() *config.RoutesConfig {
	c.RoutesCalled = true
	if c.RoutesFn != nil {
		return c.RoutesFn()
	}

	return &config.RoutesConfig{
		Symbols: &config.RoutesSymbols{},
	}
}

type Writer struct {
	WriteCalled bool
	FileCalled  bool

	WriteFn func() error
	FileFn  func() string
}

func (w *Writer) Reset() {
	w.WriteCalled = false
	w.FileCalled = false

	w.WriteFn = func() error {
		return nil
	}

	w.FileFn = func() string {
		return ""
	}
}

func (w *Writer) Write() error {
	w.WriteCalled = true
	if w.WriteFn != nil {
		return w.WriteFn()
	}
	return nil
}

func (w *Writer) File() string {
	w.FileCalled = true
	if w.FileFn != nil {
		return w.FileFn()
	}
	return ""
}
