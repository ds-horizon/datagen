package parser

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

// compareFieldList compares only the essential parts of two FieldLists
func compareFieldList(t *testing.T, expected, actual *ast.FieldList) {
	assert.Equal(t, len(expected.List), len(actual.List), "number of fields should match")

	for i, expectedField := range expected.List {
		actualField := actual.List[i]

		// Compare field names
		assert.Equal(t, len(expectedField.Names), len(actualField.Names), "number of names should match")
		for j, expectedName := range expectedField.Names {
			assert.Equal(t, expectedName.Name, actualField.Names[j].Name, "field name should match")
		}

		// Compare function types
		expectedFuncType, ok1 := expectedField.Type.(*ast.FuncType)
		actualFuncType, ok2 := actualField.Type.(*ast.FuncType)
		assert.True(t, ok1 && ok2, "both types should be FuncType")

		// Compare parameters
		if expectedFuncType.Params != nil {
			assert.Equal(t, len(expectedFuncType.Params.List), len(actualFuncType.Params.List), "number of parameters should match")
			for j, expectedParam := range expectedFuncType.Params.List {
				actualParam := actualFuncType.Params.List[j]
				assert.Equal(t, len(expectedParam.Names), len(actualParam.Names), "number of parameter names should match")
				for k, expectedName := range expectedParam.Names {
					assert.Equal(t, expectedName.Name, actualParam.Names[k].Name, "parameter name should match")
				}
				compareType(t, expectedParam.Type, actualParam.Type)
			}
		}

		// Compare results
		if expectedFuncType.Results != nil {
			assert.Equal(t, len(expectedFuncType.Results.List), len(actualFuncType.Results.List), "number of results should match")
			for j, expectedResult := range expectedFuncType.Results.List {
				actualResult := actualFuncType.Results.List[j]
				compareType(t, expectedResult.Type, actualResult.Type)
			}
		}
	}
}

// compareType compares two AST types
func compareType(t *testing.T, expected, actual ast.Expr) {
	t.Logf("Comparing types: expected=%T, actual=%T", expected, actual)

	if expected == nil {
		assert.Nil(t, actual, "expected nil type")
		return
	}
	assert.NotNil(t, actual, "actual type should not be nil")

	switch exp := expected.(type) {
	case *ast.Ident:
		act, ok := actual.(*ast.Ident)
		assert.True(t, ok, "expected Ident type")
		assert.Equal(t, exp.Name, act.Name, "type name should match")
	case *ast.SelectorExpr:
		act, ok := actual.(*ast.SelectorExpr)
		assert.True(t, ok, "expected SelectorExpr type")
		compareType(t, exp.X, act.X)
		compareType(t, exp.Sel, act.Sel)
	case *ast.ArrayType:
		act, ok := actual.(*ast.ArrayType)
		assert.True(t, ok, "expected ArrayType")
		compareType(t, exp.Elt, act.Elt)
	case *ast.MapType:
		act, ok := actual.(*ast.MapType)
		assert.True(t, ok, "expected MapType")
		compareType(t, exp.Key, act.Key)
		compareType(t, exp.Value, act.Value)
	case *ast.ChanType:
		act, ok := actual.(*ast.ChanType)
		assert.True(t, ok, "expected ChanType")
		assert.Equal(t, exp.Dir, act.Dir, "channel direction should match")
		compareType(t, exp.Value, act.Value)
	case *ast.StarExpr:
		act, ok := actual.(*ast.StarExpr)
		assert.True(t, ok, "expected StarExpr")
		compareType(t, exp.X, act.X)
	default:
		t.Fatalf("unexpected type: %T", expected)
	}
}

