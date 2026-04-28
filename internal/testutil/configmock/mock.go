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
