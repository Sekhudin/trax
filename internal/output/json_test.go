package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	appErr "github.com/sekhudin/trax/internal/errors"
)

type errWriter struct{}

func (e errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func newTestContext(w io.Writer) *Context {
	return New(w, Options{})
}

func TestAsJSON_Success(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newTestContext(buf)

	data := map[string]any{
		"b": 2,
		"a": 1,
		"nested": map[string]any{
			"d": 4,
			"c": 3,
		},
		"arr": []any{
			map[string]any{"z": 9, "y": 8},
			"plain",
		},
	}

	err := ctx.AsJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatal("invalid json output:", err)
	}

	if out["a"].(float64) != 1 {
		t.Fatal(out)
	}
	if out["b"].(float64) != 2 {
		t.Fatal(out)
	}

	nested := out["nested"].(map[string]any)
	if nested["c"].(float64) != 3 || nested["d"].(float64) != 4 {
		t.Fatal(nested)
	}

	arr := out["arr"].([]any)
	first := arr[0].(map[string]any)
	if first["y"].(float64) != 8 || first["z"].(float64) != 9 {
		t.Fatal(arr)
	}
}

func TestAsJSON_WriteError(t *testing.T) {
	ctx := newTestContext(errWriter{})

	err := ctx.AsJSON(map[string]any{"a": 1})
	if err == nil {
		t.Fatal("expected error")
	}

	coreErr, ok := err.(*appErr.CoreError)
	if !ok {
		t.Fatalf("unexpected error type: %T", err)
	}

	if coreErr.Code != appErr.ErrIO {
		t.Fatalf("expected IO error, got %s", coreErr.Code)
	}
}

func TestNormalizeMap_SortingAndRecursion(t *testing.T) {
	ctx := newTestContext(bytes.NewBuffer(nil))

	in := map[string]any{
		"b": 1,
		"a": map[string]any{
			"d": 4,
			"c": 3,
		},
	}

	out := ctx.normalizeMap(in)

	a := out["a"].(map[string]any)
	if _, ok := a["c"]; !ok {
		t.Fatal(a)
	}
	if _, ok := a["d"]; !ok {
		t.Fatal(a)
	}
}

func TestNormalizeValue_ArrayRecursion(t *testing.T) {
	ctx := newTestContext(bytes.NewBuffer(nil))

	in := []any{
		map[string]any{"b": 2, "a": 1},
		"x",
	}

	out := ctx.normalizeValue(in).([]any)

	m := out[0].(map[string]any)
	if m["a"].(float64) != 1 || m["b"].(float64) != 2 {
		t.Fatal(out)
	}
}

func TestAsJSON_OutputIsPrettyPrinted(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := newTestContext(buf)

	_ = ctx.AsJSON(map[string]any{"a": 1})

	if !strings.Contains(buf.String(), "\n") {
		t.Fatal("expected pretty printed json")
	}
}
