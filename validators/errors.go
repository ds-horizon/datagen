package validators

import (
	"fmt"
	"strings"
)

type MultiErr struct {
	Errors []error
}

func (m *MultiErr) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m *MultiErr) AddMsg(msg string) {
	if msg == "" {
		return
	}
	m.Errors = append(m.Errors, fmt.Errorf("%s", msg))
}

func (m *MultiErr) Addf(format string, args ...any) {
	m.Errors = append(m.Errors, fmt.Errorf(format, args...))
}

func (m *MultiErr) Any() bool {
	return len(m.Errors) > 0
}

func (m *MultiErr) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("multiple errors occurred\n")
	for _, err := range m.Errors {
		b.WriteString("  - ")
		b.WriteString(err.Error())
		b.WriteString("\n")
	}
	return b.String()
}

func (m *MultiErr) HasErrors() bool {
	return len(m.Errors) > 0
}

func (m *MultiErr) Count() int {
	return len(m.Errors)
}

func errorOrNil(errs *MultiErr) error {
	if errs != nil && errs.Any() {
		return errs
	}
	return nil
}
