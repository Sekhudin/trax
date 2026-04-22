package bootstrap

import (
	"errors"
	"strings"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
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

	if cfgFile == "" {
		viper.SetConfigName("trax")
		viper.AddConfigPath(".")

		return nil
	}

	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		isFileNotFound := strings.Contains(err.Error(), "no such file")

		if isFileNotFound {
			return appErr.NewConfigNotFoundError(
				"config",
				err.Error(),
			)
		}

		if errors.Is(err, viper.ConfigFileNotFoundError{}) {
			return appErr.NewConfigNotFoundError(
				"config",
				"config file not found",
			)
		}

		return appErr.NewConfigLoadError(
			"config",
			"failed to read config file",
			err,
		)
	}

	return nil
}
