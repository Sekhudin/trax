package config

type Config interface {
	Routes() *RoutesConfig
}

type config struct{}

func New() Config {
	return &config{}
}
