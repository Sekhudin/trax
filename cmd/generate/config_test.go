package generate

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/testutil/appmock"
	"github.com/sekhudin/trax/internal/testutil/cobramock"
	"github.com/sekhudin/trax/internal/testutil/configmock"
	"github.com/sekhudin/trax/internal/testutil/mock"
	"github.com/spf13/cobra"
)

type mockConfigCtx struct {
	RunECalled bool
}

func (m *mockConfigCtx) Reset() {
	m.RunECalled = false
}

func (m *mockConfigCtx) RunE(cmd *cobra.Command) error {
	m.RunECalled = true
	return nil
}

func newTestConfigCmd(c ConfigCtx) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("override", false, "")
	cmd.Flags().String("format", "toml", "")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.RunE(cmd)
	}

	return cmd
}

func TestGenConfig_Success(t *testing.T) {
	ctx := appmock.NewContext()
	mockWriter := configmock.Writer{}

	t.Run("new_ctx_initialization", func(t *testing.T) {
		mock.Reset(ctx)

		c := NewConfigCtx(ctx)
		rctx := c.(*configctx)

		if rctx.cfgFilePath("toml") == "" || rctx.cfgWriter("toml", true) == nil {
			t.Fatal("fail")
		}
	})

	t.Run("cmd_flags_exist", func(t *testing.T) {
		mock.Reset(ctx)

		cmd := NewConfigCmd(&doc.Docs{}, NewConfigCtx(ctx))
		if cmd.Flags().Lookup("override") == nil || cmd.Flags().Lookup("format") == nil {
			t.Fatal("fail")
		}
	})

	t.Run("cmd_execution_flow", func(t *testing.T) {
		m := &mockConfigCtx{}

		cmd := NewConfigCmd(&doc.Docs{}, m)
		if err := cmd.Execute(); err != nil || !m.RunECalled {
			t.Fatal("fail")
		}
	})

	t.Run("rune_write_config", func(t *testing.T) {
		mock.Reset(ctx, &mockWriter)

		c := &configctx{
			ctx:        ctx,
			cfgFormats: map[string]struct{}{"toml": {}},
			cfgFilePath: func(s string) string {
				return s
			},
			cfgWriter: func(s string, b bool) config.Writer {
				return &mockWriter
			},
		}

		if err := c.RunE(newTestConfigCmd(c)); err != nil || !mockWriter.WriteCalled || !ctx.OutMock.SuccessCalled {
			t.Fatal("fail")
		}
	})
}

func TestGenConfig_Error(t *testing.T) {
	ctx := appmock.NewContext()
	mockWriter := configmock.Writer{}

	t.Run("rune_format_flag", func(t *testing.T) {
		mock.Reset(ctx)

		c := &configctx{ctx: ctx, cfgFormats: map[string]struct{}{
			"toml": {},
		}}

		cmd := newTestConfigCmd(c)

		cmd.Flags().Lookup("format").Value = &cobramock.FlagBroken{}
		if err := c.RunE(cmd); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_override_flag", func(t *testing.T) {
		mock.Reset(ctx)

		c := &configctx{ctx: ctx, cfgFormats: map[string]struct{}{
			"toml": {},
		}}

		cmd := newTestConfigCmd(c)

		cmd.Flags().Lookup("override").Value = &cobramock.FlagBroken{}
		if err := c.RunE(cmd); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_invalid_format", func(t *testing.T) {
		mock.Reset(ctx, &mockWriter)

		c := &configctx{
			ctx: ctx,
			cfgFormats: map[string]struct{}{
				"invalid": {},
			},

			cfgWriter: func(s string, b bool) config.Writer {
				return &mockWriter
			},
		}

		if err := c.RunE(newTestConfigCmd(c)); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_write_failure", func(t *testing.T) {
		mock.Reset(ctx, &mockWriter)

		mockWriter.WriteFn = func() error { return errors.New("err") }

		c := &configctx{
			ctx: ctx,
			cfgFormats: map[string]struct{}{
				"toml": {},
			},
			cfgFilePath: func(s string) string {
				return s
			},
			cfgWriter: func(s string, b bool) config.Writer {
				return &mockWriter
			},
		}

		if err := c.RunE(newTestConfigCmd(c)); err == nil || !mockWriter.WriteCalled {
			t.Fatal("fail")
		}
	})
}

func TestGenConfig_Fallback(t *testing.T) {
	t.Run("verify_default_formats", func(t *testing.T) {
		c := NewConfigCtx(appmock.NewContext()).(*configctx)

		formats := []string{"json", "toml", "yaml", "yml"}
		for _, f := range formats {
			if _, ok := c.cfgFormats[f]; !ok {
				t.Fatalf("missing %s", f)
			}
		}
	})
}