func TestParseFieldList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.FieldList
	}{
		{
			name: "simple methods",
			input: `priority() int
name() string
description() string`,
			expected: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "priority"}},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
							Results: &ast.FieldList{
								List: []*ast.Field{
									{Type: &ast.Ident{Name: "int"}},
								},
							},
						},
					},
					{
						Names: []*ast.Ident{{Name: "name"}},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
							Results: &ast.FieldList{
								List: []*ast.Field{
									{Type: &ast.Ident{Name: "string"}},
								},
							},
						},
					},
					{
						Names: []*ast.Ident{{Name: "description"}},
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
							Results: &ast.FieldList{
								List: []*ast.Field{
									{Type: &ast.Ident{Name: "string"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "methods with parameters",
			input: `id(start int, step int) int
created_at(startDate time.Time, endDate time.Time) time.Time`,
			expected: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "id"}},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{{Name: "start"}},
										Type:  &ast.Ident{Name: "int"},
									},
									{
										Names: []*ast.Ident{{Name: "step"}},
										Type:  &ast.Ident{Name: "int"},
									},
								},
							},
							Results: &ast.FieldList{
								List: []*ast.Field{
									{Type: &ast.Ident{Name: "int"}},
								},
							},
						},
					},
					{
						Names: []*ast.Ident{{Name: "created_at"}},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{{Name: "startDate"}},
										Type: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "time"},
											Sel: &ast.Ident{Name: "Time"},
										},
									},
									{
										Names: []*ast.Ident{{Name: "endDate"}},
										Type: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "time"},
											Sel: &ast.Ident{Name: "Time"},
										},
									},
								},
							},
							Results: &ast.FieldList{
								List: []*ast.Field{
									{
										Type: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "time"},
											Sel: &ast.Ident{Name: "Time"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "empty input",
			input: "",
			expected: &ast.FieldList{
				List: []*ast.Field{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldList, err := parseFieldList(tt.input)
			assert.NoError(t, err)
			compareFieldList(t, tt.expected, fieldList)
		})
	}
}

func TestParseFunctionBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.BlockStmt
	}{
		{
			name:  "simple return",
			input: "return 42",
			expected: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.BasicLit{Value: "42"},
						},
					},
				},
			},
		},
		{
			name:  "return with binary expression",
			input: "return start + step * iter",
			expected: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.BinaryExpr{
								X:  &ast.Ident{Name: "start"},
								Op: token.ADD,
								Y: &ast.BinaryExpr{
									X:  &ast.Ident{Name: "step"},
									Op: token.MUL,
									Y:  &ast.Ident{Name: "iter"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "return with method call",
			input: "return self.id(iter)",
			expected: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "self"},
									Sel: &ast.Ident{Name: "id"},
								},
								Args: []ast.Expr{
									&ast.Ident{Name: "iter"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "multiple statements",
			input: "x := 10\nreturn x * 2",
			expected: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{&ast.BasicLit{Value: "10"}},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.BinaryExpr{
								X:  &ast.Ident{Name: "x"},
								Op: token.MUL,
								Y:  &ast.BasicLit{Value: "2"},
							},
						},
					},
				},
			},
		},
		{
			name:  "empty block",
			input: "",
			expected: &ast.BlockStmt{
				List: []ast.Stmt{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := parseFunctionBlock(tt.input)
			assert.NoError(t, err)
			compareBlock(t, tt.expected, block)
		})
	}
}

// compareBlock compares two BlockStmts
func compareBlock(t *testing.T, expected, actual *ast.BlockStmt) {
	assert.Equal(t, len(expected.List), len(actual.List), "number of statements should match")

	for i, expectedStmt := range expected.List {
		actualStmt := actual.List[i]

		switch exp := expectedStmt.(type) {
		case *ast.ReturnStmt:
			act, ok := actualStmt.(*ast.ReturnStmt)
			assert.True(t, ok, "expected ReturnStmt")
			assert.Equal(t, len(exp.Results), len(act.Results), "number of return values should match")
			for j, expectedResult := range exp.Results {
				compareExpr(t, expectedResult, act.Results[j])
			}
		case *ast.AssignStmt:
			act, ok := actualStmt.(*ast.AssignStmt)
			assert.True(t, ok, "expected AssignStmt")
			assert.Equal(t, exp.Tok, act.Tok, "assignment operator should match")
			assert.Equal(t, len(exp.Lhs), len(act.Lhs), "number of LHS expressions should match")
			assert.Equal(t, len(exp.Rhs), len(act.Rhs), "number of RHS expressions should match")
			for j, expectedLhs := range exp.Lhs {
				compareExpr(t, expectedLhs, act.Lhs[j])
			}
			for j, expectedRhs := range exp.Rhs {
				compareExpr(t, expectedRhs, act.Rhs[j])
			}
		default:
			t.Fatalf("unexpected statement type: %T", expectedStmt)
		}
	}
}

// compareExpr compares two AST expressions
func compareExpr(t *testing.T, expected, actual ast.Expr) {
	switch exp := expected.(type) {
	case *ast.Ident:
		act, ok := actual.(*ast.Ident)
		assert.True(t, ok, "expected Ident")
		assert.Equal(t, exp.Name, act.Name, "identifier name should match")
	case *ast.BasicLit:
		act, ok := actual.(*ast.BasicLit)
		assert.True(t, ok, "expected BasicLit")
		assert.Equal(t, exp.Value, act.Value, "literal value should match")
	case *ast.BinaryExpr:
		act, ok := actual.(*ast.BinaryExpr)
		assert.True(t, ok, "expected BinaryExpr")
		assert.Equal(t, exp.Op, act.Op, "binary operator should match")
		compareExpr(t, exp.X, act.X)
		compareExpr(t, exp.Y, act.Y)
	case *ast.CallExpr:
		act, ok := actual.(*ast.CallExpr)
		assert.True(t, ok, "expected CallExpr")
		compareExpr(t, exp.Fun, act.Fun)
		assert.Equal(t, len(exp.Args), len(act.Args), "number of arguments should match")
		for i, expectedArg := range exp.Args {
			compareExpr(t, expectedArg, act.Args[i])
		}
	case *ast.SelectorExpr:
		act, ok := actual.(*ast.SelectorExpr)
		assert.True(t, ok, "expected SelectorExpr")
		compareExpr(t, exp.X, act.X)
		compareExpr(t, exp.Sel, act.Sel)
	default:
		t.Fatalf("unexpected expression type: %T", expected)
	}
}

