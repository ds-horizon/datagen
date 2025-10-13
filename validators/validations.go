package validators

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/dream-sports-labs/datagen/codegen"
	"github.com/dream-sports-labs/datagen/utils"
)

type ValidatorFunc func(d *codegen.DatagenParsed, errs *MultiErr)

func Validate(d *codegen.DatagenParsed) error {
	var aggregateErrs MultiErr

	validators := []ValidatorFunc{
		RequiredSectionsValidator,
		NoDuplicateFieldNamesValidator,
		NoMissingGensValidator,
		NoExtraGensValidator,
		GenFnsReturnValidator,
		CallExprsValidator,
		FilePathModelNameValidator,
	}
	for _, validator := range validators {
		validator(d, &aggregateErrs)
	}
	return errorOrNil(&aggregateErrs)
}

func RequiredSectionsValidator(d *codegen.DatagenParsed, errs *MultiErr) {
	if d.Fields == nil || d.Fields.List == nil {
		errs.Add("model has no fields section")
	}
	if d.GenFuns == nil {
		errs.Add("model has no gens section")
	}
}

func NoDuplicateFieldNamesValidator(d *codegen.DatagenParsed, errs *MultiErr) {
	if d == nil || d.Fields == nil || d.Fields.List == nil {
		return
	}
	duplicates, _ := fetchDuplicateAndFieldSet(d)
	if len(duplicates) > 0 {
		errs.Addf("model has duplicate field names: %s", strings.Join(duplicates, ", "))
	}
}

func NoMissingGensValidator(d *codegen.DatagenParsed, errs *MultiErr) {
	if d == nil || d.Fields == nil || d.Fields.List == nil || d.GenFuns == nil {
		return
	}
	genSet := buildGenSet(d)
	missingGens := fetchMissingGens(d, genSet)
	if len(missingGens) > 0 {
		errs.Addf("model has missing gen functions: %s", strings.Join(missingGens, ", "))
	}
}

func NoExtraGensValidator(d *codegen.DatagenParsed, errs *MultiErr) {
	if d == nil || d.Fields == nil || d.Fields.List == nil || d.GenFuns == nil {
		return
	}
	genSet := buildGenSet(d)
	_, fieldSet := fetchDuplicateAndFieldSet(d)
	extrasGens := listExtrasGen(genSet, fieldSet)
	if len(extrasGens) > 0 {
		errs.Addf("found extra gen functions: %s", strings.Join(extrasGens, ", "))
	}
}

func FilePathModelNameValidator(d *codegen.DatagenParsed, errs *MultiErr) {
	splitPath := strings.Split(d.Filepath, utils.DgDirDelimeter)
	if len(splitPath) < 1 {
		errs.Addf("model should be in file named %s.dg, found in %s", d.ModelName, d.Filepath)
	}
	fileName := filepath.Base(splitPath[len(splitPath)-1])
	if fileName != d.ModelName {
		errs.Addf("model should be in file named %s.dg, found in %s.dg", d.ModelName, fileName)
	}
}

func CallExprsValidator(d *codegen.DatagenParsed, errs *MultiErr) {
	if d == nil || d.Fields == nil || d.Fields.List == nil {
		return
	}
	fieldFuncTypes := fetchFieldFuncTypes(d)
	callExprs := buildCallExprs(d)

	for fname, ftype := range fieldFuncTypes {
		if ftype.Results == nil || ftype.Results.List == nil || len(ftype.Results.List) == 0 {
			errs.Addf("field %s must declare a return type", fname)
		}
		expected := countParams(ftype)
		// Skip validation for zero-parameter function fields
		if expected == 0 {
			continue
		}
		call, ok := callExprs[fname]
		if !ok {
			errs.Addf("missing call for field %s", fname)
			continue
		}
		actual := len(call.Args)
		if expected != actual {
			errs.Addf("field %s expects %d args, got %d", fname, expected, actual)
		}
	}
	// Ensure no extra calls
	for cname := range callExprs {
		if _, ok := fieldFuncTypes[cname]; !ok {
			errs.Addf("unknown call %s", cname)
		}
	}
}

// GenFnsReturnValidator ensures each generator function has at least one return statement.
func GenFnsReturnValidator(d *codegen.DatagenParsed, errs *MultiErr) {
	for _, g := range d.GenFuns {
		if g == nil || g.Body == nil {
			errs.Addf("gen func %s must return a value", safeGenName(g))
			continue
		}
		found := false
		ast.Inspect(g.Body, func(n ast.Node) bool {
			if _, ok := n.(*ast.ReturnStmt); ok {
				found = true
				return false
			}
			return true
		})
		if !found {
			errs.Addf("gen func %s must return a value", g.Name)
		}
	}
}

func safeGenName(g *codegen.GenFn) string {
	if g == nil {
		return "<unknown>"
	}
	if strings.TrimSpace(g.Name) == "" {
		return "<unnamed>"
	}
	return g.Name
}
