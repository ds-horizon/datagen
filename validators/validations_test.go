package validators

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/dream-horizon-org/datagen/codegen"
	"github.com/dream-horizon-org/datagen/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidate_AggregatesErrors(t *testing.T) {
	// Missing fields and gens + wrong filename
	d := &codegen.DatagenParsed{
		ModelName: "User",
		Filepath:  "dir" + utils.DgDirDelimeter + "NotUser",
	}

	err := Validate(d)
	assert.Error(t, err, "expected aggregated error")

	s := err.Error()
	assert.Contains(t, s, "model has no fields section")
	assert.Contains(t, s, "model has no gens section")
	assert.Contains(t, s, "model should be in file named User.dg, found in NotUser.dg")
}

func TestRequiredSectionsValidator(t *testing.T) {
	var errs MultiErr
	RequiredSectionsValidator(&codegen.DatagenParsed{}, &errs)
	assert.Equal(t, 2, errs.Count(), "expected 2 errors for missing sections")
	msg := errs.Error()
	assert.Contains(t, msg, "model has no fields section")
	assert.Contains(t, msg, "model has no gens section")
}

func TestNoDuplicateFieldNamesValidator(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "foo"}}, Type: &ast.Ident{Name: "int"}},
		{Names: []*ast.Ident{{Name: "foo"}}, Type: &ast.Ident{Name: "int"}},
	}}
	var errs MultiErr
	NoDuplicateFieldNamesValidator(&codegen.DatagenParsed{Fields: fields}, &errs)
	assert.Equal(t, 1, errs.Count(), "expected 1 error for duplicate fields")
	assert.Contains(t, errs.Error(), "model has duplicate field names: foo")
}

func TestNoMissingGensValidator(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "a"}}},
		{Names: []*ast.Ident{{Name: "b"}}},
	}}
	d := &codegen.DatagenParsed{Fields: fields, GenFuns: []*codegen.GenFn{{Name: "a"}}}
	var errs MultiErr
	NoMissingGensValidator(d, &errs)
	assert.Equal(t, 1, errs.Count(), "expected 1 error for missing gen")
	assert.Contains(t, errs.Error(), "model has missing gen functions: b")
}

func TestNoExtraGensValidator(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: "a"}}}}}
	d := &codegen.DatagenParsed{Fields: fields, GenFuns: []*codegen.GenFn{{Name: "a"}, {Name: "c"}}}
	var errs MultiErr
	NoExtraGensValidator(d, &errs)
	assert.Equal(t, 1, errs.Count(), "expected 1 error for extra gen")
	assert.Contains(t, errs.Error(), "found extra gen functions: c")
}

func TestFilePathModelNameValidator(t *testing.T) {
	var errs MultiErr
	// ok case
	good := &codegen.DatagenParsed{ModelName: "good", Filepath: "x" + utils.DgDirDelimeter + "good"}
	FilePathModelNameValidator(good, &errs)
	assert.Equal(t, 0, errs.Count(), "expected no error for matching model/file")
	// bad case
	bad := &codegen.DatagenParsed{ModelName: "good", Filepath: "x" + utils.DgDirDelimeter + "bad"}
	FilePathModelNameValidator(bad, &errs)
	assert.Equal(t, 1, errs.Count(), "expected 1 error for mismatched filename")
	assert.Contains(t, errs.Error(), "model should be in file named good.dg, found in bad.dg")
}

func TestCallExprsValidator_MissingReturnType(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "f"}}, Type: &ast.FuncType{Params: &ast.FieldList{}, Results: nil}},
	}}
	d := &codegen.DatagenParsed{Fields: fields}
	var errs MultiErr
	CallExprsValidator(d, &errs)
	assert.Equal(t, 1, errs.Count(), "expected one error for missing return type")
	assert.Contains(t, errs.Error(), "field f must declare a return type")
}

func TestCallExprsValidator_ZeroParamsSkipsCallRequirement(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "z"}}, Type: &ast.FuncType{Params: &ast.FieldList{}, Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}}}}},
	}}
	d := &codegen.DatagenParsed{Fields: fields}
	var errs MultiErr
	CallExprsValidator(d, &errs)
	assert.Equal(t, 0, errs.Count(), "expected no error (zero-parameter function)")
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
	assert.Contains(t, msg, "field g expects 2 args, got 1")
	assert.Contains(t, msg, "unknown call h")
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
	assert.Contains(t, msg, "gen func <unknown> must return a value")
	assert.Contains(t, msg, "gen func A must return a value")
	assert.Contains(t, msg, "gen func B must return a value")
}
