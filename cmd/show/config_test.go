package show

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/testutil/appmock"
	"github.com/sekhudin/trax/internal/testutil/mock"
	"github.com/spf13/cobra"
)

type mockConfigCtx struct {
	PreRunECalled bool
	RunECalled    bool
}

func (m *mockConfigCtx) Reset() {
	m.PreRunECalled = false
	m.RunECalled = false
}

func (m *mockConfigCtx) PreRunE() error {
	m.PreRunECalled = true
	return nil
}

func (m *mockConfigCtx) RunE(cmd *cobra.Command) error {
	m.RunECalled = true
	return nil
}

func newTestConfigCmd(c ConfigCtx) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("json", false, "")
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return c.PreRunE()
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.RunE(cmd)
	}

	return cmd
}

func TestConfig_Success(t *testing.T) {
	ctx := appmock.NewContext()
	mCtx := mockConfigCtx{}

	t.Run("cmd_flags_exist", func(t *testing.T) {
		mock.Reset(&mCtx)

		cmd := NewConfigCmd(&doc.Docs{}, &mCtx)

		if cmd.Flags().Lookup("json") == nil {
			t.Fatal("fail")
		}
	})

	t.Run("cmd_execution_flow", func(t *testing.T) {
		mock.Reset(&mCtx)

		cmd := NewConfigCmd(&doc.Docs{}, &mCtx)
		cmd.SetArgs([]string{"--json"})

		if err := cmd.Execute(); err != nil || !mCtx.PreRunECalled || !mCtx.RunECalled {
			t.Fatal("fail")
		}
	})

	t.Run("cmd_closure_direct", func(t *testing.T) {
		mock.Reset(&mCtx)

		cmd := NewConfigCmd(&doc.Docs{}, &mCtx)

		if err := cmd.PreRunE(cmd, []string{}); err != nil {
			t.Fatal("fail")
		}
		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Fatal("fail")
		}
	})

	t.Run("prerun_info_output", func(t *testing.T) {
		mock.Reset(ctx)

		c := NewConfigCtx(ctx)

		if err := c.PreRunE(); err != nil || !ctx.OutMock.InfoCalled {
			t.Fatal("fail")
		}
	})

	t.Run("rune_json_output", func(t *testing.T) {
		mock.Reset(ctx)

		c := NewConfigCtx(ctx)
		cmd := newTestConfigCmd(c)

		_ = cmd.Flags().Set("json", "true")
		if err := c.RunE(cmd); err != nil || !ctx.OutMock.AsJsonCalled {
			t.Fatal("fail")
		}
	})

	t.Run("rune_flat_output", func(t *testing.T) {
		mock.Reset(ctx)

		c := NewConfigCtx(ctx)
		cmd := newTestConfigCmd(c)

		_ = cmd.Flags().Set("json", "false")
		if err := c.RunE(cmd); err != nil || !ctx.OutMock.AsFlatCalled {
			t.Fatal("fail")
		}
	})
}

func TestConfig_Error(t *testing.T) {
	ctx := appmock.NewContext()

	t.Run("rune_missing_flags", func(t *testing.T) {
		mock.Reset(ctx)

		c := NewConfigCtx(ctx)

		if err := c.RunE(&cobra.Command{}); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_json_failure", func(t *testing.T) {
		mock.Reset(ctx)

		ctx.OutMock.JSONErr = errors.New("err")

		c := NewConfigCtx(ctx)
		cmd := newTestConfigCmd(c)

		_ = cmd.Flags().Set("json", "true")
		if err := c.RunE(cmd); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_flat_failure", func(t *testing.T) {
		ctx.OutMock.FlatErr = errors.New("err")

		c := NewConfigCtx(ctx)
		cmd := newTestConfigCmd(c)

		_ = cmd.Flags().Set("json", "false")
		if err := c.RunE(cmd); err == nil {
			t.Fatal("fail")
		}
	})
}

func TestConfig_Fallback(t *testing.T) {
	t.Run("new_ctx_initialization", func(t *testing.T) {
		if c := NewConfigCtx(appmock.NewContext()); c == nil {
			t.Fatal("fail")
		}
	})
}
