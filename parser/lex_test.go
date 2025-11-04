package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ds-horizon/datagen/codegen"
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
			input:    " foobar // cmt",
			expected: " foobar // cmt",
		},
		{
			input:    "//  cmt\nfoo",
			expected: "foo",
		},
		{
			input:    "// \t\t   foobar\nfoo",
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
			input:     " foobar // cmt",
			expected:  "foobar",
			remaining: " // cmt",
		},
		{
			input:     "//  cmt\nfoo",
			expected:  "foo",
			remaining: "",
		},
		{
			input:     "// \t\t   foobar\nfoo",
			expected:  "foo",
			remaining: "",
		},
		{
			input:     "// \t\t   foobar\nfoo() bar",
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

func TestParse(t *testing.T) {
	testExpr := []struct {
		name     string
		input    string
		expected *codegen.DatagenParsed
		fail     bool
		errStr   string
	}{
		{
			name:  "empty model",
			input: "model empty {}",
			expected: &codegen.DatagenParsed{
				ModelName: "empty",
				Filepath:  "test.dg",
			},
			fail: false,
		},
		{
			name: "model with fields only",
			input: `model user {
  fields {
    name() string
    age() int
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "user",
				Filepath:  "test.dg",
			},
			fail: false,
		},
		{
			name: "model with metadata count",
			input: `model test {
  metadata {
    count: 100
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "test",
				Filepath:  "test.dg",
				Metadata: &codegen.Metadata{
					Count: 100,
					Tags:  map[string]string{},
				},
			},
			fail: false,
		},
		{
			name: "model with metadata tags",
			input: `model test {
  metadata {
    tags: {
      "service": "test",
      "team": "backend"
    }
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "test",
				Filepath:  "test.dg",
				Metadata: &codegen.Metadata{
					Tags: map[string]string{
						"service": "test",
						"team":    "backend",
					},
				},
			},
			fail: false,
		},
		{
			name: "model with misc section",
			input: `model test {
  misc {
    const Count = 100
    type TestStruct struct {
      Field1 string
    }
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "test",
				Filepath:  "test.dg",
			},
			fail: false,
		},
		{
			name: "model with gens section",
			input: `model test {
  gens {
    func name() {
      return "test"
    }
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "test",
				Filepath:  "test.dg",
			},
			fail: false,
		},
		{
			name: "model with calls section",
			input: `model test {
  fields {
    name() string
  }
  calls {
    name()
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "test",
				Filepath:  "test.dg",
			},
			fail: false,
		},
		{
			name: "model with comments",
			input: `// This is a comment
model test {
  // Another comment
  fields {
    name() string // inline comment
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "test",
				Filepath:  "test.dg",
			},
			fail: false,
		},
		{
			name: "model with metadata count and tags",
			input: `model test {
  metadata {
    count: 50
    tags: {
      "service": "test"
    }
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "test",
				Filepath:  "test.dg",
				Metadata: &codegen.Metadata{
					Count: 50,
					Tags: map[string]string{
						"service": "test",
					},
				},
			},
			fail: false,
		},
		{
			name: "model with all sections",
			input: `model complete {
  metadata {
    count: 42
    tags: {
      "key": "value"
    }
  }
  misc {
    const Test = 1
  }
  fields {
    id() int
    name() string
  }
  gens {
    func id() {
      return 1
    }
    func name() {
      return "test"
    }
  }
  calls {
    id()
    name()
  }
}`,
			expected: &codegen.DatagenParsed{
				ModelName: "complete",
				Filepath:  "test.dg",
			},
			fail: false,
		},
		{
			name:   "missing model keyword",
			input:  "user {}",
			fail:   true,
			errStr: "expected 'model'",
		},
		{
			name:   "missing model name",
			input:  "model {}",
			fail:   true,
			errStr: "expected valid model name",
		},
		{
			name:   "missing opening brace",
			input:  "model test",
			fail:   true,
			errStr: "expected '{'",
		},
		{
			name:   "missing closing brace",
			input:  "model test {",
			fail:   true,
			errStr: "",
		},
		{
			name:   "invalid section name",
			input:  "model test { invalid {} }",
			fail:   true,
			errStr: "expected section header",
		},
		{
			name:   "incomplete metadata count",
			input:  "model test { metadata { count } }",
			fail:   true,
			errStr: "",
		},
		{
			name:   "incomplete metadata tags",
			input:  "model test { metadata { tags } }",
			fail:   true,
			errStr: "",
		},
		{
			name:   "invalid metadata field",
			input:  "model test { metadata { invalid: 10 } }",
			fail:   true,
			errStr: "invalid metadata field",
		},
		{
			name:   "incomplete gens section",
			input:  "model test { gens { func } }",
			fail:   true,
			errStr: "",
		},
		{
			name:   "incomplete gen function",
			input:  "model test { gens { func name } }",
			fail:   true,
			errStr: "",
		},
		{
			name:   "empty input",
			input:  "",
			fail:   true,
			errStr: "expected 'model'",
		},
		{
			name:   "only whitespace",
			input:  "   \n\t  ",
			fail:   true,
			errStr: "expected 'model'",
		},
		{
			name:   "unclosed body",
			input:  "model test { fields {",
			fail:   true,
			errStr: "",
		},
		{
			name:   "mismatched braces - extra opening",
			input:  "model test { fields { {",
			fail:   true,
			errStr: "",
		},
	}

	for _, tt := range testExpr {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.input), "test.dg")

			if tt.fail {
				require.Error(t, err, "expected error for input: %q", tt.input)
				if tt.errStr != "" {
					assert.Contains(t, err.Error(), tt.errStr,
						"error message should contain %q, got: %v", tt.errStr, err)
				}
				return
			}

			require.NoError(t, err, "unexpected error for input: %q", tt.input)
			require.NotNil(t, got, "expected non-nil result")

			assert.Equal(t, tt.expected.ModelName, got.ModelName,
				"ModelName mismatch")
			assert.Equal(t, tt.expected.Filepath, got.Filepath,
				"Filepath mismatch")

			if tt.expected.Metadata != nil {
				require.NotNil(t, got.Metadata, "expected non-nil Metadata")
				assert.Equal(t, tt.expected.Metadata.Count, got.Metadata.Count,
					"Metadata.Count mismatch")
				if len(tt.expected.Metadata.Tags) > 0 {
					assert.Equal(t, tt.expected.Metadata.Tags, got.Metadata.Tags,
						"Metadata.Tags mismatch")
				}
			}

			if tt.expected.Fields != nil {
				assert.NotNil(t, got.Fields, "expected non-nil Fields")
			}
			if tt.expected.Misc != "" {
				assert.NotEmpty(t, got.Misc, "expected non-empty Misc")
			}
			if len(tt.expected.GenFuns) > 0 {
				assert.NotEmpty(t, got.GenFuns, "expected non-empty GenFuns")
			}
			if len(tt.expected.Calls) > 0 {
				assert.NotEmpty(t, got.Calls, "expected non-empty Calls")
			}
		})
	}
}
