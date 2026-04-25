package bootstrap

import (
	"strings"

	"github.com/spf13/viper"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type config struct {
	formatter  string
	routes     map[string]any
	formatters map[string]any
}

var cfg = config{
	formatter: "biome",
	routes: map[string]any{
		"strategy":  "next-app",
		"prefix":    "routes",
		"root":      "./src",
		"file":      "",
		"no-deps":   false,
		"output":    "src/trax/routes.ts",
		"formatter": "biome",
		"symbols": map[string]any{
			"param":    "$param",
			"wildcard": "$wildcard",
			"root":     "root",
		},
	},
	formatters: map[string]any{
		"biome": map[string]any{
			"exec": "biome",
			"args": []string{"check", "--write", "src/trax"},
		},
		"prettier": map[string]any{
			"exec": "npx",
			"args": []string{"prettier", "--write", "src/trax"},
		},
	},
}

func LoadConfig(cfgFile string) error {
	viper.SetDefault("formatter", cfg.formatter)
	viper.SetDefault("routes", cfg.routes)
	viper.SetDefault("formatters", cfg.formatters)

	viper.SetEnvPrefix("TRAX")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("trax")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return appErr.NewConfigLoadError(
			"config",
			"failed to read config",
			err,
		)
	}

	return nil
}
