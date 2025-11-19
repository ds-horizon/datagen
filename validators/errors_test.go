package validators

import (
	"fmt"
	"strings"
	"testing"
)

func TestMultiErr_AddAndCount(t *testing.T) {
	m := &MultiErr{}

	if m.Count() != 0 {
		t.Fatalf("expected empty MultiErr to have count 0, got %d", m.Count())
	}

	// Adding nil should not change state
	m.Add(nil)
	if m.Count() != 0 {
		t.Fatalf("expected count to remain 0 after adding nil, got %d", m.Count())
	}

	err1 := fmt.Errorf("first")
	m.Add(err1)
	if m.Count() != 1 {
		t.Fatalf("expected count 1 after adding an error, got %d", m.Count())
	}
	if m.Errors[0] != err1 {
		t.Fatalf("expected stored error to be err1")
	}

	// Adding another nil should not change
	m.Add(nil)
	if m.Count() != 1 {
		t.Fatalf("expected count to remain 1 after adding nil, got %d", m.Count())
	}
}

func TestMultiErr_AddMsg(t *testing.T) {
	m := &MultiErr{}

	m.AddMsg("")
	if m.Count() != 0 {
		t.Fatalf("expected empty message to be ignored, got count %d", m.Count())
	}

	m.AddMsg("hello world")
	if m.Count() != 1 {
		t.Fatalf("expected 1 error after AddMsg, got %d", m.Count())
	}
	if got := m.Errors[0].Error(); got != "hello world" {
		t.Fatalf("expected error message 'hello world', got %q", got)
	}
}

func TestMultiErr_Addf(t *testing.T) {
	m := &MultiErr{}
	m.Addf("value %d: %s", 42, "ok")

	if m.Count() != 1 {
		t.Fatalf("expected 1 error after Addf, got %d", m.Count())
	}
	if got := m.Errors[0].Error(); got != "value 42: ok" {
		t.Fatalf("expected formatted error 'value 42: ok', got %q", got)
	}
}

func TestMultiErr_AnyAndHasErrors(t *testing.T) {
	m := &MultiErr{}
	if m.Any() || m.HasErrors() {
		t.Fatalf("expected empty MultiErr to report no errors")
	}

	m.Add(fmt.Errorf("e"))
	if !m.Any() || !m.HasErrors() {
		t.Fatalf("expected non-empty MultiErr to report errors")
	}
}

func TestMultiErr_ErrorFormatting(t *testing.T) {
	// Empty should produce empty string
	var empty MultiErr
	if s := empty.Error(); s != "" {
		t.Fatalf("expected empty error string for empty MultiErr, got %q", s)
	}

	m := &MultiErr{}
	m.Add(fmt.Errorf("alpha"))
	m.Add(fmt.Errorf("beta"))

	s := m.Error()
	if !strings.HasPrefix(s, "multiple errors occurred\n") {
		t.Fatalf("expected header 'multiple errors occurred', got %q", s)
	}
	if !strings.Contains(s, "  - alpha\n") || !strings.Contains(s, "  - beta\n") {
		t.Fatalf("expected each error on its own line with bullet, got %q", s)
	}
	if !strings.HasSuffix(s, "\n") {
		t.Fatalf("expected final newline, got %q", s)
	}
}

func Test_errorOrNil(t *testing.T) {
	if got := errorOrNil(nil); got != nil {
		t.Fatalf("expected nil for nil input, got %v", got)
	}

	empty := &MultiErr{}
	if got := errorOrNil(empty); got != nil {
		t.Fatalf("expected nil for empty MultiErr, got %v", got)
	}

	nonEmpty := &MultiErr{}
	nonEmpty.Add(fmt.Errorf("x"))
	if got := errorOrNil(nonEmpty); got == nil {
		t.Fatalf("expected non-nil for MultiErr with errors")
	} else if got != nonEmpty {
		t.Fatalf("expected returned error to be the same MultiErr instance")
	}
}


