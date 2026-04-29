package fsmock

type Writer struct {
	WriteCalled bool
	WriteFn     func(string, []byte) error
}

func (w *Writer) Reset() {
	w.WriteCalled = false

	w.WriteFn = func(s string, b []byte) error {
		return nil
	}
}

func (w *Writer) Write(p string, d []byte) error {
	w.WriteCalled = true
	if w.WriteFn != nil {
		return w.WriteFn(p, d)
	}
	return nil
}
