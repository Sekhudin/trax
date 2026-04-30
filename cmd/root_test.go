package cmd

import (
	"errors"
	"io"
	"testing"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/clierror"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/testutil/appmock"
	"github.com/sekhudin/trax/internal/testutil/cobramock"
	"github.com/sekhudin/trax/internal/testutil/mock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// --- Group: New & Dependencies ---

func TestNew_Success(t *testing.T) {
	ctx := appmock.NewContext()

	t.Run("default_factory_init", func(t *testing.T) {
		cmd := New(ctx)
		if cmd == nil || cmd.Use != "trax" {
			t.Fatal("fail")
		}
	})
}

func TestNewWithDeps_Success(t *testing.T) {
	ctx := appmock.NewContext()
	deps := DefaultDependencies(ctx)

	t.Run("metadata_injection", func(t *testing.T) {
		mock.Reset(ctx)
		d := &doc.Docs{Use: "custom-trax"}
		deps.Docs.Root = *d

		cmd := NewWithDependencies(ctx, deps)
		if cmd.Use != "custom-trax" {
			t.Fatal("fail")
		}
	})

	t.Run("subcommand_injection", func(t *testing.T) {
		deps.NewGenerateCmd = func(a app.Context) *cobra.Command {
			return &cobra.Command{Use: "mock-gen"}
		}
		cmd := NewWithDependencies(ctx, deps)

		found := false
		for _, c := range cmd.Commands() {
			if c.Name() == "mock-gen" {
				found = true
			}
		}
		if !found {
			t.Fatal("missing mock-gen")
		}
	})

	t.Run("trigger_prerun_closure", func(t *testing.T) {
		viper.Reset()
		deps.NewShowCmd = func(a app.Context) *cobra.Command {
			return &cobra.Command{
				Use: "test-run",
				Run: func(cmd *cobra.Command, args []string) {},
			}
		}
		cmd := NewWithDependencies(ctx, deps)
		cmd.SetArgs([]string{"test-run"})
		cmd.SetOut(io.Discard)

		if err := cmd.Execute(); err != nil {
			t.Fatalf("fail: %v", err)
		}
	})
}

// --- Group: Context Logic (PreRun & Flags) ---

func TestContext_Success(t *testing.T) {
	ctx := appmock.NewContext()
	deps := DefaultDependencies(ctx)

	t.Run("prerun_viper_binding", func(t *testing.T) {
		mock.Reset(ctx)
		viper.Reset()
		cmd := NewWithDependencies(ctx, deps)

		cmd.SetArgs([]string{"--config", "", "--debug"})
		_ = cmd.ParseFlags([]string{"--config", "", "--debug"})

		if err := deps.Ctx.PersistentPreRunE(cmd); err != nil {
			t.Fatal(err)
		}
		if !viper.GetBool("debug") {
			t.Fatal("viper binding failed")
		}
	})
}

func TestContext_Error(t *testing.T) {
	ctx := appmock.NewContext()
	deps := DefaultDependencies(ctx)

	t.Run("flag_error_direct", func(t *testing.T) {
		err := deps.Ctx.FlagErrorFn(&cobra.Command{}, errors.New("boom"))
		if err == nil {
			t.Fatal("should return validation error")
		}
	})

	t.Run("flag_read_failure", func(t *testing.T) {
		cmd := NewWithDependencies(ctx, deps)
		cmd.PersistentFlags().Lookup("config").Value = &cobramock.FlagBroken{}
		if err := deps.Ctx.PersistentPreRunE(cmd); err == nil {
			t.Fatal("should error on broken flag")
		}
	})

	t.Run("config_load_failure", func(t *testing.T) {
		cmd := NewWithDependencies(ctx, deps)
		cmd.SetArgs([]string{"--config", "non-existent.yaml"})
		_ = cmd.ParseFlags([]string{"--config", "non-existent.yaml"})

		if err := deps.Ctx.PersistentPreRunE(cmd); err == nil {
			t.Fatal("should error on missing file")
		}
	})
}

// --- Group: Execute & Global Handler ---

func TestExecute_Fallback(t *testing.T) {
	ctx := appmock.NewContext()

	t.Run("execute_success_path", func(t *testing.T) {
		oldCmd := Command
		defer func() { Command = oldCmd }()

		Command = func() *cobra.Command {
			c := New(ctx)
			c.SetArgs([]string{"--help"})
			c.SetOut(io.Discard)
			return c
		}
		Execute()
	})

	t.Run("execute_error_path", func(t *testing.T) {
		oldCmd := Command
		oldHandler := ErrorHanler
		defer func() {
			Command = oldCmd
			ErrorHanler = oldHandler
		}()

		Command = func() *cobra.Command {
			c := New(ctx)
			c.SetArgs([]string{"--invalid"})
			c.SetOut(io.Discard)
			c.SetErr(io.Discard)
			return c
		}

		handlerCalled := false
		ErrorHanler = func(err error, h clierror.Handler) {
			handlerCalled = true
		}

		Execute()
		if !handlerCalled {
			t.Fatal("ErrorHandler should be called")
		}
	})
}