func TestParseCallList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []*ast.CallExpr
	}{
		{
			name: "simple function calls",
			input: `id(1, 1)
created_at(time.Now(), time.Now())`,
			expected: []*ast.CallExpr{
				{
					Fun: &ast.Ident{Name: "id"},
					Args: []ast.Expr{
						&ast.BasicLit{Value: "1"},
						&ast.BasicLit{Value: "1"},
					},
				},
				{
					Fun: &ast.Ident{Name: "created_at"},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "time"},
								Sel: &ast.Ident{Name: "Now"},
							},
						},
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "time"},
								Sel: &ast.Ident{Name: "Now"},
							},
						},
					},
				},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []*ast.CallExpr{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls, err := parseCallList(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(calls), "number of calls should match")

			for i, expectedCall := range tt.expected {
				actualCall := calls[i]
				compareCallExpr(t, expectedCall, actualCall)
			}
		})
	}
}

// compareCallExpr compares two CallExprs
func compareCallExpr(t *testing.T, expected, actual *ast.CallExpr) {
    // Compare the function being called. It can be a FuncType (for param list parsing)
    // or a regular callable expression like Ident/Selector (for call list parsing).
    switch expectedFun := expected.Fun.(type) {
    case *ast.FuncType:
        actualFun, ok := actual.Fun.(*ast.FuncType)
        assert.True(t, ok, "actual Fun should be FuncType")

        if expectedFun.Params != nil {
            assert.NotNil(t, actualFun.Params, "actual params should not be nil")
            assert.Equal(t, len(expectedFun.Params.List), len(actualFun.Params.List), "number of parameters should match")

            for i, expectedField := range expectedFun.Params.List {
                actualField := actualFun.Params.List[i]

                // Compare field names
                assert.Equal(t, len(expectedField.Names), len(actualField.Names), "number of names should match")
                for j, expectedName := range expectedField.Names {
                    assert.Equal(t, expectedName.Name, actualField.Names[j].Name, "parameter name should match")
                }

                // Compare types
                compareType(t, expectedField.Type, actualField.Type)
            }
        } else {
            assert.Nil(t, actualFun.Params, "actual params should be nil")
        }
    default:
        // For non-FuncType call targets, compare as general expressions
        compareExpr(t, expected.Fun, actual.Fun)
    }

    // Compare arguments if any
    assert.Equal(t, len(expected.Args), len(actual.Args), "number of arguments should match")
    for i, expectedArg := range expected.Args {
        compareExpr(t, expectedArg, actual.Args[i])
    }
}

func TestParseParamList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.CallExpr
	}{
		{
			name:  "empty parameters",
			input: "",
			expected: &ast.CallExpr{
				Fun: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{},
					},
				},
			},
		},
		{
			name:  "single parameter",
			input: "start time.Time",
			expected: &ast.CallExpr{
				Fun: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "start"}},
								Type: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "time"},
									Sel: &ast.Ident{Name: "Time"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "multiple parameters with same type",
			input: "a, b, c int",
			expected: &ast.CallExpr{
				Fun: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "a"},
									{Name: "b"},
									{Name: "c"},
								},
								Type: &ast.Ident{Name: "int"},
							},
						},
					},
				},
			},
		},
		{
			name:  "multiple parameters with different types",
			input: "id int, name string, createdAt time.Time",
			expected: &ast.CallExpr{
				Fun: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "id"}},
								Type:  &ast.Ident{Name: "int"},
							},
							{
								Names: []*ast.Ident{{Name: "name"}},
								Type:  &ast.Ident{Name: "string"},
							},
							{
								Names: []*ast.Ident{{Name: "createdAt"}},
								Type: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "time"},
									Sel: &ast.Ident{Name: "Time"},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call, err := parseParamList(tt.input)
			assert.NoError(t, err)
			compareCallExpr(t, tt.expected, call)
		})
	}
}
