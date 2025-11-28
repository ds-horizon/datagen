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
		name              string
		input             string
		expectedMetadata  *codegen.Metadata
		expectedModelName string
		expectedFilepath  string
		expectedFields    bool
		expectedMisc      bool
		expectedGenFuncs  bool
		expectedCalls     bool
		fail              bool
		errStr            string
	}{
		{
			name:              "empty model",
			input:             "model empty {}",
			expectedMetadata:  nil,
			expectedModelName: "empty",
			expectedFilepath:  "test.dg",
			expectedFields:    false,
			expectedMisc:      false,
			expectedGenFuncs:  false,
			expectedCalls:     false,
			fail:              false,
		},
		{
			name: "model with fields only",
			input: `model user {
  fields {
    name() string
    age() int
  }
}`,
			expectedMetadata:  nil,
			expectedModelName: "user",
			expectedFilepath:  "test.dg",
			expectedFields:    true,
			expectedMisc:      false,
			expectedGenFuncs:  false,
			expectedCalls:     false,
			fail:              false,
		},
		{
			name: "model with metadata count",
			input: `model test {
  metadata {
    count: 100
  }
}`,
			expectedMetadata: &codegen.Metadata{
				Count: 100,
				Tags:  map[string]string{},
			},
			expectedModelName: "test",
			expectedFilepath:  "test.dg",
			expectedFields:    false,
			expectedMisc:      false,
			expectedGenFuncs:  false,
			expectedCalls:     false,
			fail:              false,
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
			expectedMetadata: &codegen.Metadata{
				Tags: map[string]string{
					"service": "test",
					"team":    "backend",
				},
			},
			expectedModelName: "test",
			expectedFilepath:  "test.dg",
			expectedFields:    false,
			expectedMisc:      false,
			expectedGenFuncs:  false,
			expectedCalls:     false,
			fail:              false,
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
			expectedMetadata:  nil,
			expectedModelName: "test",
			expectedFilepath:  "test.dg",
			expectedFields:    false,
			expectedMisc:      true,
			expectedGenFuncs:  false,
			expectedCalls:     false,
			fail:              false,
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
			expectedMetadata:  nil,
			expectedModelName: "test",
			expectedFilepath:  "test.dg",
			expectedFields:    false,
			expectedMisc:      false,
			expectedGenFuncs:  true,
			expectedCalls:     false,
			fail:              false,
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
			expectedMetadata:  nil,
			expectedModelName: "test",
			expectedFilepath:  "test.dg",
			expectedFields:    true,
			expectedMisc:      false,
			expectedGenFuncs:  false,
			expectedCalls:     true,
			fail:              false,
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
			expectedMetadata:  nil,
			expectedModelName: "test",
			expectedFilepath:  "test.dg",
			expectedFields:    true,
			expectedMisc:      false,
			expectedGenFuncs:  false,
			expectedCalls:     false,
			fail:              false,
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
			expectedMetadata: &codegen.Metadata{
				Count: 50,
				Tags: map[string]string{
					"service": "test",
				},
			},
			expectedModelName: "test",
			expectedFilepath:  "test.dg",
			expectedFields:    false,
			expectedMisc:      false,
			expectedGenFuncs:  false,
			expectedCalls:     false,
			fail:              false,
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
			expectedMetadata: &codegen.Metadata{
				Count: 42,
				Tags: map[string]string{
					"key": "value",
				},
			},
			expectedModelName: "complete",
			expectedFilepath:  "test.dg",
			expectedFields:    true,
			expectedMisc:      true,
			expectedGenFuncs:  true,
			expectedCalls:     true,
			fail:              false,
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
			errStr: "expected '}', got '{'",
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
			errStr: "expected colon (:)",
		},
		{
			name:   "incomplete metadata tags",
			input:  "model test { metadata { tags } }",
			fail:   true,
			errStr: "expected colon (:)",
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
			errStr: "expected gen fn name",
		},
		{
			name:   "incomplete gen function",
			input:  "model test { gens { func name } }",
			fail:   true,
			errStr: "expected '(', got '}'",
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
			errStr: "invalid fields body incomplete body",
		},
		{
			name:   "mismatched braces - extra opening",
			input:  "model test { fields { {",
			fail:   true,
			errStr: "invalid fields body incomplete body",
		},
	}

	for _, tt := range testExpr {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.input), "test.dg")

			if tt.fail {
				require.Error(t, err, "expected error for input: %q", tt.input)
				assert.Contains(t, err.Error(), tt.errStr, "error message should contain %q, got: %v", tt.errStr, err)
				return
			}

			require.NoError(t, err, "unexpected error for input: %q", tt.input)
			require.NotNil(t, got, "expected non-nil result")

			assert.Equal(t, tt.expectedModelName, got.ModelName,
				"ModelName mismatch")
			assert.Equal(t, tt.expectedFilepath, got.Filepath,
				"Filepath mismatch")

			if tt.expectedMetadata != nil {
				require.NotNil(t, got.Metadata, "expected non-nil Metadata")
				assert.Equal(t, tt.expectedMetadata.Count, got.Metadata.Count,
					"Metadata.Count mismatch")
				if len(tt.expectedMetadata.Tags) > 0 {
					assert.Equal(t, tt.expectedMetadata.Tags, got.Metadata.Tags,
						"Metadata.Tags mismatch")
				}
			}

			assert.Equal(t, tt.expectedFields, got.Fields != nil, "Fields presence mismatch")
			assert.Equal(t, tt.expectedMisc, got.Misc != "", "Misc presence mismatch")
			assert.Equal(t, tt.expectedGenFuncs, len(got.GenFuns) > 0, "GenFuns presence mismatch")
			assert.Equal(t, tt.expectedCalls, len(got.Calls) > 0, "Calls presence mismatch")
		})
	}
}
