package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/ds-horizon/datagen/utils"
)

const (
	tmplPackage           = "templates/package.tmpl"
	tmplMain              = "templates/main.tmpl"
	tmplCommands          = "templates/commands.tmpl"
	tmplModelManager      = "templates/model_manager.tmpl"
	tmplSinkManager       = "templates/sink_manager.tmpl"
	tmplSinkMysqlModel    = "templates/sink_mysql_model.tmpl"
	tmplTags              = "templates/tags.go.tmpl"
	tmplConfig            = "templates/config.go.tmpl"
	tmplLinks             = "templates/links.go.tmpl"
	tmplMySQLConfig       = "templates/mysql_config.tmpl"
	tmplPostgresConfig    = "templates/postgres_config.tmpl"
	tmplKafkaConfig       = "templates/kafka_config.tmpl"
	tmplWriters           = "templates/writers.tmpl"
	tmplGoMod             = "templates/go.mod.tmpl"
	tmplGoSum             = "templates/go.sum.tmpl"
	tmplStdlib            = "templates/stdlib.go.tmpl"
	tmplLogger            = "templates/logger.go.tmpl"
	tmplMetadata          = "templates/metadata_struct.tmpl"
	tmplBaseStruct        = "templates/base_struct.tmpl"
	tmplGenerator         = "templates/generator_struct.tmpl"
	tmplDataHolder        = "templates/data_holder_struct.tmpl"
	tmplWrapperFunc       = "templates/wrapper_func.tmpl"
	tmplGenFunction       = "templates/gen_function.tmpl"
	tmplInitFunction      = "templates/init_function.tmpl"
	tmplCSV               = "templates/csv_function.tmpl"
	tmplJSON              = "templates/json_function.tmpl"
	tmplXML               = "templates/xml_function.tmpl"
	tmplMysqlSink         = "templates/load_mysql.tmpl"
	tmplMysqlInit         = "templates/init_mysql.tmpl"
	tmplPostgresSink      = "templates/load_postgres.tmpl"
	tmplPostgresInit      = "templates/init_postgres.tmpl"
	tmplSinkPostgresModel = "templates/sink_postgres_model.tmpl"
)

type SectionGenerator func(d *DatagenParsed) (string, error)

type modelNameData struct {
	DgDir                    *utils.DgDir
	SanitisedModelNames      []string
	FullyQualifiedModelNames []string
}

func Codegen(parsed []*DatagenParsed, dirPath string, dgDir *utils.DgDir) error {
	if len(parsed) == 0 {
		return nil
	}

	for _, result := range parsed {
		if err := codegenModel(result, dirPath); err != nil {
			return fmt.Errorf("failed to generate code for model\n  model: %s\n  cause: %w", result.ModelName, err)
		}
	}

	if err := codegenCommons(parsed, dirPath, dgDir); err != nil {
		return fmt.Errorf("error generating main.go: %v", err)
	}

	return nil
}

