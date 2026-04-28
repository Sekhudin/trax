package routes

import (
	"testing"
)

func TestCleanPath_AllBranches(t *testing.T) {
	b := route{
		cfg: &Config{},
	}

	tests := []struct {
		name    string
		rw      RawRoute
		wantErr bool
		want    string
	}{
		{
			name:    "no leading slash",
			rw:      RawRoute{Name: "a", Path: "users"},
			wantErr: true,
		},
		{
			name:    "double slash",
			rw:      RawRoute{Name: "a", Path: "/users//list"},
			wantErr: true,
		},
		{
			name: "with query trimmed",
			rw:   RawRoute{Name: "a", Path: "/users?id=1"},
			want: "/users",
		},
		{
			name: "trailing slash trimmed",
			rw:   RawRoute{Name: "a", Path: "/users/"},
			want: "/users",
		},
		{
			name: "root path stays",
			rw:   RawRoute{Name: "a", Path: "/"},
			want: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := b.cleanPath(tt.rw)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatal(err)
				}
				if out != tt.want {
					t.Fatalf("want %s got %s", tt.want, out)
				}
			}
		})
	}
}

func TestSplitPath_WithAndWithoutPrefix(t *testing.T) {
	b := route{
		cfg: &Config{
			Prefix: "/api",
		},
	}

	rw := RawRoute{Path: "/users/list"}
	parts := b.splitPath(rw)

	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %v", parts)
	}

	b.cfg.Prefix = ""
	parts = b.splitPath(rw)
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %v", parts)
	}
}

func TestValidateParts_AllBranches(t *testing.T) {
	b := route{
		cfg: &Config{},
	}

	tests := []struct {
		name    string
		parts   []string
		wantErr bool
	}{
		{
			name:    "wildcard not last",
			parts:   []string{"users", "*", "list"},
			wantErr: true,
		},
		{
			name:    "wildcard mixed chars",
			parts:   []string{"users", "ab*cd"},
			wantErr: true,
		},
		{
			name:    "valid wildcard last",
			parts:   []string{"users", "*"},
			wantErr: false,
		},
		{
			name:    "no wildcard",
			parts:   []string{"users", "list"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := b.validateParts(RawRoute{Name: "a"}, tt.parts)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBuild_AllBranches(t *testing.T) {
	b := NewRouteBuilder(&Config{})

	rws := []RawRoute{
		{Name: "ok", Path: "/users"},
		{Name: "bad", Path: "no-slash"},
	}

	_, err := b.Build(rws)
	if err == nil {
		t.Fatal("expected error from cleanPath")
	}

	rws = []RawRoute{
		{Name: "ok", Path: "/users/*/list"},
	}

	_, err = b.Build(rws)
	if err == nil {
		t.Fatal("expected error from validateParts")
	}

	rws = []RawRoute{
		{Name: "ok", Path: "/users"},
	}

	rs, err := b.Build(rws)
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != 1 {
		t.Fatal("expected 1 route")
	}
}
