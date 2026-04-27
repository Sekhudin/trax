package app

import (
	"runtime/debug"
	"testing"
)

func TestVersion_ExplicitVersion(t *testing.T) {
	got := Version("1.2.3")

	if got != "1.2.3" {
		t.Fatalf("expected 1.2.3, got %s", got)
	}
}

func TestVersion_BuildInfo_ValidVersion(t *testing.T) {
	old := readBuildInfo
	defer func() { readBuildInfo = old }()

	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{
				Version: "v1.0.0",
			},
		}, true
	}

	got := Version("")

	if got != "v1.0.0" {
		t.Fatalf("expected v1.0.0, got %s", got)
	}
}

func TestVersion_BuildInfo_EmptyVersion(t *testing.T) {
	old := readBuildInfo
	defer func() { readBuildInfo = old }()

	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{
				Version: "",
			},
		}, true
	}

	got := Version("")

	if got != "dev" {
		t.Fatalf("expected dev, got %s", got)
	}
}

func TestVersion_BuildInfo_DevelVersion(t *testing.T) {
	old := readBuildInfo
	defer func() { readBuildInfo = old }()

	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{
				Version: "(devel)",
			},
		}, true
	}

	got := Version("")

	if got != "dev" {
		t.Fatalf("expected dev, got %s", got)
	}
}

func TestVersion_BuildInfo_NotAvailable(t *testing.T) {
	old := readBuildInfo
	defer func() { readBuildInfo = old }()

	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return nil, false
	}

	got := Version("")

	if got != "dev" {
		t.Fatalf("expected dev, got %s", got)
	}
}
