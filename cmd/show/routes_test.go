package show

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/config"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/testutil/appmock"
	"github.com/sekhudin/trax/internal/testutil/cobramock"
	"github.com/sekhudin/trax/internal/testutil/configmock"
	"github.com/sekhudin/trax/internal/testutil/mock"
	"github.com/sekhudin/trax/internal/testutil/routesmock"
	"github.com/sekhudin/trax/modules/routes"
	"github.com/spf13/cobra"
)

type mockRoutesCtx struct {
	PreRunECalled bool
	RunECalled    bool
}

func (m *mockRoutesCtx) Reset() {
	m.PreRunECalled = false
	m.RunECalled = false
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

func TestRoutes_Success(t *testing.T) {
	ctx := appmock.NewContext()
	mockConfig := configmock.Config{}
	mockRoutesConfig := routesmock.RoutesConfig{}
	mockBuilder := routesmock.Builder{}

	t.Run("new_ctx_initialization", func(t *testing.T) {
		mock.Reset(ctx)

		c := NewRoutesCtx(ctx)
		rctx := c.(*routesctx)

		if rctx.config() == nil || rctx.routeConfig(&config.RoutesConfig{}) == nil || rctx.routeBuilder(&routes.Config{}) == nil {
			t.Fatal("fail")
		}
	})

	t.Run("cmd_flags_exist", func(t *testing.T) {
		mock.Reset(ctx)

		cmd := NewRoutesCmd(&doc.Docs{}, NewRoutesCtx(ctx))

		flags := []string{"strategy", "root", "file", "key", "json"}
		for _, f := range flags {
			if cmd.Flags().Lookup(f) == nil {
				t.Fatalf("missing %s", f)
			}
		}
	})

	t.Run("cmd_execution_flow", func(t *testing.T) {
		m := &mockRoutesCtx{}
		cmd := NewRoutesCmd(&doc.Docs{}, m)

		if err := cmd.Execute(); err != nil || !m.PreRunECalled || !m.RunECalled {
			t.Fatal("fail")
		}
	})

	t.Run("prerun_load_config", func(t *testing.T) {
		mock.Reset(ctx, &mockConfig, &mockRoutesConfig)

		c := &routesctx{
			ctx: ctx,
			config: func() config.Config {
				return &mockConfig
			},
			routeConfig: func(*config.RoutesConfig) routes.RoutesConfig {
				return &mockRoutesConfig
			},
		}

		if err := c.PreRunE(newTestRoutesCmd(c)); err != nil || c.cfg == nil || !ctx.OutMock.InfoCalled {
			t.Fatal("fail")
		}
	})

	t.Run("rune_flat_output", func(t *testing.T) {
		mock.Reset(ctx, &mockBuilder)

		c := &routesctx{
			ctx: ctx,
			cfg: &routes.Config{},
			routeBuilder: func(*routes.Config) routes.Builder {
				return &mockBuilder
			},
		}

		cmd := newTestRoutesCmd(c)
		_ = cmd.Flags().Set("json", "false")
		if err := c.RunE(cmd); err != nil || !ctx.OutMock.AsFlatCalled {
			t.Fatal("fail")
		}
	})

	t.Run("rune_json_output", func(t *testing.T) {
		mock.Reset(ctx, &mockBuilder)

		c := &routesctx{
			ctx: ctx,
			cfg: &routes.Config{},
			routeBuilder: func(*routes.Config) routes.Builder {
				return &mockBuilder
			},
		}

		cmd := newTestRoutesCmd(c)
		_ = cmd.Flags().Set("json", "true")
		if err := c.RunE(cmd); err != nil || !ctx.OutMock.AsJsonCalled {
			t.Fatal("fail")
		}
	})
}

func TestRoutes_Error(t *testing.T) {
	ctx := appmock.NewContext()
	mockRoutesConfig := routesmock.RoutesConfig{}
	mockBuilder := routesmock.Builder{}

	t.Run("prerun_load_failure", func(t *testing.T) {
		mock.Reset(ctx, &mockBuilder)

		mockRoutesConfig.LoadFn = func() (*routes.Config, error) {
			return nil, errors.New("boom")
		}

		c := &routesctx{
			ctx:    ctx,
			config: func() config.Config { return &configmock.Config{} },
			routeConfig: func(rc *config.RoutesConfig) routes.RoutesConfig {
				return &mockRoutesConfig
			},
		}
		if err := c.PreRunE(newTestRoutesCmd(c)); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_key_flag", func(t *testing.T) {
		mock.Reset(ctx)

		c := &routesctx{ctx: ctx, cfg: &routes.Config{}}
		cmd := newTestRoutesCmd(c)

		cmd.Flags().Lookup("key").Value = &cobramock.FlagBroken{}
		if err := c.RunE(cmd); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_json_flag", func(t *testing.T) {
		mock.Reset(ctx)

		c := &routesctx{ctx: ctx, cfg: &routes.Config{}}
		cmd := newTestRoutesCmd(c)

		cmd.Flags().Lookup("json").Value = &cobramock.FlagBroken{}
		if err := c.RunE(cmd); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_build_failure", func(t *testing.T) {
		mock.Reset(ctx, &mockBuilder)

		mockBuilder.BuildFn = func() (routes.BuildResult, error) {
			return nil, errors.New("err")
		}

		c := &routesctx{
			ctx: ctx,
			cfg: &routes.Config{},
			routeBuilder: func(*routes.Config) routes.Builder {
				return &mockBuilder
			},
		}
		if err := c.RunE(newTestRoutesCmd(c)); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_selector_failure", func(t *testing.T) {
		mock.Reset(ctx, &mockBuilder)

		mockBuilder.BuildFn = func() (routes.BuildResult, error) {
			return &routesmock.BuildResult{
				SelectFn: func(key string) (map[string]any, error) { return nil, errors.New("err") },
			}, nil
		}

		c := &routesctx{
			ctx: ctx,
			cfg: &routes.Config{},
			routeBuilder: func(*routes.Config) routes.Builder {
				return &mockBuilder
			},
		}
		if err := c.RunE(newTestRoutesCmd(c)); err == nil {
			t.Fatal("fail")
		}
	})
}

func TestRoutes_Fallback(t *testing.T) {
	t.Run("new_ctx_smoke", func(t *testing.T) {
		if c := NewRoutesCtx(appmock.NewContext()); c == nil {
			t.Fatal("fail")
		}
	})
}
