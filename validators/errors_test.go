package validators

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMultiErr_AddAndCount(t *testing.T) {
	m := &MultiErr{}

	assert.Equal(t, 0, m.Count(), "expected empty MultiErr to have count 0")

	// Adding nil should not change state
	m.Add(nil)
	assert.Equal(t, 0, m.Count(), "expected count to remain 0 after adding nil")

	err1 := fmt.Errorf("first")
	m.Add(err1)
	assert.Equal(t, 1, m.Count(), "expected count 1 after adding an error")
	assert.Equal(t, err1, m.Errors[0], "expected stored error to be err1")

	// Adding another nil should not change
	m.Add(nil)
	assert.Equal(t, 1, m.Count(), "expected count to remain 1 after adding nil")
	assert.Equal(t, "first", m.Error(), "expected single error to render its own message")
}

func TestMultiErr_AddMsg(t *testing.T) {
	m := &MultiErr{}

	m.AddMsg("")
	assert.Equal(t, 0, m.Count(), "expected empty message to be ignored")

	m.AddMsg("hello world")
	assert.Equal(t, 1, m.Count(), "expected 1 error after AddMsg")
	assert.Equal(t, "hello world", m.Errors[0].Error(), "expected error message")
}

func TestMultiErr_Addf(t *testing.T) {
	m := &MultiErr{}
	m.Addf("value %d: %s", 42, "ok")

	assert.Equal(t, 1, m.Count(), "expected 1 error after Addf")
	assert.Equal(t, "value 42: ok", m.Errors[0].Error(), "expected formatted error")
}

func TestMultiErr_AnyAndHasErrors(t *testing.T) {
	m := &MultiErr{}
	assert.False(t, m.Any(), "expected empty MultiErr to report no errors")
	assert.False(t, m.HasErrors(), "expected empty MultiErr to report no errors")

	m.Add(fmt.Errorf("e"))
	assert.True(t, m.Any(), "expected non-empty MultiErr to report errors")
	assert.True(t, m.HasErrors(), "expected non-empty MultiErr to report errors")
}

func TestMultiErr_ErrorFormatting(t *testing.T) {
	// Empty should produce empty string
	var empty MultiErr
	assert.Equal(t, "", empty.Error(), "expected empty error string for empty MultiErr")

	m := &MultiErr{}
	m.Add(fmt.Errorf("alpha"))
	m.Add(fmt.Errorf("beta"))

	s := m.Error()
	assert.Equal(t, "multiple errors occurred\n  - alpha\n  - beta\n", s)
	assert.True(t, strings.HasSuffix(s, "\n"), "expected final newline")
}

func Test_errorOrNil(t *testing.T) {
	assert.Nil(t, errorOrNil(nil), "expected nil for nil input")

	empty := &MultiErr{}
	assert.Nil(t, errorOrNil(empty), "expected nil for empty MultiErr")

	nonEmpty := &MultiErr{}
	nonEmpty.Add(fmt.Errorf("x"))
	got := errorOrNil(nonEmpty)
	assert.NotNil(t, got, "expected non-nil for MultiErr with errors")
	assert.Same(t, nonEmpty, got, "expected returned error to be the same MultiErr instance")
}
