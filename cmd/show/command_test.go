package show

import (
	"testing"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/testutil/appmock"
	"github.com/sekhudin/trax/internal/testutil/mock"
)

func TestShow_Success(t *testing.T) {
	ctx := appmock.NewContext()

	t.Run("default_deps_structure", func(t *testing.T) {
		mock.Reset(ctx)

		d := DefaultDependencies()
		if d.Docs == nil || d.NewConfigCtx == nil || d.NewRoutesCtx == nil {
			t.Fatal("fail")
		}

		if d.Docs.Root.Use != "show" || d.Docs.Config.Use != "config" || d.Docs.Routes.Use != "routes" {
			t.Fatal("fail")
		}
	})

	t.Run("build_command_tree", func(t *testing.T) {
		mock.Reset(ctx)

		NewConfigCalled, NewRoutesCalled := false, false

		deps := &Dependencies{
			Docs: &Docs{
				Root:   doc.Docs{Use: "root"},
				Config: doc.Docs{Use: "cfg"},
				Routes: doc.Docs{Use: "rts"},
			},
			NewConfigCtx: func(app.Context) ConfigCtx {
				NewConfigCalled = true
				return &mockConfigCtx{}
			},
			NewRoutesCtx: func(app.Context) RoutesCtx {
				NewRoutesCalled = true
				return &mockRoutesCtx{}
			},
		}

		cmd := NewWithDependencies(ctx, deps)
		if cmd.Use != "root" || len(cmd.Commands()) != 2 {
			t.Fatal("fail")
		}

		if !NewConfigCalled || !NewRoutesCalled {
			t.Fatal("fail")
		}
	})

	t.Run("new_default_instance", func(t *testing.T) {
		mock.Reset(ctx)

		cmd := New(ctx)
		if cmd.Use != "show" || len(cmd.Commands()) != 2 {
			t.Fatal("fail")
		}
	})
}

func TestShow_Error(t *testing.T) {
	t.Run("command_attachment_check", func(t *testing.T) {
		deps := &Dependencies{
			Docs: &Docs{
				Root:   doc.Docs{Use: "root"},
				Config: doc.Docs{Use: "cfg"},
				Routes: doc.Docs{Use: "rts"},
			},

			NewConfigCtx: func(app.Context) ConfigCtx {
				return &mockConfigCtx{}
			},

			NewRoutesCtx: func(app.Context) RoutesCtx {
				return &mockRoutesCtx{}
			},
		}

		cmd := NewWithDependencies(appmock.NewContext(), deps)

		foundCfg, foundRts := false, false
		for _, c := range cmd.Commands() {
			if c.Use == "cfg" {
				foundCfg = true
			}
			if c.Use == "rts" {
				foundRts = true
			}
		}
		if !foundCfg || !foundRts {
			t.Fatal("fail")
		}
	})
}

func TestShow_Fallback(t *testing.T) {
	t.Run("verify_subcommand_uses", func(t *testing.T) {
		d := DefaultDependencies()

		if d.Docs.Config.Use != "config" || d.Docs.Routes.Use != "routes" {
			t.Fatal("fail")
		}
	})
}
