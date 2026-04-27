package config

import "github.com/spf13/viper"

type RoutesSymbols struct {
	Param    string
	Wildcard string
	Root     string
}

type RoutesConfig struct {
	Strategy string
	Root     string
	NoDeps   bool
	File     string
	Output   string
	Prefix   string
	Symbols  *RoutesSymbols
}

func (*config) Routes() *RoutesConfig {
	return &RoutesConfig{
		Strategy: viper.GetString("routes.strategy"),
		Root:     viper.GetString("routes.root"),
		NoDeps:   viper.GetBool("routes.no-deps"),
		File:     viper.GetString("routes.file"),
		Output:   viper.GetString("routes.output"),
		Prefix:   viper.GetString("routes.prefix"),
		Symbols: &RoutesSymbols{
			Param:    viper.GetString("routes.symbols.param"),
			Wildcard: viper.GetString("routes.symbols.wildcard"),
			Root:     viper.GetString("routes.symbols.root"),
		},
	}
}
