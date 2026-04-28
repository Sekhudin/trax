package show

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/testutil/appmock"
	"github.com/spf13/cobra"
)

type mockConfigCtx struct {
	PreRunECalled bool
	RunECalled    bool
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
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error { return c.PreRunE() }
	cmd.RunE = func(cmd *cobra.Command, args []string) error { return c.RunE(cmd) }
	return cmd
}

func Test_NewConfigCmd_flags_exist(t *testing.T) {
	ctx := &mockConfigCtx{}

	cmd := NewConfigCmd(&doc.Docs{}, ctx)

	if cmd.Flags().Lookup("json") == nil {
		t.Fatal("json flag missing")
	}
}

func Test_NewConfigCmd_execute_full_path(t *testing.T) {
	ctx := &mockConfigCtx{}

	cmd := NewConfigCmd(&doc.Docs{}, ctx)

	cmd.SetArgs([]string{"--json"})

	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	if !ctx.PreRunECalled {
		t.Fatal("PreRunE not called")
	}

	if !ctx.RunECalled {
		t.Fatal("RunE not called")
	}
}

func Test_NewConfigCmd_direct_closure(t *testing.T) {
	ctx := &mockConfigCtx{}

	cmd := NewConfigCmd(&doc.Docs{}, ctx)

	if err := cmd.PreRunE(cmd, []string{}); err != nil {
		t.Fatal(err)
	}

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatal(err)
	}
}

func Test_PreRunE(t *testing.T) {
	ctx := appmock.NewContext()
	c := NewConfigCtx(ctx)

	err := c.PreRunE()
	if err != nil {
		t.Fatal(err)
	}
	if !ctx.OutMock.InfoCalled {
		t.Fatal("info not called")
	}
}

func Test_RunE_json_true(t *testing.T) {
	ctx := appmock.NewContext()
	c := NewConfigCtx(ctx)
	cmd := newTestConfigCmd(c)
	_ = cmd.Flags().Set("json", "true")

	err := c.RunE(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if !ctx.OutMock.AsJsonCalled || ctx.OutMock.AsFlatCalled {
		t.Fatal("wrong path")
	}
}

func Test_RunE_json_false(t *testing.T) {
	ctx := appmock.NewContext()
	c := NewConfigCtx(ctx)
	cmd := newTestConfigCmd(c)
	_ = cmd.Flags().Set("json", "false")

	err := c.RunE(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if !ctx.OutMock.AsFlatCalled || ctx.OutMock.AsJsonCalled {
		t.Fatal("wrong path")
	}
}

func Test_RunE_flag_error(t *testing.T) {
	ctx := appmock.NewContext()
	c := NewConfigCtx(ctx)
	cmd := &cobra.Command{}

	err := c.RunE(cmd)
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_RunE_json_error(t *testing.T) {
	ctx := appmock.NewContext()
	ctx.OutMock.JSONErr = errors.New("err")

	c := NewConfigCtx(ctx)
	cmd := newTestConfigCmd(c)
	_ = cmd.Flags().Set("json", "true")

	err := c.RunE(cmd)
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_RunE_flat_error(t *testing.T) {
	ctx := appmock.NewContext()
	ctx.OutMock.FlatErr = errors.New("err")

	c := NewConfigCtx(ctx)
	cmd := newTestConfigCmd(c)
	_ = cmd.Flags().Set("json", "false")

	err := c.RunE(cmd)
	if err == nil {
		t.Fatal("expected error")
	}
}
