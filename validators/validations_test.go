package validators

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/ds-horizon/datagen/codegen"
	"github.com/ds-horizon/datagen/utils"
)

func TestValidate_AggregatesErrors(t *testing.T) {
	// Missing fields and gens + wrong filename
	d := &codegen.DatagenParsed{
		ModelName: "User",
		Filepath:  "dir" + utils.DgDirDelimeter + "NotUser",
	}

	err := Validate(d)
	if err == nil {
		t.Fatalf("expected aggregated error, got nil")
	}

	s := err.Error()
	if !strings.Contains(s, "model has no fields section") {
		t.Fatalf("expected missing fields message, got: %q", s)
	}
	if !strings.Contains(s, "model has no gens section") {
		t.Fatalf("expected missing gens message, got: %q", s)
	}
	if !strings.Contains(s, "model should be in file named User.dg, found in NotUser.dg") {
		t.Fatalf("expected model/filepath mismatch message, got: %q", s)
	}
}

func TestRequiredSectionsValidator(t *testing.T) {
	var errs MultiErr
	RequiredSectionsValidator(&codegen.DatagenParsed{}, &errs)
	if errs.Count() != 2 {
		t.Fatalf("expected 2 errors for missing sections, got %d", errs.Count())
	}
}

func TestNoDuplicateFieldNamesValidator(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "foo"}}, Type: &ast.Ident{Name: "int"}},
		{Names: []*ast.Ident{{Name: "foo"}}, Type: &ast.Ident{Name: "int"}},
	}}
	var errs MultiErr
	NoDuplicateFieldNamesValidator(&codegen.DatagenParsed{Fields: fields}, &errs)
	if errs.Count() != 1 {
		t.Fatalf("expected 1 error for duplicate fields, got %d", errs.Count())
	}
	if !strings.Contains(errs.Error(), "model has duplicate field names: foo") {
		t.Fatalf("unexpected error text: %q", errs.Error())
	}
}

func TestNoMissingGensValidator(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "a"}}},
		{Names: []*ast.Ident{{Name: "b"}}},
	}}
	d := &codegen.DatagenParsed{Fields: fields, GenFuns: []*codegen.GenFn{{Name: "a"}}}
	var errs MultiErr
	NoMissingGensValidator(d, &errs)
	if errs.Count() != 1 {
		t.Fatalf("expected 1 error for missing gen, got %d", errs.Count())
	}
	if !strings.Contains(errs.Error(), "model has missing gen functions: b") {
		t.Fatalf("unexpected error text: %q", errs.Error())
	}
}

func TestNoExtraGensValidator(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: "a"}}}}}
	d := &codegen.DatagenParsed{Fields: fields, GenFuns: []*codegen.GenFn{{Name: "a"}, {Name: "c"}}}
	var errs MultiErr
	NoExtraGensValidator(d, &errs)
	if errs.Count() != 1 {
		t.Fatalf("expected 1 error for extra gen, got %d", errs.Count())
	}
	if !strings.Contains(errs.Error(), "found extra gen functions: c") {
		t.Fatalf("unexpected error text: %q", errs.Error())
	}
}

func TestFilePathModelNameValidator(t *testing.T) {
	var errs MultiErr
	// ok case
	good := &codegen.DatagenParsed{ModelName: "good", Filepath: "x" + utils.DgDirDelimeter + "good"}
	FilePathModelNameValidator(good, &errs)
	if errs.Count() != 0 {
		t.Fatalf("expected no error for matching model/file, got %d", errs.Count())
	}
	// bad case
	bad := &codegen.DatagenParsed{ModelName: "good", Filepath: "x" + utils.DgDirDelimeter + "bad"}
	FilePathModelNameValidator(bad, &errs)
	if errs.Count() != 1 {
		t.Fatalf("expected 1 error for mismatched filename, got %d", errs.Count())
	}
	if !strings.Contains(errs.Error(), "model should be in file named good.dg, found in bad.dg") {
		t.Fatalf("unexpected error text: %q", errs.Error())
	}
}

func TestCallExprsValidator_MissingReturnType(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "f"}}, Type: &ast.FuncType{Params: &ast.FieldList{}, Results: nil}},
	}}
	d := &codegen.DatagenParsed{Fields: fields}
	var errs MultiErr
	CallExprsValidator(d, &errs)
	if errs.Count() != 1 || !strings.Contains(errs.Error(), "field f must declare a return type") {
		t.Fatalf("expected missing return type error, got: %q", errs.Error())
	}
}

func TestCallExprsValidator_ZeroParamsSkipsCallRequirement(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "z"}}, Type: &ast.FuncType{Params: &ast.FieldList{}, Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}}}}},
	}}
	d := &codegen.DatagenParsed{Fields: fields}
	var errs MultiErr
	CallExprsValidator(d, &errs)
	if errs.Count() != 0 {
		t.Fatalf("expected no error (zero-parameter function), got: %q", errs.Error())
	}
}

func TestCallExprsValidator_ParamCountMismatchAndUnknownCall(t *testing.T) {
	// Field g expects 2 params
	ft := &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "a"}}, Type: &ast.Ident{Name: "int"}},
		{Names: []*ast.Ident{{Name: "b"}}, Type: &ast.Ident{Name: "int"}},
	}}, Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}}}}
	fields := &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: "g"}}, Type: ft}}}
	// Provide one arg for g (mismatch) and an extra call h (unknown)
	calls := []*ast.CallExpr{
		{Fun: &ast.Ident{Name: "g"}, Args: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "1"}}},
		{Fun: &ast.Ident{Name: "h"}, Args: []ast.Expr{}},
	}
	d := &codegen.DatagenParsed{Fields: fields, Calls: calls}
	var errs MultiErr
	CallExprsValidator(d, &errs)
	msg := errs.Error()
	if !strings.Contains(msg, "field g expects 2 args, got 1") {
		t.Fatalf("expected arg mismatch error, got: %q", msg)
	}
	if !strings.Contains(msg, "unknown call h") {
		t.Fatalf("expected unknown call error, got: %q", msg)
	}
}

func TestGenFnsReturnValidator(t *testing.T) {
	// nil gen
	d := &codegen.DatagenParsed{GenFuns: []*codegen.GenFn{nil}}
	var errs MultiErr
	GenFnsReturnValidator(d, &errs)
	// nil body
	d.GenFuns = []*codegen.GenFn{{Name: "A", Body: nil}}
	GenFnsReturnValidator(d, &errs)
	// body without return
	bodyNoReturn := &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{X: &ast.BasicLit{Kind: token.STRING, Value: "\"x\""}}}}
	d.GenFuns = []*codegen.GenFn{{Name: "B", Body: bodyNoReturn}}
	GenFnsReturnValidator(d, &errs)
	// body with return
	bodyWithReturn := &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}}
	d.GenFuns = []*codegen.GenFn{{Name: "C", Body: bodyWithReturn}}
	GenFnsReturnValidator(d, &errs)

	msg := errs.Error()
	if !strings.Contains(msg, "gen func <unknown> must return a value") {
		t.Fatalf("expected error for nil gen, got: %q", msg)
	}
	if !strings.Contains(msg, "gen func A must return a value") {
		t.Fatalf("expected error for nil body, got: %q", msg)
	}
	if !strings.Contains(msg, "gen func B must return a value") {
		t.Fatalf("expected error for body without return, got: %q", msg)
	}
}