func codegenModel(parsed *DatagenParsed, dirPath string) error {
	modelDir := dirPath
	if err := os.MkdirAll(modelDir, 0o750); err != nil {
		return err
	}

	generators := map[string]SectionGenerator{
		"misc":             generateMiscSection,
		"metadata":         generateMetadataSection,
		"base_struct":      generateBaseStruct,
		"generator_struct": generateGeneratorStruct,
		"data_holder":      generateDataHolderStruct,
		"generator_funcs":  generateGeneratorFuncs,
		"gen_function":     generateGenFunction,
		"init_function":    generateInitFunction,
		"csv_functions":    generateCSVFunctions,
		"json_functions":   generateJSONFunctions,
		"xml_functions":    generateXMLFunctions,
	}

	sections := make(map[string]string, len(generators))
	for name, generator := range generators {
		section, err := generator(parsed)
		if err != nil {
			return err
		}
		sections[name] = section
	}

	tmpl, err := template.ParseFS(templates, tmplPackage)
	if err != nil {
		return fmt.Errorf("failed to parse package template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, sections); err != nil {
		return fmt.Errorf("failed to generate code: %v", err)
	}

	modelPath := filepath.Join(modelDir, fmt.Sprintf("%s.go", parsed.FullyQualifiedModelName))
	if err := writeFormattedGoFile(modelPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write generated model file\n  path: %s\n  cause: %w", modelPath, err)
	}

	if err := parsed.generateMySQLInitFile(modelDir); err != nil {
		return fmt.Errorf("failed to generate MySQL init file\n  model: %s\n  cause: %w", parsed.FullyQualifiedModelName, err)
	}
	if err := parsed.generateMySQLLoadFile(modelDir); err != nil {
		return fmt.Errorf("failed to generate MySQL load file\n  model: %s\n  cause: %w", parsed.FullyQualifiedModelName, err)
	}
	if err := parsed.generateMySQLSinkFile(modelDir); err != nil {
		return fmt.Errorf("failed to generate MySQL sink file\n  model: %s\n  cause: %w", parsed.FullyQualifiedModelName, err)
	}

	if err := parsed.generatePostgresInitFile(modelDir); err != nil {
		return fmt.Errorf("failed to generate Postgres init file\n  model: %s\n  cause: %w", parsed.FullyQualifiedModelName, err)
	}
	if err := parsed.generatePostgresLoadFile(modelDir); err != nil {
		return fmt.Errorf("failed to generate Postgres load file\n  model: %s\n  cause: %w", parsed.FullyQualifiedModelName, err)
	}
	if err := parsed.generatePostgresSinkFile(modelDir); err != nil {
		return fmt.Errorf("failed to generate Postgres sink file\n  model: %s\n  cause: %w", parsed.FullyQualifiedModelName, err)
	}

	return nil
}

func codegenCommons(parsed []*DatagenParsed, dirPath string, dgDir *utils.DgDir) error {
	modelNames := make([]string, 0, len(parsed))
	sanitisedModelNames := make([]string, 0, len(parsed))
	for _, p := range parsed {
		modelNames = append(modelNames, p.FullyQualifiedModelName)
		sanitisedModelNames = append(sanitisedModelNames, strings.ReplaceAll(p.FullyQualifiedModelName, utils.DgDirDelimeter, "."))
	}

	if err := generateMainFile(dirPath); err != nil {
		return fmt.Errorf("failed to generate main.go: %v", err)
	}

	if err := generateCommandsFile(dirPath, modelNames); err != nil {
		return fmt.Errorf("failed to generate commands.go: %v", err)
	}

	if err := generateModelManagerFile(dirPath, modelNames, dgDir); err != nil {
		return fmt.Errorf("failed to generate model_manager.go: %v", err)
	}

	if err := generateSinkManagerFile(dirPath, &modelNameData{SanitisedModelNames: sanitisedModelNames, FullyQualifiedModelNames: modelNames}); err != nil {
		return fmt.Errorf("failed to generate sink_manager.go: %v", err)
	}

	// Generate tags.go
	tagsTmpl, err := template.ParseFS(templates, tmplTags)
	if err != nil {
		return fmt.Errorf("failed to parse tags template: %v", err)
	}

	var buf bytes.Buffer
	if err := tagsTmpl.Execute(&buf, sanitisedModelNames); err != nil {
		return fmt.Errorf("failed to generate tags.go: %v", err)
	}
	tagsPath := filepath.Join(dirPath, "tags.go")
	if err := writeFormattedGoFile(tagsPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write tags.go: %v", err)
	}

	// Generate config.go
	cfgTmpl, err := template.ParseFS(templates, tmplConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config template: %v", err)
	}
	buf.Reset()
	if err := cfgTmpl.Execute(&buf, modelNames); err != nil {
		return fmt.Errorf("failed to generate config.go: %v", err)
	}
	cfgPath := filepath.Join(dirPath, "config.go")
	if err := writeFormattedGoFile(cfgPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write config.go: %v", err)
	}

	staticFiles := map[string]string{
		tmplWriters:        "writers.go",
		tmplGoMod:          "go.mod",
		tmplGoSum:          "go.sum",
		tmplStdlib:         "stdlib.go",
		tmplLogger:         "logger.go",
		tmplMySQLConfig:    "mysql_config.go",
		tmplPostgresConfig: "postgres_config.go",
		tmplKafkaConfig:    "kafka_config.go",
		tmplLinks:          "links.go",
	}
	if err := copyStaticTemplates(dirPath, staticFiles); err != nil {
		return fmt.Errorf("failed to copy static templates\n  output_dir: %s\n  cause: %w", dirPath, err)
	}

	return nil
}

func generateMetadataSection(d *DatagenParsed) (string, error) {
	s, err := renderFS(tmplMetadata, metadataVars(d))
	if err != nil {
		return "", fmt.Errorf("error generating metadata function: %w", err)
	}
	return s, nil
}

func generateMiscSection(d *DatagenParsed) (string, error) {
	var buf bytes.Buffer
	_, err := fmt.Fprintf(&buf, "%s\n", d.Misc)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func generateBaseStruct(d *DatagenParsed) (string, error) {
	s, err := renderFS(tmplBaseStruct, fieldsVars(d))
	if err != nil {
		return "", fmt.Errorf("error generating base struct: %w", err)
	}
	return s, nil
}

func generateGeneratorStruct(d *DatagenParsed) (string, error) {
	s, err := renderFS(tmplGenerator, fieldsVars(d))
	if err != nil {
		return "", fmt.Errorf("error generating generator struct: %w", err)
	}
	return s, nil
}

func generateDataHolderStruct(d *DatagenParsed) (string, error) {
	s, err := renderFS(tmplDataHolder, fieldsVars(d))
	if err != nil {
		return "", fmt.Errorf("failed to generate data holder section\n  model: %s\n  cause: %w", d.FullyQualifiedModelName, err)
	}
	return s, nil
}

func generateGeneratorFuncs(d *DatagenParsed) (string, error) {
	tmpl, err := template.ParseFS(templates, tmplWrapperFunc)
	if err != nil {
		return "", fmt.Errorf("failed to parse template\n  template: %s\n  cause: %w", tmplWrapperFunc, err)
	}

	var buf bytes.Buffer
	fset := token.NewFileSet()

	for _, genFn := range d.GenFuns {
		fieldType, fieldName := findFieldTypeAndName(d, genFn.Name)
		if fieldName == "" {
			continue
		}

		// Convert function parameters to string
		var paramsBuf bytes.Buffer
		var varsBuf bytes.Buffer
		if genFn.Calls != nil {
			if funcType, ok := genFn.Calls.Fun.(*ast.FuncType); ok && funcType.Params != nil {
				for i, field := range funcType.Params.List {
					if i > 0 {
						paramsBuf.WriteString(", ")
						varsBuf.WriteString(", ")
					}
					for j, name := range field.Names {
						if j > 0 {
							paramsBuf.WriteString(", ")
							varsBuf.WriteString(", ")
						}
						paramsBuf.WriteString(name.Name)
						varsBuf.WriteString(name.Name)
					}
					paramsBuf.WriteString(" ")
					varsBuf.WriteString(" ")
					err := printer.Fprint(&paramsBuf, fset, field.Type)
					if err != nil {
						return "", err
					}
				}
			}
		}

		var bodyBuf bytes.Buffer
		if genFn.Body != nil {
			err := printer.Fprint(&bodyBuf, fset, genFn.Body)
			if err != nil {
				return "", err
			}
		}

		data := wrapperFuncData{
			ModelName:               d.ModelName,
			FullyQualifiedModelName: d.FullyQualifiedModelName,
			FieldName:               fieldName,
			FieldType:               fieldType,
			GenFuncParams:           paramsBuf.String(),
			GenFuncVars:             varsBuf.String(),
			GenFuncBody:             bodyBuf.String(),
		}

		if err := tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("failed to generate wrapper function\n  model: %s\n  field: %s\n  cause: %w", d.FullyQualifiedModelName, fieldName, err)
		}
	}

	return buf.String(), nil
}

func generateGenFunction(d *DatagenParsed) (string, error) {
	s, err := renderFS(tmplGenFunction, fieldsVars(d))
	if err != nil {
		return "", fmt.Errorf("failed to generate Gen section\n  model: %s\n  cause: %w", d.FullyQualifiedModelName, err)
	}
	return s, nil
}

func generateInitFunction(d *DatagenParsed) (string, error) {
	s, err := renderFS(tmplInitFunction, fieldsVars(d))
	if err != nil {
		return "", fmt.Errorf("failed to generate init section\n  model: %s\n  cause: %w", d.FullyQualifiedModelName, err)
	}
	return s, nil
}

func generateCSVFunctions(d *DatagenParsed) (string, error) {
	s, err := renderFS(tmplCSV, fieldsVars(d))
	if err != nil {
		return "", fmt.Errorf("failed to generate CSV functions section\n  model: %s\n  cause: %w", d.FullyQualifiedModelName, err)
	}
	return s, nil
}

func generateJSONFunctions(d *DatagenParsed) (string, error) {
	s, err := renderFS(tmplJSON, fieldsVars(d))
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON functions section\n  model: %s\n  cause: %w", d.FullyQualifiedModelName, err)
	}
	return s, nil
}

