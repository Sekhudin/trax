package show

import (
	"encoding/json"
	"fmt"
	"io"

	"trax/internal/docs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	appErr "trax/internal/errors"
)

var sConfigCmd = docs.ApplyDocs(sConfigDocs, &cobra.Command{
	RunE: sConfigRunE,
})

func init() {
	flags := sConfigCmd.Flags()

	flags.Bool("json", false, "Output as json")
}

func sConfigRunE(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	asJSON, err := flags.GetBool("json")
	if err != nil {
		return appErr.NewFlagReadError("json", err)
	}

	settings := viper.AllSettings()
	w := cmd.OutOrStdout()

	if asJSON {
		return printJSON(w, settings)
	}

	return printFlat(w, "", settings)
}

func printJSON(w io.Writer, data map[string]any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return appErr.NewInternalError("config", "failed to marshal config to json", err)
	}

	if _, err := fmt.Fprintln(w, string(b)); err != nil {
		return appErr.NewIOError("stdout", "failed to write json output", err)
	}

	return nil
}

func printFlat(w io.Writer, prefix string, data map[string]any) error {
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]any:
			if err := printFlat(w, key, val); err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprintf(w, "%s = %v\n", key, val); err != nil {
				return appErr.NewIOError("stdout", "failed to write config output", err)
			}
		}
	}

	return nil
}
