package show

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/testutil/appmock"
	"github.com/sekhudin/trax/internal/testutil/cobramock"
	"github.com/sekhudin/trax/internal/testutil/configmock"
	"github.com/sekhudin/trax/internal/testutil/routesmock"
	"github.com/sekhudin/trax/modules/routes"
	"github.com/spf13/cobra"
)

type mockRoutesCtx struct {
	PreRunECalled bool
	RunECalled    bool
}

func (m *mockRoutesCtx) PreRunE(cmd *cobra.Command) error {
	m.PreRunECalled = true
	return nil
}

func (m *mockRoutesCtx) RunE(cmd *cobra.Command) error {
	m.RunECalled = true
	return nil
}

func newTestRoutesCmd(c RoutesCtx) *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Flags().String("strategy", "", "")
	cmd.Flags().String("root", "", "")
	cmd.Flags().String("file", "", "")
	cmd.Flags().String("key", "", "")
	cmd.Flags().Bool("json", false, "")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return c.PreRunE(cmd)
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.RunE(cmd)
	}

	return cmd
}

func Test_NewRoutesCtx_executes_internal_functions(t *testing.T) {
	ctx := appmock.NewContext()

	c := NewRoutesCtx(ctx)

	conf := c.(*routesctx).config()
	if conf == nil {
		t.Fatal("config nil")
	}

	rc := &config.RoutesConfig{}
	routesCfg := c.(*routesctx).routeConfig(rc)
	if routesCfg == nil {
		t.Fatal("routeConfig nil")
	}

	builder := c.(*routesctx).routeBuilder(&routes.Config{})
	if builder == nil {
		t.Fatal("routeBuilder nil")
	}
}

func Test_NewRoutesCmd_flags_exist(t *testing.T) {
	ctx := appmock.NewContext()
	c := NewRoutesCtx(ctx)

	cmd := NewRoutesCmd(&doc.Docs{}, c)

	flags := []string{"strategy", "root", "file", "key", "json"}

	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Fatalf("%s missing", f)
		}
	}
}

func Test_NewRoutesCmd_execute_full_path(t *testing.T) {
	ctx := &mockRoutesCtx{}
	cmd := NewRoutesCmd(&doc.Docs{}, ctx)

	cmd.SetArgs([]string{})
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

func Test_PreRunE_success(t *testing.T) {
	ctx := appmock.NewContext()

	c := &routesctx{
		ctx: ctx,
		config: func() config.Config {
			return &configmock.Config{}
		},
		routeConfig: func(*config.RoutesConfig) routes.RoutesConfig {
			return &routesmock.RoutesConfig{}
		},
	}

	cmd := newTestRoutesCmd(c)

	err := c.PreRunE(cmd)
	if err != nil {
		t.Fatal(err)
	}

	if c.cfg == nil {
		t.Fatal("cfg not set")
	}

	if !ctx.OutMock.InfoCalled {
		t.Fatal("info not called")
	}
}

func Test_PreRunE_load_error(t *testing.T) {
	ctx := appmock.NewContext()

	c := &routesctx{
		ctx: ctx,
		config: func() config.Config {
			return &configmock.Config{}
		},
		routeConfig: func(*config.RoutesConfig) routes.RoutesConfig {
			return &routesmock.RoutesConfig{
				LoadFn: func() (*routes.Config, error) {
					return nil, errors.New("boom")
				},
			}
		},
	}

	err := c.PreRunE(newTestRoutesCmd(c))
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_RunE_flat_success(t *testing.T) {
	ctx := appmock.NewContext()

	c := &routesctx{
		ctx: ctx,
		cfg: &routes.Config{},
		routeBuilder: func(*routes.Config) routes.Builder {
			return &routesmock.Builder{}
		},
	}

	cmd := newTestRoutesCmd(c)
	_ = cmd.Flags().Set("json", "false")

	err := c.RunE(cmd)
	if err != nil {
		t.Fatal(err)
	}

	if !ctx.OutMock.AsFlatCalled {
		t.Fatal("flat not called")
	}
}

func Test_RunE_json_success(t *testing.T) {
	ctx := appmock.NewContext()

	c := &routesctx{
		ctx: ctx,
		cfg: &routes.Config{},
		routeBuilder: func(*routes.Config) routes.Builder {
			return &routesmock.Builder{}
		},
	}

	cmd := newTestRoutesCmd(c)
	_ = cmd.Flags().Set("json", "true")

	err := c.RunE(cmd)
	if err != nil {
		t.Fatal(err)
	}

	if !ctx.OutMock.AsJsonCalled {
		t.Fatal("json not called")
	}
}

func Test_RunE_key_flag_error(t *testing.T) {
	c := &routesctx{
		ctx: appmock.NewContext(),
		cfg: &routes.Config{},
		routeBuilder: func(*routes.Config) routes.Builder {
			t.Fatal("should not reach builder")
			return nil
		},
	}

	cmd := newTestRoutesCmd(c)
	cmd.Flags().Lookup("key").Value = &cobramock.FlagBroken{}

	err := c.RunE(cmd)
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_RunE_json_flag_error(t *testing.T) {
	c := &routesctx{
		ctx: &appmock.Context{},
		cfg: &routes.Config{},
		routeBuilder: func(*routes.Config) routes.Builder {
			t.Fatal("should not reach builder")
			return nil
		},
	}

	cmd := newTestRoutesCmd(c)
	cmd.Flags().Lookup("json").Value = &cobramock.FlagBroken{}

	err := c.RunE(cmd)
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_RunE_build_error(t *testing.T) {
	ctx := appmock.NewContext()

	c := &routesctx{
		ctx: ctx,
		cfg: &routes.Config{},
		routeBuilder: func(*routes.Config) routes.Builder {
			return &routesmock.Builder{
				BuildFn: func() (routes.BuildResult, error) {
					return nil, errors.New("build error")
				},
			}
		},
	}

	err := c.RunE(newTestRoutesCmd(c))
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_RunE_selector_error(t *testing.T) {
	ctx := appmock.NewContext()

	c := &routesctx{
		ctx: ctx,
		cfg: &routes.Config{},
		routeBuilder: func(*routes.Config) routes.Builder {
			return &routesmock.Builder{
				BuildFn: func() (routes.BuildResult, error) {
					return &routesmock.BuildResult{
						SelectFn: func(key string) (map[string]any, error) {
							return nil, errors.New("select error")
						},
					}, nil
				},
			}
		},
	}

	err := c.RunE(newTestRoutesCmd(c))
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_NewRoutesCtx_smoke(t *testing.T) {
	ctx := appmock.NewContext()

	c := NewRoutesCtx(ctx)
	if c == nil {
		t.Fatal("nil ctx")
	}
}
