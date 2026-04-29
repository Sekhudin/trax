package routes

import (
	"errors"
	"testing"

	"github.com/sekhudin/trax/internal/testutil/fsmock"
	"github.com/sekhudin/trax/internal/testutil/mock"
)

func TestGenerator_Success(t *testing.T) {
	mockWriter := fsmock.Writer{}
	mockTemplate := mockTemplateBuilder{}

	t.Run("generate_and_write_file", func(t *testing.T) {
		g := NewGenerator(&mockWriter, &mockTemplate)

		err := g.Generate("/tmp/trax.ts")
		if err != nil {
			t.Fatal("should_generate_successfully")
		}
		if !mockWriter.WriteCalled {
			t.Fatal("writer_payload_mismatch")
		}
	})
}

func TestGenerator_Error(t *testing.T) {
	mockWriter := fsmock.Writer{}
	mockTemplate := mockTemplateBuilder{}

	t.Run("template_build_failed", func(t *testing.T) {
		mock.Reset(&mockWriter, &mockTemplate)
		mockTemplate.BuildFn = func() (string, error) {
			return "", errors.New("error")
		}

		g := NewGenerator(&mockWriter, &mockTemplate)

		if err := g.Generate("any"); err == nil {
			t.Fatal("should_catch_template_err")
		}

		if mockWriter.WriteCalled {
			t.Fatal("should_not_called")
		}
	})

	t.Run("file_write_failed", func(t *testing.T) {
		mock.Reset(&mockWriter, &mockTemplate)
		mockWriter.WriteFn = func(s string, b []byte) error {
			return errors.New("error")
		}

		g := NewGenerator(&mockWriter, &mockTemplate)

		if err := g.Generate("any"); err == nil {
			t.Fatal("should_catch_write_err")
		}
	})
}

func TestGenerator_Fallback(t *testing.T) {
	t.Run("interface_compliance_check", func(t *testing.T) {
		var _ Generator = (*generator)(nil)
	})
}
