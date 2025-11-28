package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type wrapperFunc func(input string, wrapper func(string) string) (ast.Expr, error)

func parseWrappedExpr(input string, wrapper func(string) string) (ast.Expr, error) {
	wrappedCode := wrapper(input)
	expr, err := parser.ParseExpr(wrappedCode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %w", err)
	}
	return expr, nil
}

// parseFieldList parses a string containing field definitions into an *ast.FieldList.
// The input string should be in the format of Go interface methods.
func parseFieldList(input string, wrapperFunc wrapperFunc) (*ast.FieldList, error) {
	expr, err := wrapperFunc(input, func(s string) string {
		return fmt.Sprintf("interface {\n %s \n}", s)
	})
	if err != nil {
		return nil, err
	}

	// Get the interface type
	interfaceType, ok := expr.(*ast.InterfaceType)
	if !ok {
		return nil, fmt.Errorf("expected InterfaceType, got %T", expr)
	}

	return interfaceType.Methods, nil
}

// parseFunctionBlock parses a string containing arbitrary Go code into an *ast.BlockStmt.
// The input string should be valid Go code that can be wrapped in a function body.
func parseFunctionBlock(input string, wrapperFunc wrapperFunc) (*ast.BlockStmt, error) {
	expr, err := wrapperFunc(input, func(s string) string {
		return fmt.Sprintf("func() {\n%s\n}", s)
	})
	if err != nil {
		return nil, err
	}

	// Get the function literal
	funcLit, ok := expr.(*ast.FuncLit)
	if !ok {
		return nil, fmt.Errorf("expected FuncLit, got %T", expr)
	}

	// Get the function body
	if funcLit.Body == nil {
		return nil, fmt.Errorf("function body is nil")
	}

	return funcLit.Body, nil
}

// parseCallList parses a string containing a list of function calls into []*ast.CallExpr.
// The input string should be a list of function calls, one per line.
// Example:
//
//	id(1, 1)
//	created_at(time.Now(), time.Now())
func parseCallList(input string, wrapperFunc wrapperFunc) ([]*ast.CallExpr, error) {
	// Wrap the input in a block statement
	block, err := parseFunctionBlock(input, wrapperFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse call list: %w", err)
	}

	// Extract call expressions from the block
	return processCallBlock(block)
}

func processCallBlock(block *ast.BlockStmt) ([]*ast.CallExpr, error) {
	calls := make([]*ast.CallExpr, 0, len(block.List))
	for _, stmt := range block.List {
		// Each statement should be an expression statement
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			return nil, fmt.Errorf("expected ExprStmt, got %T", stmt)
		}

		// The expression should be a call expression
		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			return nil, fmt.Errorf("expected CallExpr, got %T", exprStmt.X)
		}

		calls = append(calls, call)
	}

	return calls, nil
}

// parseTags parses a string containing key-value pairs into map[string]string.
// The input string should be key-value pairs in JSON format (without outer braces).
// Example:
//
//	"service_name": "pluto",
//	"team_name": "platform"
func parseTags(input string, wrapperFunc wrapperFunc) (map[string]string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return map[string]string{}, nil
	}

	if strings.HasSuffix(trimmed, ",") {
		return nil, fmt.Errorf("failed to parse tags: trailing comma is not allowed")
	}

	expr, err := wrapperFunc(trimmed, func(s string) string {
		return fmt.Sprintf("map[string]string{\n%s,\n}", s)
	})
	if err != nil {
		return nil, err
	}

	return processTagExpr(expr)
}

func processTagExpr(expr ast.Expr) (map[string]string, error) {
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("expected CompositeLit, got %T", expr)
	}

	// Ensure the type is a map
	if _, ok := compLit.Type.(*ast.MapType); !ok {
		return nil, fmt.Errorf("expected map literal, got %T", compLit.Type)
	}

	result := make(map[string]string)
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			return nil, fmt.Errorf("expected key-value pair, got %T", elt)
		}

		// Keys and values must be string literals
		kLit, ok := kv.Key.(*ast.BasicLit)
		if !ok || kLit.Kind != token.STRING {
			return nil, fmt.Errorf("map key must be a string literal, got %T", kv.Key)
		}
		vLit, ok := kv.Value.(*ast.BasicLit)
		if !ok || vLit.Kind != token.STRING {
			return nil, fmt.Errorf("map value must be a string literal, got %T", kv.Value)
		}

		key, err := strconv.Unquote(kLit.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid string key: %w", err)
		}
		val, err := strconv.Unquote(vLit.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid string value for key %q: %w", key, err)
		}

		result[key] = val
	}

	return result, nil
}

// parseParamList parses a string containing a parameter list into *ast.CallExpr.
// The input string should be a Go parameter list.
// Examples:
//
//	''
//	start time.Time
//	a, b, c int
func parseParamList(input string, wrapperFunc wrapperFunc) (*ast.CallExpr, error) {
	// Wrap the input in a function type
	expr, err := wrapperFunc(input, func(s string) string {
		return fmt.Sprintf("func ( %s )", s)
	})
	if err != nil {
		return nil, err
	}

	// Get the function type
	funcType, ok := expr.(*ast.FuncType)
	if !ok {
		return nil, fmt.Errorf("expected FuncType, got %T", expr)
	}

	// Create a call expression with the function type
	return &ast.CallExpr{
		Fun: funcType,
	}, nil
}
