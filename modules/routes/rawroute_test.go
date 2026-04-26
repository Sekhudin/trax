package routes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sekhudin/trax/internal/path"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	fp := filepath.Join(dir, "routes.yaml")

	if err := os.WriteFile(fp, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	return fp
}

func fileCfg(fp string) *Config {
	return &Config{
		IsFileStrategy: true,
		File:           &path.FilePath{Full: fp},
	}
}

func TestReadFile_AllBranches(t *testing.T) {
	tests := []struct {
		name    string
		content string
		path    string
		wantErr bool
	}{
		{
			name:    "config not found",
			path:    "/not/exist.yaml",
			wantErr: true,
		},
		{
			name: "unmarshal error",
			content: `
routes
wrong_field: "oops"
`,
			wantErr: true,
		},
		{
			name: "empty routes",
			content: `
routes: []
`,
			wantErr: true,
		},
		{
			name: "success",
			content: `
routes:
  - name: users
    path: /users
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := tt.path
			if tt.content != "" {
				fp = writeFile(t, tt.content)
			}

			b := rawroutebuilder{
				cfg: fileCfg(fp),
			}

			rws, err := b.readFile()
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatal(err)
				}
				if len(rws) != 1 {
					t.Fatal("expected 1 route")
				}
			}
		})
	}
}

func TestBuild_FileStrategy(t *testing.T) {
	fp := writeFile(t, `
routes:
  - name: users
    path: /users
`)
	b := newRawRouteBuilder(fileCfg(fp))

	rws, err := b.build()
	if err != nil {
		t.Fatal(err)
	}
	if len(rws) != 1 {
		t.Fatal("expected 1 route")
	}
}

func TestReadDisc_InvalidStrategy(t *testing.T) {
	b := rawroutebuilder{
		cfg: &Config{
			IsFileStrategy: false,
			Strategy:       "invalid",
		},
	}

	_, err := b.readDisc()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBuild_ReadDisc_WalkerError(t *testing.T) {
	tests := []string{"next-app", "next-page"}

	for _, s := range tests {
		t.Run(s, func(t *testing.T) {
			cfg := &Config{
				IsFileStrategy: false,
				Strategy:       s,
				Root:           "/definitely/not/exist",
			}

			b := newRawRouteBuilder(cfg)

			_, err := b.build()
			if err == nil {
				t.Fatal("expected walker error")
			}
		})
	}
}