func generateXMLFunctions(d *DatagenParsed) (string, error) {
	funcs := template.FuncMap{
		"XmlPrefix": func(s string) string {
			if s == "" {
				return s
			}
			return "Xml_" + s
		},
	}
	s, err := renderFSWithFuncs(tmplXML, funcs, "xml_function.tmpl", fieldsVars(d))
	if err != nil {
		return "", fmt.Errorf("failed to generate XML functions section\n  model: %s\n  cause: %w", d.FullyQualifiedModelName, err)
	}
	return s, nil
}

// generateMySQLLoadFile renders templates/load_mysql.tmpl into <ModelName>_mysql.go
func (d *DatagenParsed) generateMySQLLoadFile(modelDir string) error {
	if len(getFieldData(d)) == 0 {
		return nil
	}

	ib, err := renderFS(tmplMysqlSink, fieldsVars(d))
	if err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplMysqlSink, err)
	}

	outPath := filepath.Join(modelDir, fmt.Sprintf("%s_mysql.go", d.FullyQualifiedModelName))
	if err := writeFormattedGoFile(outPath, []byte(ib)); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", outPath, err)
	}
	return nil
}

// generateMySQLInitFile renders templates/init_mysql.tmpl into <ModelName>_init_mysql.go
func (d *DatagenParsed) generateMySQLInitFile(modelDir string) error {
	ib, err := renderFS(tmplMysqlInit, fieldsVars(d))
	if err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplMysqlInit, err)
	}
	initMySQLPath := filepath.Join(modelDir, fmt.Sprintf("%s_init_mysql.go", d.FullyQualifiedModelName))

	if err := writeFormattedGoFile(initMySQLPath, []byte(ib)); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", initMySQLPath, err)
	}
	return nil
}

