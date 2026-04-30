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
	"github.com/sekhudin/trax/internal/testutil/routesmock"
	"github.com/sekhudin/trax/modules/routes"
	"github.com/spf13/cobra"
)

type mockRoutesCtx struct {
	PreRunECalled  bool
	RunECalled     bool
	PostRunECalled bool
}

func (m *mockRoutesCtx) Reset() {
	m.PreRunECalled = true
	m.RunECalled = true
	m.PostRunECalled = true
}

func (m *mockRoutesCtx) PreRunE(cmd *cobra.Command) error {
	m.PreRunECalled = true
	return nil
}

func (m *mockRoutesCtx) RunE() error {
	m.RunECalled = true
	return nil
}

func (m *mockRoutesCtx) PostRunE(cmd *cobra.Command) error {
	m.PostRunECalled = true
	return nil
}

func newTestRoutesCmd(c RoutesCtx) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().String("strategy", "", "")
	cmd.Flags().String("root", "", "")
	cmd.Flags().String("file", "", "")
	cmd.Flags().String("output", "", "")
	cmd.Flags().String("formatter", "", "")
	cmd.Flags().Bool("no-deps", false, "")
	cmd.Flags().Bool("no-format", false, "")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return c.PreRunE(cmd)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.RunE()
	}

	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return c.PostRunE(cmd)
	}

	return cmd
}

func TestGenRoutes_Success(t *testing.T) {
	ctx := appmock.NewContext()
	mockConfig := configmock.Config{}
	mockRoutesConfig := routesmock.RoutesConfig{}
	mockTemplate := routesmock.TemplateBuilder{}

	t.Run("cmd_initialization_and_flags", func(t *testing.T) {
		mock.Reset(ctx)

		cmd := NewRoutesCmd(&doc.Docs{}, &mockRoutesCtx{})

		flags := []string{"strategy", "root", "file", "output", "formatter", "no-deps", "no-format"}
		for _, f := range flags {
			if cmd.Flags().Lookup(f) == nil {
				t.Fatalf("missing %s", f)
			}
		}
	})

	t.Run("full_execution_flow", func(t *testing.T) {
		m := &mockRoutesCtx{}

		cmd := NewRoutesCmd(&doc.Docs{}, m)
		if err := cmd.Execute(); err != nil || !m.PreRunECalled || !m.RunECalled || !m.PostRunECalled {
			t.Fatal("fail")
		}
	})

	t.Run("prerun_config_loading", func(t *testing.T) {
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

	t.Run("postrun_formatting_enabled", func(t *testing.T) {
		mock.Reset(ctx, &mockRoutesConfig)
		mock.Set(&mockRoutesConfig)

		mockRoutesConfig.Config.Output.Filename = "out.ts"

		c := &routesctx{
			ctx: ctx,
			cfg: mockRoutesConfig.Config,
		}

		cmd := newTestRoutesCmd(c)

		_ = cmd.Flags().Set("no-format", "false")
		if err := c.PostRunE(cmd); err != nil || !ctx.OutMock.SuccessCalled {
			t.Fatal("fail")
		}
	})

	t.Run("execute_factory_functions", func(t *testing.T) {
		mock.Reset(ctx, &mockRoutesConfig)
		mock.Set(&mockRoutesConfig)

		c := NewRoutesCtx(ctx)
		rctx := c.(*routesctx)

		_ = rctx.config()
		_ = rctx.routeConfig(mockConfig.Routes())
		_ = rctx.routeBuilder(mockRoutesConfig.Config)
		_ = rctx.routeTemplate([]routes.Route{}, nil, mockRoutesConfig.Config)

		_ = rctx.routeGenerator(&mockTemplate)
	})
}

func TestGenRoutes_Error(t *testing.T) {
	ctx := appmock.NewContext()
	mockConfig := configmock.Config{}
	mockRoutesConfig := routesmock.RoutesConfig{}
	mockBuilder := routesmock.Builder{}
	mockTemplate := routesmock.TemplateBuilder{}
	mockGenerator := routesmock.Generator{}

	t.Run("prerun_load_failure", func(t *testing.T) {
		mock.Reset(ctx, &mockConfig, &mockRoutesConfig)

		mockRoutesConfig.LoadFn = func() (*routes.Config, error) {
			return nil, errors.New("erro")
		}

		c := &routesctx{
			ctx: ctx,
			config: func() config.Config {
				return &mockConfig
			},
			routeConfig: func(*config.RoutesConfig) routes.RoutesConfig {
				return &mockRoutesConfig
			},
		}

		if err := c.PreRunE(newTestRoutesCmd(c)); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_build_failure", func(t *testing.T) {
		mock.Reset(ctx, &mockRoutesConfig, &mockBuilder)
		mock.Set(&mockRoutesConfig)

		mockBuilder.BuildFn = func() (routes.BuildResult, error) {
			return nil, errors.New("error")
		}

		c := &routesctx{
			ctx: ctx,
			routeBuilder: func(*routes.Config) routes.Builder {
				return &mockBuilder
			},
		}

		if err := c.RunE(); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("rune_generate_failure", func(t *testing.T) {
		mock.Reset(&mockRoutesConfig, &mockBuilder, &mockGenerator, &mockTemplate)
		mock.Set(&mockRoutesConfig)

		mockGenerator.GenerateFn = func(s string) error {
			return errors.New("error")
		}

		c := &routesctx{
			cfg: mockRoutesConfig.Config,
			routeBuilder: func(*routes.Config) routes.Builder {
				return &mockBuilder
			},
			routeGenerator: func(t routes.Template) routes.Generator {
				return &mockGenerator
			},
			routeTemplate: func(r []routes.Route, ts routes.TreeSelector, c *routes.Config) routes.Template {
				return &mockTemplate
			},
		}

		if err := c.RunE(); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("postrun_runner_failure", func(t *testing.T) {
		mock.Reset(ctx, &mockRoutesConfig, &mockConfig, &mockGenerator, &mockTemplate, &mockBuilder)
		mock.Set(&mockRoutesConfig)

		ctx.RunnerMock.RunFn = func() error {
			return errors.New("err")
		}

		c := &routesctx{
			ctx: ctx,
			cfg: mockRoutesConfig.Config,
			config: func() config.Config {
				return &mockConfig
			},
			routeConfig: func(rc *config.RoutesConfig) routes.RoutesConfig {
				return &mockRoutesConfig
			},
			routeGenerator: func(t routes.Template) routes.Generator {
				return &mockGenerator
			},
			routeTemplate: func(r []routes.Route, ts routes.TreeSelector, c *routes.Config) routes.Template {
				return &mockTemplate
			},
			routeBuilder: func(*routes.Config) routes.Builder {
				return &mockBuilder
			},
		}

		cmd := newTestRoutesCmd(c)
		if err := c.PostRunE(cmd); err == nil {
			t.Fatal("fail")
		}
	})

	t.Run("postrun_flag_error", func(t *testing.T) {
		c := &routesctx{
			ctx: ctx,
		}

		cmd := newTestRoutesCmd(c)

		cmd.Flags().Lookup("no-format").Value = &cobramock.FlagBroken{}
		if err := c.PostRunE(cmd); err == nil {
			t.Fatal("fail")
		}
	})
}

func TestGenRoutes_Fallback(t *testing.T) {
	t.Run("new_ctx_initialization", func(t *testing.T) {
		c := NewRoutesCtx(appmock.NewContext())

		rctx := c.(*routesctx)
		if rctx.config == nil || rctx.routeGenerator == nil {
			t.Fatal("fail")
		}
	})
}
