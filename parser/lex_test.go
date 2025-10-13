package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDitchSpaces(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    " foobar",
			expected: "foobar",
		},
		{
			input:    "\nfoobar",
			expected: "foobar",
		},
		{
			input:    "\t\t\n   foobar",
			expected: "foobar",
		},
		{
			input:    "foobar",
			expected: "foobar",
		},
	}

	for _, test := range tests {
		l := lex{input: test.input}
		l.ditchSpaces()
		assert.Equal(t, test.expected, l.input[l.curPos:])
	}

}

func TestDitchComments(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    " foobar ;; cmt",
			expected: " foobar ;; cmt",
		},
		{
			input:    ";;  cmt\nfoo",
			expected: "foo",
		},
		{
			input:    ";; \t\t   foobar\nfoo",
			expected: "foo",
		},
	}

	for _, test := range tests {
		l := lex{input: test.input}
		l.ditchComments()
		assert.Equal(t, test.expected, l.input[l.curPos:])
	}

}

func TestConsumeString(t *testing.T) {
	tests := []struct {
		input     string
		expected  string
		remaining string
	}{
		{
			input:     " foobar ;; cmt",
			expected:  "foobar",
			remaining: " ;; cmt",
		},
		{
			input:     ";;  cmt\nfoo",
			expected:  "foo",
			remaining: "",
		},
		{
			input:     ";; \t\t   foobar\nfoo",
			expected:  "foo",
			remaining: "",
		},
		{
			input:     ";; \t\t   foobar\nfoo() bar",
			expected:  "foo",
			remaining: "() bar",
		},
		{
			input:     "foo{} bar",
			expected:  "foo",
			remaining: "{} bar",
		},
		{
			input:     "  {} bar",
			expected:  "",
			remaining: "{} bar",
		},
	}

	for _, test := range tests {
		l := lex{input: test.input}
		val := l.consumeString()
		assert.Equal(t, test.expected, val)
		assert.Equal(t, test.remaining, l.input[l.curPos:])
	}
}
