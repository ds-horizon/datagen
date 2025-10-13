package validators

import (
	"fmt"
	"strings"
)

type MultiErr []string

func (m *MultiErr) Add(msg string) {
	if msg == "" {
		return
	}
	*m = append(*m, msg)
}

func (m *MultiErr) Addf(format string, args ...any) {
	*m = append(*m, fmt.Sprintf(format, args...))
}

func (m *MultiErr) Any() bool {
	if m == nil {
		return false
	}
	return len(*m) > 0
}

func (m *MultiErr) Error() string {
	if m == nil {
		return ""
	}
	return strings.Join(*m, ";\n")
}

func errorOrNil(errs *MultiErr) error {
	if errs != nil && errs.Any() {
		return errs
	}
	return nil
}
