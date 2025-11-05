package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testWrapperExpr = func(t *testing.T, expected string, wantErr error, wantRes ast.Expr) wrapperFunc {
	return func(input string, wrapper func(string) string) (ast.Expr, error) {
		wrappedString := wrapper(input)
		assert.Equal(t, expected, wrappedString)
		return wantRes, wantErr
	}
}

func TestParseFieldList(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		returnErr   error
		returnRes   ast.Expr
		expectedErr error
	}{
		{
			name: "simple methods",
			input: `priority() int
name() string
description() string`,
			expected: `interface {
 priority() int
name() string
description() string 
}`,
			returnErr:   nil,
			returnRes:   &ast.InterfaceType{},
			expectedErr: nil,
		},
		{
			name: "methods with parameters",
			input: `id(start int, step int) int
created_at(startDate time.Time, endDate time.Time) time.Time`,
			expected: `interface {
 id(start int, step int) int
created_at(startDate time.Time, endDate time.Time) time.Time 
}`,
			returnErr:   nil,
			returnRes:   &ast.InterfaceType{},
			expectedErr: nil,
		},
		{
			name:  "empty input",
			input: "",
			expected: `interface {
  
}`,
			returnErr:   nil,
			returnRes:   &ast.InterfaceType{},
			expectedErr: nil,
		},
		{
			name:  "expect error is returned as is",
			input: "foo",
			expected: `interface {
 foo 
}`,
			returnErr:   errors.New("error"),
			returnRes:   &ast.InterfaceType{},
			expectedErr: errors.New("error"),
		},
		{
			name:  "invalid input",
			input: "++",
			expected: `interface {
 ++ 
}`,
			returnErr:   nil,
			returnRes:   nil,
			expectedErr: fmt.Errorf("expected InterfaceType, got %T", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseFieldList(tt.input, testWrapperExpr(t, tt.expected, tt.returnErr, tt.returnRes))
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestParseFunctionBlock(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		returnErr   error
		returnRes   ast.Expr
		expectedErr error
	}{
		{
			name:  "simple return",
			input: "return 42",
			expected: `func() {
return 42
}`,
			returnErr:   nil,
			returnRes:   &ast.FuncLit{Body: &ast.BlockStmt{List: []ast.Stmt{}}},
			expectedErr: nil,
		},
		{
			name: "methods with parameters",
			input: `x := 10
return x * 2`,
			expected: `func() {
x := 10
return x * 2
}`,
			returnErr:   nil,
			returnRes:   &ast.FuncLit{Body: &ast.BlockStmt{List: []ast.Stmt{}}},
			expectedErr: nil,
		},
		{
			name:  "empty input",
			input: "",
			expected: `func() {

}`,
			returnErr:   nil,
			returnRes:   &ast.FuncLit{Body: &ast.BlockStmt{List: []ast.Stmt{}}},
			expectedErr: nil,
		},
		{
			name:  "error in wrapperFunc",
			input: "foo",
			expected: `func() {
foo
}`,
			returnErr:   errors.New("error"),
			returnRes:   &ast.FuncLit{Body: &ast.BlockStmt{List: []ast.Stmt{}}},
			expectedErr: errors.New("error"),
		},
		{
			name:  "invalid input",
			input: "++",
			expected: `func() {
++
}`,
			returnErr:   nil,
			returnRes:   nil,
			expectedErr: fmt.Errorf("expected FuncLit, got %T", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseFunctionBlock(tt.input, testWrapperExpr(t, tt.expected, tt.returnErr, tt.returnRes))
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestParseParamList(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		returnErr   error
		returnRes   ast.Expr
		expectedErr error
	}{
		{
			name:        "empty input",
			input:       "",
			expected:    `func (  )`,
			returnErr:   nil,
			returnRes:   &ast.FuncType{},
			expectedErr: nil,
		},
		{
			name:        "simple arg",
			input:       "a int",
			expected:    `func ( a int )`,
			returnErr:   nil,
			returnRes:   &ast.FuncType{},
			expectedErr: nil,
		},
		{
			name:        "multiple args",
			input:       `start, stop time.Time, a int, b int`,
			expected:    `func ( start, stop time.Time, a int, b int )`,
			returnErr:   nil,
			returnRes:   &ast.FuncType{},
			expectedErr: nil,
		},
		{
			name:        "error in wrapperFunc",
			input:       "foo",
			expected:    `func ( foo )`,
			returnErr:   errors.New("error"),
			returnRes:   &ast.FuncType{},
			expectedErr: errors.New("error"),
		},
		{
			name:        "invalid input",
			input:       "++",
			expected:    `func ( ++ )`,
			returnErr:   nil,
			returnRes:   nil,
			expectedErr: fmt.Errorf("expected FuncType, got %T", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := parseParamList(tt.input, testWrapperExpr(t, tt.expected, tt.returnErr, tt.returnRes))
			assert.Equal(t, tt.expectedErr, err)
			if err == nil {
				assert.Equal(t, &ast.CallExpr{Fun: tt.returnRes}, val)
			}
		})
	}
}

func TestProcessCallBlock(t *testing.T) {
	mkCall := func(name string, args ...ast.Expr) *ast.CallExpr {
		return &ast.CallExpr{Fun: ast.NewIdent(name), Args: args}
	}

	call1 := mkCall("foo")
	call2 := mkCall("bar", ast.NewIdent("x"))

	tests := []struct {
		name           string
		block          *ast.BlockStmt
		wantCallsLen   int
		wantErrSubstrs []string
		wantPanic      bool
	}{
		{
			name: "success-two-calls",
			block: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: call1},
					&ast.ExprStmt{X: call2},
				},
			},
			wantCallsLen: 2,
		},
		{
			name:         "empty-block",
			block:        &ast.BlockStmt{List: nil},
			wantCallsLen: 0,
		},
		{
			name: "non-expr-stmt",
			block: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{},
				},
			},
			wantErrSubstrs: []string{"expected ExprStmt", "ReturnStmt"},
		},
		{
			name: "exprstmt-with-non-callexpr",
			block: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: ast.NewIdent("notACall")},
				},
			},
			wantErrSubstrs: []string{"expected CallExpr", "Ident"},
		},
		{
			name:      "nil-block-panics",
			block:     nil,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					_, _ = processCallBlock(tt.block)
				}, "expected panic for input: %+v", tt.block)
				return
			}

			got, err := processCallBlock(tt.block)

			if len(tt.wantErrSubstrs) > 0 {
				if assert.Error(t, err, "expected an error") {
					for _, sub := range tt.wantErrSubstrs {
						assert.Contains(t, err.Error(), sub, "error should contain substring %q", sub)
					}
				}
				assert.Nil(t, got, "expected nil result on error")
				return
			}

			assert.NoError(t, err, "unexpected error")
			assert.NotNil(t, got, "expected non-nil slice")
			assert.Len(t, got, tt.wantCallsLen, "unexpected number of calls")
		})
	}
}

