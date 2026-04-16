package generate

import (
	"fmt"

	"trax/internal/docs"
	"trax/internal/output"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

var gConfigCmd = docs.ApplyDocs(gConfigDocs, &cobra.Command{
	RunE: gConfigRunE,
})

func init() {
	flags := gConfigCmd.Flags()

	flags.Bool("override", false, "overwrite existing config file")
	flags.StringP("format", "f", "toml", "config format (json|yaml|toml`)")
}

func gConfigRunE(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	out := output.New(cmd.OutOrStdout())

	isOverride, err := flags.GetBool("override")
	if err != nil {
		return appErr.NewFlagReadError("override", err)
	}

	writeConfig := viper.SafeWriteConfigAs
	if isOverride {
		writeConfig = viper.WriteConfigAs
	}

	format, err := flags.GetString("format")
	if err != nil {
		return appErr.NewFlagReadError("format", err)
	}

	ext := format
	switch format {
	case "json":
		ext = "json"
	case "yaml":
		ext = "yaml"
	case "toml":
		ext = format
	default:
		return appErr.NewValidationError("config", "Invalid format")
	}

	fileName := fmt.Sprintf("./trax.%s", ext)
	if err := writeConfig(fileName); err != nil {
		return appErr.NewIOError("config", "Failed to write configuration file", err)
	}

	out.Success("config", fmt.Sprintf("Config written %s", fileName))
	return nil
}