// generateMySQLSinkFile renders templates/sink_mysql_model.tmpl into <ModelName>_sink_mysql.go
func (d *DatagenParsed) generateMySQLSinkFile(modelDir string) error {
	ib, err := renderFS(tmplSinkMysqlModel, fieldsVars(d))
	if err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplSinkMysqlModel, err)
	}
	sinkMySQLPath := filepath.Join(modelDir, fmt.Sprintf("%s_sink_mysql.go", d.FullyQualifiedModelName))

	if err := writeFormattedGoFile(sinkMySQLPath, []byte(ib)); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", sinkMySQLPath, err)
	}
	return nil
}

// generatePostgresLoadFile renders templates/load_postgres.tmpl into <ModelName>_postgres.go
func (d *DatagenParsed) generatePostgresLoadFile(modelDir string) error {
	if len(getFieldData(d)) == 0 {
		return nil
	}

	ib, err := renderFS(tmplPostgresSink, fieldsVars(d))
	if err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplPostgresSink, err)
	}

	outPath := filepath.Join(modelDir, fmt.Sprintf("%s_postgres.go", d.FullyQualifiedModelName))
	if err := writeFormattedGoFile(outPath, []byte(ib)); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", outPath, err)
	}
	return nil
}

// generatePostgresInitFile renders templates/init_postgres.tmpl into <ModelName>_init_postgres.go
func (d *DatagenParsed) generatePostgresInitFile(modelDir string) error {
	ib, err := renderFS(tmplPostgresInit, fieldsVars(d))
	if err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplPostgresInit, err)
	}
	initPostgresPath := filepath.Join(modelDir, fmt.Sprintf("%s_init_postgres.go", d.FullyQualifiedModelName))

	if err := writeFormattedGoFile(initPostgresPath, []byte(ib)); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", initPostgresPath, err)
	}
	return nil
}

