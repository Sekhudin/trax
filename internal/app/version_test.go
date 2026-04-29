package app

import (
	"runtime/debug"
	"testing"
)

func TestVersion_Success(t *testing.T) {
	t.Run("return_explicit_version", func(t *testing.T) {
		got := Version("1.2.3")
		if got != "1.2.3" {
			t.Fatalf("expected 1.2.3, got %s", got)
		}
	})

	t.Run("use_build_info", func(t *testing.T) {
		old := readBuildInfo
		defer func() { readBuildInfo = old }()

		readBuildInfo = func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{Main: debug.Module{Version: "v1.0.0"}}, true
		}

		got := Version("")
		if got != "v1.0.0" {
			t.Fatalf("expected v1.0.0, got %s", got)
		}
	})
}

func TestVersion_Fallback(t *testing.T) {
	old := readBuildInfo
	defer func() { readBuildInfo = old }()

	t.Run("handle_empty_info", func(t *testing.T) {
		readBuildInfo = func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{Main: debug.Module{Version: ""}}, true
		}
		if got := Version(""); got != "dev" {
			t.Fatalf("expected dev, got %s", got)
		}
	})

	t.Run("handle_devel_tag", func(t *testing.T) {
		readBuildInfo = func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, true
		}
		if got := Version(""); got != "dev" {
			t.Fatalf("expected dev, got %s", got)
		}
	})

	t.Run("handle_missing_info", func(t *testing.T) {
		readBuildInfo = func() (*debug.BuildInfo, bool) { return nil, false }
		if got := Version(""); got != "dev" {
			t.Fatalf("expected dev, got %s", got)
		}
	})
}
