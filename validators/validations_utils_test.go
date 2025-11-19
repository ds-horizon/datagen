package validators

import (
	"testing"

	"go/ast"
	"go/token"

	"github.com/ds-horizon/datagen/codegen"
)

func TestBuildGenSet(t *testing.T) {
	d := &codegen.DatagenParsed{GenFuns: []*codegen.GenFn{{Name: "a"}, {Name: "  "}, nil}}
	set := buildGenSet(d)
	if _, ok := set["a"]; !ok || len(set) != 1 {
		t.Fatalf("expected set to contain only 'a', got %#v", set)
	}
}

func TestFetchMissingGens(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: "a"}}}, {Names: []*ast.Ident{{Name: "b"}}}}}
	d := &codegen.DatagenParsed{Fields: fields}
	set := map[string]struct{}{"a": {}}
	missing := fetchMissingGens(d, set)
	if len(missing) != 1 || missing[0] != "b" {
		t.Fatalf("expected missing ['b'], got %#v", missing)
	}
}

func TestFetchDuplicateAndFieldSet(t *testing.T) {
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "x"}}},
		{Names: []*ast.Ident{{Name: "y"}}},
		{Names: []*ast.Ident{{Name: "x"}}},
	}}
	d := &codegen.DatagenParsed{Fields: fields}
	dups, set := fetchDuplicateAndFieldSet(d)
	if len(dups) != 1 || dups[0] != "x" {
		t.Fatalf("expected duplicates ['x'], got %#v", dups)
	}
	if _, ok := set["x"]; !ok {
		t.Fatalf("expected field set to contain 'x'")
	}
	if _, ok := set["y"]; !ok {
		t.Fatalf("expected field set to contain 'y'")
	}
}

func TestFetchFieldFuncTypes(t *testing.T) {
	ft := &ast.FuncType{Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "int"}}}}}
	fields := &ast.FieldList{List: []*ast.Field{
		{Names: []*ast.Ident{{Name: "f"}}, Type: ft},
		{Names: []*ast.Ident{{Name: "n"}}, Type: &ast.Ident{Name: "int"}},
	}}
	d := &codegen.DatagenParsed{Fields: fields}
	m := fetchFieldFuncTypes(d)
	if len(m) != 1 {
		t.Fatalf("expected only one func-typed field, got %d", len(m))
	}
	if _, ok := m["f"]; !ok {
		t.Fatalf("expected map to contain key 'f'")
	}
}

func TestListExtrasGen(t *testing.T) {
	genSet := map[string]struct{}{"a": {}, "c": {}}
	fieldSet := map[string]struct{}{"a": {}}
	extra := listExtrasGen(genSet, fieldSet)
	if len(extra) != 1 || extra[0] != "c" {
		t.Fatalf("expected extras ['c'], got %#v", extra)
	}
}

func TestBuildCallExprs(t *testing.T) {
	calls := []*ast.CallExpr{
		{Fun: &ast.Ident{Name: "foo"}},
		{Fun: &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "Bar"}}},
	}
	d := &codegen.DatagenParsed{Calls: calls}
	m := buildCallExprs(d)
	if len(m) != 1 {
		t.Fatalf("expected only ident-named call to be collected, got %d", len(m))
	}
	if _, ok := m["foo"]; !ok {
		t.Fatalf("expected map to contain key 'foo'")
	}
}

func TestCountParams(t *testing.T) {
	if got := countParams(nil); got != 0 {
		t.Fatalf("expected 0 for nil func type, got %d", got)
	}
	if got := countParams(&ast.FuncType{}); got != 0 {
		t.Fatalf("expected 0 for no params, got %d", got)
	}
	ft1 := &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{{Names: nil, Type: &ast.Ident{Name: "int"}}}}}
	if got := countParams(ft1); got != 1 {
		t.Fatalf("expected 1 for unnamed param, got %d", got)
	}
	ft2 := &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: "a"}, {Name: "b"}}, Type: &ast.Ident{Name: "int"}}}}}
	if got := countParams(ft2); got != 2 {
		t.Fatalf("expected 2 for two named params in same field, got %d", got)
	}
	ft3 := &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{
		{Names: nil, Type: &ast.Ident{Name: "string"}},
		{Names: []*ast.Ident{{Name: "x"}, {Name: "y"}, {Name: "z"}}, Type: &ast.Ident{Name: "int"}},
	}}}
	if got := countParams(ft3); got != 4 {
		t.Fatalf("expected 4 for mixed params, got %d", got)
	}

	_ = token.INT // silence unused import if optimised by toolchain
}


