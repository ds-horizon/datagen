package parser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"strings"
)

// parseWrappedExpr wraps the input in a specific syntax and parses it as an expression
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
func parseFieldList(input string) (*ast.FieldList, error) {
	expr, err := parseWrappedExpr(input, func(s string) string {
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
func parseFunctionBlock(input string) (*ast.BlockStmt, error) {
	expr, err := parseWrappedExpr(input, func(s string) string {
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
func parseCallList(input string) ([]*ast.CallExpr, error) {
	// Wrap the input in a block statement
	block, err := parseFunctionBlock(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse call list: %w", err)
	}

	// Extract call expressions from the block
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
func parseTags(input string) (map[string]string, error) {
	cleaned := removeLineComments(input)
	trimmed := strings.TrimSpace(cleaned)

	if trimmed == "" {
		return map[string]string{}, nil
	}

	// Disallow a trailing comma in the tags section of the DSL.
	if strings.HasSuffix(trimmed, ",") {
		return nil, fmt.Errorf("failed to parse tags: trailing comma is not allowed")
	}

	payload := "{" + trimmed + "}"
	var out map[string]string
	if err := json.Unmarshal([]byte(payload), &out); err != nil {
		return nil, fmt.Errorf("failed to parse tags: %w", err)
	}
	return out, nil
}

// removeLineComments strips Go-style line comments (// ...) from the input.
// This is used to reliably detect a trailing comma even when a comment follows it.
func removeLineComments(s string) string {
	// Fast path: if there's no comment marker, return as-is
	if !strings.Contains(s, "//") {
		return s
	}
	var b strings.Builder
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}
		b.WriteString(line)
		// Preserve original line structure
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// parseParamList parses a string containing a parameter list into *ast.CallExpr.
// The input string should be a Go parameter list.
// Examples:
//
//	''
//	start time.Time
//	a, b, c int
func parseParamList(input string) (*ast.CallExpr, error) {
	// Wrap the input in a function type
	expr, err := parseWrappedExpr(input, func(s string) string {
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
