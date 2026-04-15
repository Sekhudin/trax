package bootstrap

import (
	"errors"
	"strings"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

func LoadConfig() error {
	cfgFile := viper.GetString("config")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".trax.config")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("TRAX")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {

		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return appErr.NewConfigNotFoundError(
				"config",
				"config file not found, using defaults",
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
