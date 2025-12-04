package codegen

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dream-horizon-org/datagen/utils"
)

//go:embed templates
var templates embed.FS

type fieldData struct {
	Name     string
	Type     string
	InitArgs string
}

type templateVars struct {
	ModelName               string
	FullyQualifiedModelName string
	Fields                  []fieldData
	Metadata                Metadata
}

type wrapperFuncData struct {
	ModelName               string
	FullyQualifiedModelName string
	FieldName               string
	FieldType               string
	GenFuncParams           string
	GenFuncVars             string
	GenFuncBody             string
}

type serialiserFuncData struct {
	ModelName               string
	FullyQualifiedModelName string
	SerialiserFuncBody      string
}

type GenFn struct {
	Name  string
	Calls *ast.CallExpr
	Body  *ast.BlockStmt
}

type SerialiserFunc struct {
	Body *ast.BlockStmt
}

type DatagenParsed struct {
	ModelName               string
	FullyQualifiedModelName string
	Fields                  *ast.FieldList
	Misc                    string
	GenFuns                 []*GenFn
	SerialiserFunc          *SerialiserFunc
	Calls                   []*ast.CallExpr
	Metadata                *Metadata
	Filepath                string
}

type Metadata struct {
	Count int
	Tags  map[string]string
}

func getMetadata(d *DatagenParsed) Metadata {
	if d.Metadata == nil {
		return Metadata{
			Count: utils.DefaultMetadataCount,
			Tags:  map[string]string{},
		}
	}

	metadata := *(d.Metadata)

	// if count is improperly set, use the default metadata count
	if metadata.Count <= 0 {
		metadata.Count = utils.DefaultMetadataCount
	}
	return metadata
}

func getFieldData(d *DatagenParsed) []fieldData {
	var fields []fieldData
	if d.Fields != nil {
		for _, field := range d.Fields.List {
			for _, name := range field.Names {
				initArgs := ""
				for _, call := range d.Calls {
					if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == name.Name {
						var argsBuf bytes.Buffer
						for i, arg := range call.Args {
							if i > 0 {
								argsBuf.WriteString(", ")
							}
							err := printer.Fprint(&argsBuf, token.NewFileSet(), arg)
							if err != nil {
								return nil
							}
						}
						initArgs = argsBuf.String()
						break
					}
				}

				fields = append(fields, fieldData{
					Name:     name.Name,
					Type:     getTypeString(field.Type),
					InitArgs: initArgs,
				})
			}
		}
	}
	return fields
}

func getTypeString(expr ast.Expr) string {
	if ft, ok := expr.(*ast.FuncType); ok {
		if ft.Results != nil && len(ft.Results.List) > 0 {
			return stringifyExpr(ft.Results.List[0].Type)
		}
		return "func"
	}
	return stringifyExpr(expr)
}

func stringifyExpr(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name

	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", stringifyExpr(t.X), t.Sel.Name)

	case *ast.StarExpr:
		return "*" + stringifyExpr(t.X)

	case *ast.ArrayType:
		if t.Len != nil {
			return fmt.Sprintf("[%s]%s", stringifyExpr(t.Len), stringifyExpr(t.Elt))
		}
		return "[]" + stringifyExpr(t.Elt)

	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", stringifyExpr(t.Key), stringifyExpr(t.Value))

	case *ast.ParenExpr:
		return "(" + stringifyExpr(t.X) + ")"

	default:
		return fmt.Sprintf("%T", e)
	}
}

func findFieldTypeAndName(d *DatagenParsed, fieldName string) (string, string) {
	if d.Fields == nil {
		return "", ""
	}

	for _, field := range d.Fields.List {
		for _, name := range field.Names {
			if name.Name == fieldName {
				return getTypeString(field.Type), name.Name
			}
		}
	}
	return "", ""
}

func fieldsVars(d *DatagenParsed) templateVars {
	return templateVars{ModelName: d.ModelName, Fields: getFieldData(d), FullyQualifiedModelName: d.FullyQualifiedModelName}
}

func metadataVars(d *DatagenParsed) templateVars {
	return templateVars{ModelName: d.ModelName, Metadata: getMetadata(d), FullyQualifiedModelName: d.FullyQualifiedModelName}
}

func copyStaticTemplates(dirPath string, files map[string]string) error {
	for src, dst := range files {
		content, err := templates.ReadFile(src)
		if err != nil {
			return fmt.Errorf("failed to read template\n  template: %s\n  cause: %w", src, err)
		}
		outPath := filepath.Join(dirPath, dst)
		if err := os.WriteFile(outPath, content, 0o600); err != nil {
			return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", outPath, err)
		}
	}
	return nil
}

// shared helpers to reduce duplication when rendering templates
func renderFS(tmplPath string, data templateVars) (string, error) {
	tmpl, err := template.ParseFS(templates, tmplPath)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderFSWithFuncs(tmplPath string, funcs template.FuncMap, execName string, data templateVars) (string, error) {
	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(funcs).ParseFS(templates, tmplPath)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	name := execName
	if name == "" {
		name = filepath.Base(tmplPath)
	}
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// writeFormattedGoFile formats Go source code and writes it to the given path.
func writeFormattedGoFile(path string, src []byte) error {
	formatted, err := format.Source(src)
	if err != nil {
		return fmt.Errorf("failed to format generated file\n  file: %s\n  cause: %w", filepath.Base(path), err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("failed to create directory\n  path: %s\n  cause: %w", path, err)
	}

	if err := os.WriteFile(path, formatted, 0o600); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", path, err)
	}
	return nil
}