// generatePostgresSinkFile renders templates/sink_postgres_model.tmpl into <ModelName>_sink_postgres.go
func (d *DatagenParsed) generatePostgresSinkFile(modelDir string) error {
	ib, err := renderFS(tmplSinkPostgresModel, fieldsVars(d))
	if err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplSinkPostgresModel, err)
	}
	sinkPostgresPath := filepath.Join(modelDir, fmt.Sprintf("%s_sink_postgres.go", d.FullyQualifiedModelName))

	if err := writeFormattedGoFile(sinkPostgresPath, []byte(ib)); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", sinkPostgresPath, err)
	}
	return nil
}

// generateMainFile generates the main.go file (CLI entry point)
func generateMainFile(dirPath string) error {
	content, err := templates.ReadFile(tmplMain)
	if err != nil {
		return fmt.Errorf("failed to read template\n  template: %s\n  cause: %w", tmplMain, err)
	}

	mainPath := filepath.Join(dirPath, "main.go")
	if err := writeFormattedGoFile(mainPath, content); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", mainPath, err)
	}
	return nil
}

// generateCommandsFile generates the commands.go file
func generateCommandsFile(dirPath string, modelNames []string) error {
	tmpl, err := template.ParseFS(templates, tmplCommands)
	if err != nil {
		return fmt.Errorf("failed to parse template\n  template: %s\n  cause: %w", tmplCommands, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, modelNames); err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplCommands, err)
	}

	commandsPath := filepath.Join(dirPath, "commands.go")
	if err := writeFormattedGoFile(commandsPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", commandsPath, err)
	}
	return nil
}

// generateModelManagerFile generates the model_manager.go file
func generateModelManagerFile(dirPath string, modelNames []string, dgDir *utils.DgDir) error {
	funcs := template.FuncMap{
		"last": func(full string) string {
			parts := strings.Split(full, utils.DgDirDelimeter)
			return parts[len(parts)-1]
		},
		"dirName": func(d *utils.DgDir) string {
			name := d.Name
			if idx := strings.LastIndex(name, "/"); idx >= 0 {
				name = name[idx+1:]
			}
			return name
		},
		"models": func(d *utils.DgDir) []string {
			keys := make([]string, 0, d.Models.Len())
			for k := range d.Models.AllFromFront() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			return keys
		},
		"children": func(d *utils.DgDir) []*utils.DgDir {
			out := make([]*utils.DgDir, len(d.Children))
			copy(out, d.Children)
			sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
			return out
		},
		"dot": func(full string) string {
			return strings.ReplaceAll(full, utils.DgDirDelimeter, ".")
		},
	}

	tmpl, err := template.New(filepath.Base(tmplModelManager)).Funcs(funcs).ParseFS(templates, tmplModelManager)
	if err != nil {
		return fmt.Errorf("failed to parse template\n  template: %s\n  cause: %w", tmplModelManager, err)
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, &modelNameData{DgDir: dgDir, SanitisedModelNames: modelNames}); err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplModelManager, err)
	}

	modelManagerPath := filepath.Join(dirPath, "model_manager.go")
	if err := writeFormattedGoFile(modelManagerPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", modelManagerPath, err)
	}
	return nil
}

// generateSinkManagerFile generates the sink_manager.go file
func generateSinkManagerFile(dirPath string, modelNameData *modelNameData) error {
	funcs := template.FuncMap{
		"dot": func(full string) string {
			return strings.ReplaceAll(full, utils.DgDirDelimeter, ".")
		},
	}
	tmpl, err := template.New(filepath.Base(tmplSinkManager)).Funcs(funcs).ParseFS(templates, tmplSinkManager)
	if err != nil {
		return fmt.Errorf("failed to parse template\n  template: %s\n  cause: %w", tmplSinkManager, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, modelNameData); err != nil {
		return fmt.Errorf("failed to render template\n  template: %s\n  cause: %w", tmplSinkManager, err)
	}

	sinkManagerPath := filepath.Join(dirPath, "sink_manager.go")
	if err := writeFormattedGoFile(sinkManagerPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write generated file\n  path: %s\n  cause: %w", sinkManagerPath, err)
	}

	return nil
}
