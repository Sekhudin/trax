package configmock

import "github.com/sekhudin/trax/internal/config"

type Config struct {
	RoutesCalled bool
	RoutesFn     func() *config.RoutesConfig
}

func (c *Config) Routes() *config.RoutesConfig {
	c.RoutesCalled = true
	if c.RoutesFn != nil {
		return c.RoutesFn()
	}
	return &config.RoutesConfig{}
}

type Writer struct {
	WriteCalled bool
	FileCalled  bool

	WriteFn func() error
	FileFn  func() string
}

func (c *Writer) Write() error {
	c.WriteCalled = true
	if c.WriteFn != nil {
		return c.WriteFn()
	}
	return nil
}

func (c *Writer) File() string {
	c.FileCalled = true
	if c.FileFn != nil {
		return c.FileFn()
	}
	return ""
}
