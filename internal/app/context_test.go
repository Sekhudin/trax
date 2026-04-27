package app

import (
	"bytes"
	"testing"

	"github.com/sekhudin/trax/internal/output"
	"github.com/spf13/cobra"
)

func TestNew_ContextInitialization(t *testing.T) {
	opt := output.Options{
		NoColor: true,
	}

	ctx := New(opt)

	if ctx == nil {
		t.Fatal("expected context, got nil")
	}

	if ctx.Color() == nil {
		t.Fatal("expected Colorizer to be initialized")
	}

	if ctx.Out() == nil {
		t.Fatal("expected Output Context to be initialized")
	}

	if ctx.Runner() == nil {
		t.Fatal("expected Runner to be initialized")
	}
}

func TestApplyOptions_BasicUpdate(t *testing.T) {
	opt := output.Options{
		NoColor: true,
	}

	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	ctx := New(output.Options{})

	ctx.ApplyOptions(cmd, opt)

	if ctx.Color() == nil {
		t.Fatal("color should not be nil")
	}

	if ctx.Out() == nil {
		t.Fatal("out should not be nil")
	}

	if ctx.Runner() == nil {
		t.Fatal("runner should not be nil")
	}
}

func TestApplyOptions_WritersPropagation(t *testing.T) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	cmd := &cobra.Command{}
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)

	ctx := New(output.Options{})

	opt := output.Options{
		NoColor: true,
	}

	ctx.ApplyOptions(cmd, opt)

	if ctx.Out() == nil {
		t.Fatal("expected Out to be set")
	}
}

func TestApplyOptions_NilSafety(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	ctx := New(output.Options{})

	cmd := &cobra.Command{}
	opt := output.Options{}

	ctx.ApplyOptions(cmd, opt)
}

func TestContext_FullFlow(t *testing.T) {
	opt := output.Options{
		Debug:   true,
		NoColor: true,
	}

	ctx := New(opt)

	if ctx.Color() == nil || ctx.Out() == nil || ctx.Runner() == nil {
		t.Fatal("invalid initial context")
	}

	cmd := &cobra.Command{}
	ctx.ApplyOptions(cmd, opt)

	if ctx.Out() == nil {
		t.Fatal("Out should still be valid after ApplyOptions")
	}
}