func TestParseTags(t *testing.T) {
	// empty string is handled
	got, err := parseTags("", testWrapperExpr(t, "", nil, nil))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{}, got)

	// only whitespace string is handled
	got, err = parseTags("   \t \n", testWrapperExpr(t, "", nil, nil))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{}, got)

	// trailing comma is handled
	got, err = parseTags(`"k": "v", `, testWrapperExpr(t, "", nil, nil))
	require.Error(t, err)
	require.Nil(t, got)
	assert.Contains(t, err.Error(), "trailing comma is not allowed")

	// invalid data
	got, err = parseTags(`++`, testWrapperExpr(t, `map[string]string{
++,
}`, errors.New("invalid map"), nil))
	require.Error(t, err)
	assert.Nil(t, got)

	// happy path
	got, err = parseTags(`"a": "b"`, testWrapperExpr(t, `map[string]string{
"a": "b",
}`, nil, mustParseExpr(t, `map[string]string{"a":"b"}`)))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "b"}, got)
}

func mustParseExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	expr, err := parser.ParseExpr(src)
	if err != nil {
		t.Fatalf("parse error for %q: %v", src, err)
	}
	return expr
}

func TestProcessTagExpr(t *testing.T) {
	type tc struct {
		name      string
		expr      ast.Expr
		want      map[string]string
		wantErr   bool
		errSubstr string
	}

	tests := []tc{
		{
			name: "simple map",
			expr: mustParseExpr(t, `map[string]string{"a":"1","b":"2"}`),
			want: map[string]string{"a": "1", "b": "2"},
		},
		{
			name: "simple map with url",
			expr: mustParseExpr(t, `map[string]string{"a":"http://example.com"}`),
			want: map[string]string{"a": "http://example.com"},
		},
		{
			name: "empty map",
			expr: mustParseExpr(t, `map[string]string{}`),
			want: map[string]string{},
		},
		{
			name: "escapes are unquoted",
			expr: mustParseExpr(t, `map[string]string{"key\"with\"quotes":"line\nbreak"}`),
			want: map[string]string{`key"with"quotes`: "line\nbreak"},
		},
		{
			name:      "not a CompositeLit",
			expr:      mustParseExpr(t, `someIdent`),
			wantErr:   true,
			errSubstr: "expected CompositeLit",
		},
		{
			name:      "composite is not a map",
			expr:      mustParseExpr(t, `[]string{"a"}`),
			wantErr:   true,
			errSubstr: "expected map literal",
		},
		{
			name:      "key is not a string literal",
			expr:      mustParseExpr(t, `map[string]string{1:"x"}`),
			wantErr:   true,
			errSubstr: "map key must be a string literal",
		},
		{
			name:      "value is not a string literal",
			expr:      mustParseExpr(t, `map[string]string{"a": 1}`),
			wantErr:   true,
			errSubstr: "map value must be a string literal",
		},
		{
			name: "duplicate keys keep last value",
			expr: mustParseExpr(t, `map[string]string{"a":"first","b":"ok","a":"last"}`),
			want: map[string]string{"a": "last", "b": "ok"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := processTagExpr(tt.expr)

			if tt.wantErr {
				require.Error(t, err, "expected error but got none")

				if tt.errSubstr != "" {
					assert.Contains(t, err.Error(), tt.errSubstr,
						"error message should contain expected substring")
				}
				return
			}

			require.NoError(t, err, "unexpected error")
			assert.Equal(t, tt.want, got, "result mismatch")
		})
	}
}

func TestParseWrappedExpr(t *testing.T) {
	// invalid go code
	wrapper := func(string) string { return "this is not valid go code" }

	expr, err := parseWrappedExpr("ignored", wrapper)
	require.Error(t, err)
	assert.Nil(t, expr)
	assert.Contains(t, err.Error(), "failed to parse expression")

	// happy path
	wrapper = func(s string) string { return "(" + s + ")" }
	expr, err = parseWrappedExpr("1 + 2", wrapper)
	require.NoError(t, err)
	require.NotNil(t, expr)
	require.IsType(t, &ast.ParenExpr{}, expr)

	pe := expr.(*ast.ParenExpr)
	be, ok := pe.X.(*ast.BinaryExpr)
	require.True(t, ok, "inner expr should be a BinaryExpr")
	assert.Equal(t, token.ADD, be.Op)
}
