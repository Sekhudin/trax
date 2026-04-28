package show

import (
	"testing"

	"github.com/sekhudin/trax/internal/app"
	"github.com/sekhudin/trax/internal/doc"
	"github.com/sekhudin/trax/internal/testutil/appmock"
)

func Test_DefaultDependencies_structure(t *testing.T) {
	d := DefaultDependencies()

	if d.Docs == nil {
		t.Fatal("Docs nil")
	}

	if d.Docs.Root.Use != "show" {
		t.Fatal("root docs wrong")
	}

	if d.Docs.Config.Use != "config" {
		t.Fatal("config docs wrong")
	}

	if d.Docs.Routes.Use != "routes" {
		t.Fatal("routes docs wrong")
	}

	if d.NewConfigCtx == nil {
		t.Fatal("NewConfigCtx nil")
	}

	if d.NewRoutesCtx == nil {
		t.Fatal("NewRoutesCtx nil")
	}
}

func Test_NewWithDependencies_builds_command_tree(t *testing.T) {
	calledConfig := false
	calledRoutes := false

	deps := &Dependencies{
		Docs: &Docs{
			Root:   doc.Docs{Use: "root"},
			Config: doc.Docs{Use: "cfg"},
			Routes: doc.Docs{Use: "rts"},
		},
		NewConfigCtx: func(app.Context) ConfigCtx {
			calledConfig = true
			return &mockConfigCtx{}
		},
		NewRoutesCtx: func(app.Context) RoutesCtx {
			calledRoutes = true
			return &mockRoutesCtx{}
		},
	}

	cmd := NewWithDependencies(appmock.NewContext(), deps)

	if cmd.Use != "root" {
		t.Fatalf("unexpected root use: %s", cmd.Use)
	}

	children := cmd.Commands()
	if len(children) != 2 {
		t.Fatalf("expected 2 subcommands, got %d", len(children))
	}

	if children[0].Use != "cfg" && children[1].Use != "cfg" {
		t.Fatal("config command not attached")
	}

	if children[0].Use != "rts" && children[1].Use != "rts" {
		t.Fatal("routes command not attached")
	}

	if !calledConfig {
		t.Fatal("NewConfigCtx not called")
	}

	if !calledRoutes {
		t.Fatal("NewRoutesCtx not called")
	}
}

func Test_New_uses_DefaultDependencies(t *testing.T) {
	cmd := New(appmock.NewContext())

	if cmd.Use != "show" {
		t.Fatalf("expected root 'show', got %s", cmd.Use)
	}

	children := cmd.Commands()
	if len(children) != 2 {
		t.Fatalf("expected 2 subcommands, got %d", len(children))
	}
}
