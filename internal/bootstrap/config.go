package bootstrap

import (
	"errors"
	"strings"

	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

func LoadConfig(cfgFile string) error {
	viper.SetDefault("debug", false)

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
