package app

import (
	"bytes"
	"testing"

	"github.com/sekhudin/trax/internal/output"
	"github.com/spf13/cobra"
)

func TestNew_Success(t *testing.T) {
	t.Run("initialize_all_components", func(t *testing.T) {
		opt := output.Options{NoColor: true}
		ctx := New(opt)

		if ctx == nil {
			t.Fatal("ctx_is_nil")
		}
		if ctx.Color() == nil {
			t.Fatal("color_is_nil")
		}
		if ctx.Out() == nil {
			t.Fatal("out_is_nil")
		}
		if ctx.Runner() == nil {
			t.Fatal("runner_is_nil")
		}
	})
}

func TestApplyOptions_Success(t *testing.T) {
	ctx := New(output.Options{})
	opt := output.Options{NoColor: true}

	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	t.Run("update_internal_fields", func(t *testing.T) {
		ctx.ApplyOptions(cmd, opt)
		if ctx.Color() == nil || ctx.Out() == nil || ctx.Runner() == nil {
			t.Fatal("fields_not_updated")
		}
	})

	t.Run("propagate_command_writers", func(t *testing.T) {
		var outBuf bytes.Buffer
		cmd.SetOut(&outBuf)

		ctx.ApplyOptions(cmd, opt)
		if ctx.Out() == nil {
			t.Fatal("out_is_nil")
		}
	})
}

func TestApplyOptions_Error(t *testing.T) {
	t.Run("handle_nil_safety", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic_detected: %v", r)
			}
		}()

		ctx := New(output.Options{})
		ctx.ApplyOptions(&cobra.Command{}, output.Options{})
	})
}
