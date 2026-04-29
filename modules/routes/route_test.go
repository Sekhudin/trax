package routes

import (
	"testing"
)

func TestRouteBuilder_Success(t *testing.T) {
	b := NewRouteBuilder(&Config{Prefix: "/api"})

	t.Run("split_with_prefix", func(t *testing.T) {
		rw := RawRoute{Path: "/users"}
		parts := b.(*route).splitPath(rw)
		if len(parts) != 2 || parts[0] != "api" {
			t.Fatal("should_include_prefix_parts")
		}
	})

	t.Run("build_complete_route", func(t *testing.T) {
		rws := []RawRoute{{Name: "home", Path: "/dashboard?ref=test"}}
		rs, err := b.Build(rws)
		if err != nil || rs[0].Path != "/dashboard" {
			t.Fatal("should_build_without_query")
		}
	})
}

func TestRouteBuilder_Error(t *testing.T) {
	b := NewRouteBuilder(&Config{})
	r := b.(*route)

	t.Run("missing_leading_slash", func(t *testing.T) {
		_, err := r.cleanPath(RawRoute{Path: "no-slash"})
		if err == nil {
			t.Fatal("should_reject_no_slash")
		}
	})

	t.Run("contains doule slash", func(t *testing.T) {
		_, err := r.cleanPath(RawRoute{Path: "/foo//bar"})
		if err == nil {
			t.Fatal("should_reject_double_slash")
		}
	})

	t.Run("wildcard_mid_segment", func(t *testing.T) {
		err := r.validateParts(RawRoute{}, []string{"a", "*", "b"})
		if err == nil {
			t.Fatal("should_reject_misplaced_wildcard")
		}
	})

	t.Run("invalid_wildcard_name", func(t *testing.T) {
		err := r.validateParts(RawRoute{}, []string{"user*name"})
		if err == nil {
			t.Fatal("should_reject_mixed_wildcard")
		}
	})

	t.Run("build_invalid_path", func(t *testing.T) {
		rws := []RawRoute{{Name: "bad", Path: "invalid"}}
		_, err := b.Build(rws)
		if err == nil {
			t.Fatal("should_stop_on_error")
		}
	})

	t.Run("build_invalid_segment", func(t *testing.T) {
		rws := []RawRoute{{Name: "bad", Path: "/foo/*/bar"}}
		_, err := b.Build(rws)
		if err == nil {
			t.Fatal("should_stop_on_error")
		}
	})
}

func TestRouteBuilder_Fallback(t *testing.T) {
	t.Run("root_path_exception", func(t *testing.T) {
		r := &route{}
		out, _ := r.cleanPath(RawRoute{Path: "/"})
		if out != "/" {
			t.Fatal("should_keep_root_slash")
		}
	})

	t.Run("interface_compliance_check", func(t *testing.T) {
		var _ RouteBuilder = (*route)(nil)
	})
}
