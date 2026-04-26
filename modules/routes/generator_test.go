package routes

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/fs"
)

type fakeWriter struct {
	path   string
	data   []byte
	called bool
	err    error
}

func (f *fakeWriter) Write(p string, d []byte) error {
	f.called = true
	f.path = p
	f.data = d
	return f.err
}

type fakeTemplate struct {
	out string
	err error
}

func (f fakeTemplate) Build() (string, error) {
	return f.out, f.err
}

func makeGen(w fs.FileWriter, tpl Template) Generator {
	return NewGenerator(w, tpl)
}

func TestGenerator_Generate_Success(t *testing.T) {
	w := &fakeWriter{}
	tpl := fakeTemplate{out: "generated content"}

	g := makeGen(w, tpl)

	err := g.Generate("/tmp/out.ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !w.called {
		t.Fatal("writer not called")
	}

	if w.path != "/tmp/out.ts" {
		t.Fatalf("wrong path: %s", w.path)
	}

	if string(w.data) != "generated content" {
		t.Fatalf("wrong content: %s", string(w.data))
	}
}

func TestGenerator_Generate_TemplateError(t *testing.T) {
	w := &fakeWriter{}
	tpl := fakeTemplate{err: errors.New("template failed")}

	g := makeGen(w, tpl)

	err := g.Generate("/tmp/out.ts")
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if w.called {
		t.Fatal("writer should not be called when template fails")
	}
}

func TestGenerator_Generate_WriteError(t *testing.T) {
	w := &fakeWriter{err: errors.New("write failed")}
	tpl := fakeTemplate{out: "data"}

	g := makeGen(w, tpl)

	err := g.Generate("/tmp/out.ts")
	if err == nil {
		t.Fatal("expected write error")
	}
}
