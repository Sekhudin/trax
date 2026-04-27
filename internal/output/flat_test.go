package output

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type flatErrWriter struct{}

func (e flatErrWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write error")
}

func newFlatTestContext(w io.Writer) Context {
	return New(w, Options{})
}

func newFlatTestContextStruct(w io.Writer) *context {
	opt := Options{}

	return &context{
		w:     w,
		opt:   opt,
		color: NewColorizer(opt.NoColor),
	}
}

func TestAsFlat_SimpleMap(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newFlatTestContext(buf)

	data := map[string]any{
		"b": 2,
		"a": 1,
	}

	err := ctx.AsFlat("", data)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	if !strings.Contains(out, "a = 1") || !strings.Contains(out, "b = 2") {
		t.Fatal(out)
	}
}

func TestAsFlat_WithPrefix(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newFlatTestContext(buf)

	data := map[string]any{
		"x": 10,
	}

	err := ctx.AsFlat("root", data)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "root.x = 10") {
		t.Fatal(buf.String())
	}
}

func TestAsFlat_NestedMap(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newFlatTestContext(buf)

	data := map[string]any{
		"parent": map[string]any{
			"child": 5,
		},
	}

	err := ctx.AsFlat("", data)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "parent.child = 5") {
		t.Fatal(buf.String())
	}
}

func TestAsFlat_ArrayHandling(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newFlatTestContext(buf)

	data := map[string]any{
		"arr": []any{"a", "b"},
	}

	err := ctx.AsFlat("", data)
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "arr.0 = a") ||
		!strings.Contains(out, "arr.1 = b") {
		t.Fatal(out)
	}
}

func TestAsFlat_DeepNestedMix(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newFlatTestContext(buf)

	data := map[string]any{
		"a": map[string]any{
			"b": []any{
				map[string]any{"c": 7},
			},
		},
	}

	err := ctx.AsFlat("", data)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "a.b.0 = map[c:7]") {
		t.Fatal(buf.String())
	}
}

func TestAsFlat_WriteErrorPropagation(t *testing.T) {
	ctx := newFlatTestContext(flatErrWriter{})

	err := ctx.AsFlat("", map[string]any{"a": 1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPrintFlatValue_Direct(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newFlatTestContextStruct(buf)

	err := ctx.printFlatValue("key", 123)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != "key = 123\n" {
		t.Fatal(buf.String())
	}
}
