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

type generateconfig struct{}

var (
	gc        = generateconfig{}
	gcCommand = docs.ApplyDocs(&doc.config, &cobra.Command{
		RunE: gc.runE,
	})
)

func init() {
	flags := gcCommand.Flags()

	flags.Bool("override", false, "overwrite existing config file")
	flags.StringP("format", "f", "toml", "config format")
}

func (*generateconfig) runE(cmd *cobra.Command, args []string) error {
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
			"invalid value %q (allowed: %q)",
			format, strings.Join(formats, " | "),
		))
	}

	fileName := fmt.Sprintf("./trax.%s", format)
	if err := writeConfig(fileName); err != nil {
		return err
	}

	out.Success("config", fmt.Sprintf("config written %s", fileName))
	return nil
}
