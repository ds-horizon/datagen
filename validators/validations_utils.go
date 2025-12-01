package validators

import (
	"go/ast"
	"strings"

	"github.com/dream-horizon-org/datagen/codegen"
)

// buildGenSet creates a set of declared generator function names.
func buildGenSet(d *codegen.DatagenParsed) map[string]struct{} {
	genSet := map[string]struct{}{}
	for _, g := range d.GenFuns {
		if g == nil {
			continue
		}
		name := strings.TrimSpace(g.Name)
		if name == "" {
			continue
		}
		genSet[name] = struct{}{}
	}
	return genSet
}

func fetchMissingGens(d *codegen.DatagenParsed, genSet map[string]struct{}) []string {
	if d == nil || d.Fields == nil || d.Fields.List == nil {
		return nil
	}
	var missingGens []string
	for _, f := range d.Fields.List {
		for _, n := range f.Names {
			if n == nil {
				continue
			}
			field := strings.TrimSpace(n.Name)
			if field == "" {
				continue
			}
			if _, ok := genSet[field]; !ok {
				missingGens = append(missingGens, field)
			}
		}
	}
	return missingGens
}

func fetchDuplicateAndFieldSet(d *codegen.DatagenParsed) ([]string, map[string]struct{}) {
	fieldSet := map[string]struct{}{}
	if d == nil || d.Fields == nil || d.Fields.List == nil {
		return nil, fieldSet
	}
	seen := map[string]struct{}{}
	var duplicates []string

	for _, f := range d.Fields.List {
		for _, n := range f.Names {
			if n == nil {
				continue
			}
			field := strings.TrimSpace(n.Name)
			if field == "" {
				continue
			}
			fieldSet[field] = struct{}{}
			// check duplicate fields
			if _, exists := seen[field]; exists {
				duplicates = append(duplicates, field)
			} else {
				seen[field] = struct{}{}
			}
		}
	}
	return duplicates, fieldSet
}

func fetchFieldFuncTypes(d *codegen.DatagenParsed) map[string]*ast.FuncType {
	fieldFuncTypes := map[string]*ast.FuncType{}
	if d == nil || d.Fields == nil || d.Fields.List == nil {
		return fieldFuncTypes
	}
	for _, f := range d.Fields.List {
		for _, n := range f.Names {
			if n == nil {
				continue
			}
			field := strings.TrimSpace(n.Name)
			if field == "" {
				continue
			}
			if ft, ok := f.Type.(*ast.FuncType); ok {
				fieldFuncTypes[field] = ft
			}
		}
	}
	return fieldFuncTypes
}

// listExtrasGen returns gens declared without a corresponding field.
func listExtrasGen(genSet, fieldSet map[string]struct{}) []string {
	var extras []string
	for gen := range genSet {
		if _, ok := fieldSet[gen]; !ok {
			extras = append(extras, gen)
		}
	}
	return extras
}

// buildCallExprs builds call name -> *ast.CallExpr
func buildCallExprs(d *codegen.DatagenParsed) map[string]*ast.CallExpr {
	callExprs := map[string]*ast.CallExpr{}
	if d == nil || d.Calls == nil {
		return callExprs
	}
	for _, c := range d.Calls {
		if ident, ok := c.Fun.(*ast.Ident); ok {
			callExprs[ident.Name] = c
		}
	}
	return callExprs
}

func countParams(ft *ast.FuncType) int {
	if ft == nil || ft.Params == nil || ft.Params.List == nil {
		return 0
	}
	count := 0
	for _, p := range ft.Params.List {
		if len(p.Names) == 0 {
			count += 1
		} else {
			count += len(p.Names)
		}
	}
	return count
}
