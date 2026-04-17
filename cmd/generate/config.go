package generate

import (
	"fmt"
	"slices"
	"strings"

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
	flags.StringP("format", "f", "toml", "config format")
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

	formats := []string{"json", "toml", "yaml", "yml"}
	if !slices.Contains(formats, format) {
		return appErr.NewValidationError("config", fmt.Sprintf(
			"invalid value '%s' (allowed: %s)",
			format, strings.Join(formats, ", "),
		))
	}

	fileName := fmt.Sprintf("./trax.%s", format)
	if err := writeConfig(fileName); err != nil {
		return appErr.NewIOError("config", "failed to write configuration file", err)
	}

	out.Success("config", fmt.Sprintf("config written %s", fileName))
	return nil
}
