package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ds-horizon/datagen/codegen"
	"github.com/ds-horizon/datagen/utils"
)

func TestTranspile(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) *utils.DgDir
		expectedError bool
		validate      func(t *testing.T, result []*codegen.DatagenParsed, err error)
	}{
		{
			name: "valid single model",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name: "test",
					Models: map[string][]byte{
						"TestModel": []byte(`model TestModel {
	fields {
		id() int
	}
	gens {
		func id() {
			return iter
		}
	}
}`),
					},
					Children: []*utils.DgDir{},
				}
			},
			expectedError: false,
			validate: func(t *testing.T, result []*codegen.DatagenParsed, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(result))
				assert.Equal(t, "TestModel", result[0].ModelName)
			},
		},
		{
			name: "multiple models",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name: "test",
					Models: map[string][]byte{
						"Model1": []byte(`model Model1 {
	fields {
		field1() string
	}
	gens {
		func field1() {
			return "test"
		}
	}
}`),
						"Model2": []byte(`model Model2 {
	fields {
		field2() int
	}
	gens {
		func field2() {
			return iter
		}
	}
}`),
					},
					Children: []*utils.DgDir{},
				}
			},
			expectedError: false,
			validate: func(t *testing.T, result []*codegen.DatagenParsed, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 2, len(result))
			},
		},
		{
			name: "invalid syntax",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name: "test",
					Models: map[string][]byte{
						"BadModel": []byte(`model BadModel { invalid syntax here`),
					},
					Children: []*utils.DgDir{},
				}
			},
			expectedError: true,
			validate: func(t *testing.T, result []*codegen.DatagenParsed, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "empty model (should error)",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name: "test",
					Models: map[string][]byte{
						"EmptyModel": []byte(`model EmptyModel {}`),
					},
					Children: []*utils.DgDir{},
				}
			},
			expectedError: true,
			validate: func(t *testing.T, result []*codegen.DatagenParsed, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "model has no fields section")
			},
		},
		{
			name: "model with metadata",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name: "test",
					Models: map[string][]byte{
						"MetadataModel": []byte(`model MetadataModel {
	metadata {
		count: 5
	}
	fields {
		id() int
	}
	gens {
		func id() {
			return iter
		}
	}
}`),
					},
					Children: []*utils.DgDir{},
				}
			},
			expectedError: false,
			validate: func(t *testing.T, result []*codegen.DatagenParsed, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(result))
				assert.Equal(t, "MetadataModel", result[0].ModelName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dgDir := tt.setupFunc(t)
			result, err := transpile(dgDir)
			tt.validate(t, result, err)
		})
	}
}

func TestProcessDgDirData(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) *utils.DgDir
		expectedError bool
		expectedCount int
	}{
		{
			name: "single level with one model",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name: "root",
					Models: map[string][]byte{
						"Model1": []byte(`model Model1 { fields { id() int } gens { func id() { return iter } } }`),
					},
					Children: []*utils.DgDir{},
				}
			},
			expectedError: false,
			expectedCount: 1,
		},
		{
			name: "nested structure with multiple models",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name: "root",
					Models: map[string][]byte{
						"RootModel": []byte(`model RootModel { fields { id() int } gens { func id() { return iter } } }`),
					},
					Children: []*utils.DgDir{
						{
							Name: "child1",
							Models: map[string][]byte{
								"ChildModel1": []byte(`model ChildModel1 { fields { name() string } gens { func name() { return "test" } } }`),
							},
							Children: []*utils.DgDir{},
						},
						{
							Name: "child2",
							Models: map[string][]byte{
								"ChildModel2": []byte(`model ChildModel2 { fields { age() int } gens { func age() { return iter } } }`),
							},
							Children: []*utils.DgDir{},
						},
					},
				}
			},
			expectedError: false,
			expectedCount: 3,
		},
		{
			name: "deeply nested structure",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name: "root",
					Models: map[string][]byte{
						"L0": []byte(`model L0 { fields { id() int } gens { func id() { return iter } } }`),
					},
					Children: []*utils.DgDir{
						{
							Name: "level1",
							Models: map[string][]byte{
								"L1": []byte(`model L1 { fields { id() int } gens { func id() { return iter } } }`),
							},
							Children: []*utils.DgDir{
								{
									Name: "level2",
									Models: map[string][]byte{
										"L2": []byte(`model L2 { fields { id() int } gens { func id() { return iter } } }`),
									},
									Children: []*utils.DgDir{
										{
											Name: "level3",
											Models: map[string][]byte{
												"L3": []byte(`model L3 { fields { id() int } gens { func id() { return iter } } }`),
											},
											Children: []*utils.DgDir{},
										},
									},
								},
							},
						},
					},
				}
			},
			expectedError: false,
			expectedCount: 4,
		},
		{
			name: "nil directory",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return nil
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "empty directory structure",
			setupFunc: func(t *testing.T) *utils.DgDir {
				return &utils.DgDir{
					Name:     "empty",
					Models:   map[string][]byte{},
					Children: []*utils.DgDir{},
				}
			},
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dgDir := tt.setupFunc(t)
			outDir := t.TempDir()

			result, err := processDgDirData(dgDir, outDir, []*codegen.DatagenParsed{})

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(result))
			}
		})
	}
}

func TestProcessDgDirDataWithInvalidModel(t *testing.T) {
	t.Run("should error on invalid model in nested structure", func(t *testing.T) {
		dgDir := &utils.DgDir{
			Name: "root",
			Models: map[string][]byte{
				"ValidModel": []byte(`model ValidModel { fields { id() int } gens { func id() { return iter } } }`),
			},
			Children: []*utils.DgDir{
				{
					Name: "child",
					Models: map[string][]byte{
						"InvalidModel": []byte(`model InvalidModel { invalid syntax`),
					},
					Children: []*utils.DgDir{},
				},
			},
		}
		outDir := t.TempDir()

		result, err := processDgDirData(dgDir, outDir, []*codegen.DatagenParsed{})

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestFindAndTranspileDatagenModels(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (string, string)
		expectedError bool
		errorContains string
	}{
		{
			name: "valid directory with models",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				outDir := t.TempDir()

				file := filepath.Join(tmpDir, "Test.dg")
				content := []byte(`model Test {
	fields {
		id() int
	}
	gens {
		func id() {
			return iter
		}
	}
}`)
				err := os.WriteFile(file, content, 0600)
				require.NoError(t, err)

				return tmpDir, outDir
			},
			expectedError: false,
		},
		{
			name: "no .dg files found",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				outDir := t.TempDir()

				// Create some non-.dg files
				file := filepath.Join(tmpDir, "readme.txt")
				err := os.WriteFile(file, []byte("Not a dg file"), 0600)
				require.NoError(t, err)

				return tmpDir, outDir
			},
			expectedError: true,
			errorContains: "no .dg files found",
		},
		{
			name: "non-existent input directory",
			setupFunc: func(t *testing.T) (string, string) {
				outDir := t.TempDir()
				return "/non/existent/path", outDir
			},
			expectedError: true,
			errorContains: "failed to read input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputPath, outDir := tt.setupFunc(t)

			err := findAndTranspileDatagenModels(outDir, inputPath)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildAndRunGenFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (*cobra.Command, []string)
		expectedError bool
		errorContains string
	}{
		{
			name: "valid flags with noexec",
			setupFunc: func(t *testing.T) (*cobra.Command, []string) {
				tmpDir := t.TempDir()
				file := filepath.Join(tmpDir, "Test.dg")
				err := os.WriteFile(file, []byte(`model Test { fields { id() int } gens { func id() { return iter } } }`), 0600)
				require.NoError(t, err)

				cmd := &cobra.Command{}
				cmd.Flags().Int("count", 10, "")
				cmd.Flags().String("tags", "", "")
				cmd.Flags().String("output", tmpDir, "")
				cmd.Flags().String("format", "json", "")
				cmd.Flags().Int64("seed", 0, "")
				cmd.Flags().Bool("noexec", true, "")
				cmd.Flags().Bool("verbose", false, "")

				return cmd, []string{file}
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := tt.setupFunc(t)

			err := BuildAndRunGen(cmd, args)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildAndRunExecuteFlagParsing(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (*cobra.Command, []string)
		expectedError bool
	}{
		{
			name: "valid flags with noexec",
			setupFunc: func(t *testing.T) (*cobra.Command, []string) {
				tmpDir := t.TempDir()
				file := filepath.Join(tmpDir, "Test.dg")
				err := os.WriteFile(file, []byte(`model Test { fields { id() int } gens { func id() { return iter } } }`), 0600)
				require.NoError(t, err)

				configFile := filepath.Join(tmpDir, "config.yaml")
				err = os.WriteFile(configFile, []byte("config: test"), 0600)
				require.NoError(t, err)

				cmd := &cobra.Command{}
				cmd.Flags().String("config", configFile, "")
				cmd.Flags().String("output", tmpDir, "")
				cmd.Flags().Bool("noexec", true, "")
				cmd.Flags().Bool("verbose", false, "")

				return cmd, []string{file}
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := tt.setupFunc(t)

			if len(args) == 0 {
				return
			}

			err := BuildAndRunExecute(cmd, args)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogCapturedOutput(t *testing.T) {
	// This is a simple test to ensure the function doesn't panic
	t.Run("empty output", func(t *testing.T) {
		var stdout, stderr strings.Builder
		logCapturedOutput(stdout, stderr)
		// Should not panic
	})

	t.Run("output with content", func(t *testing.T) {
		var stdout, stderr strings.Builder
		stdout.WriteString("stdout content")
		stderr.WriteString("stderr content")
		logCapturedOutput(stdout, stderr)
		// Should not panic
	})
}

func TestLogOutputLines(t *testing.T) {
	t.Run("empty output", func(t *testing.T) {
		// Should handle empty strings gracefully
		logOutputLines("", "test")
		logOutputLines("   ", "test")
		logOutputLines("\n\n", "test")
	})

	t.Run("output with content", func(t *testing.T) {
		// Should not panic with valid content
		logOutputLines("test output", "command stdout")
		logOutputLines("error output", "command stderr")
	})
}

func TestBuildAndRunGenErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (*cobra.Command, []string)
		expectedError bool
		errorContains string
	}{
		{
			name: "invalid count flag type",
			setupFunc: func(t *testing.T) (*cobra.Command, []string) {
				tmpDir := t.TempDir()
				file := filepath.Join(tmpDir, "Test.dg")
				err := os.WriteFile(file, []byte(`model Test { fields { id() int } gens { func id() { return iter } } }`), 0600)
				require.NoError(t, err)

				cmd := &cobra.Command{}
				// Intentionally not setting up flags to test error paths
				return cmd, []string{file}
			},
			expectedError: true,
			errorContains: "count",
		},
		{
			name: "file with syntax error",
			setupFunc: func(t *testing.T) (*cobra.Command, []string) {
				tmpDir := t.TempDir()
				file := filepath.Join(tmpDir, "Bad.dg")
				err := os.WriteFile(file, []byte(`model Bad { invalid syntax`), 0600)
				require.NoError(t, err)

				cmd := &cobra.Command{}
				cmd.Flags().Int("count", 10, "")
				cmd.Flags().String("tags", "", "")
				cmd.Flags().String("output", tmpDir, "")
				cmd.Flags().String("format", "json", "")
				cmd.Flags().Int64("seed", 0, "")
				cmd.Flags().Bool("noexec", true, "")
				cmd.Flags().Bool("verbose", false, "")

				return cmd, []string{file}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := tt.setupFunc(t)

			if len(args) == 0 {
				return
			}

			err := BuildAndRunGen(cmd, args)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			}
		})
	}
}

func TestBuildAndRunExecuteErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (*cobra.Command, []string)
		expectedError bool
		errorContains string
	}{
		{
			name: "invalid config flag type",
			setupFunc: func(t *testing.T) (*cobra.Command, []string) {
				tmpDir := t.TempDir()
				file := filepath.Join(tmpDir, "Test.dg")
				err := os.WriteFile(file, []byte(`model Test { fields { id() int } gens { func id() { return iter } } }`), 0600)
				require.NoError(t, err)

				cmd := &cobra.Command{}
				// Intentionally not setting up flags to test error paths
				return cmd, []string{file}
			},
			expectedError: true,
			errorContains: "config",
		},
		{
			name: "file with validation error",
			setupFunc: func(t *testing.T) (*cobra.Command, []string) {
				tmpDir := t.TempDir()
				file := filepath.Join(tmpDir, "Invalid.dg")
				// Empty model that will fail validation
				err := os.WriteFile(file, []byte(`model Invalid {}`), 0600)
				require.NoError(t, err)

				configFile := filepath.Join(tmpDir, "config.yaml")
				err = os.WriteFile(configFile, []byte("config: test"), 0600)
				require.NoError(t, err)

				cmd := &cobra.Command{}
				cmd.Flags().String("config", configFile, "")
				cmd.Flags().String("output", tmpDir, "")
				cmd.Flags().Bool("noexec", true, "")
				cmd.Flags().Bool("verbose", false, "")

				return cmd, []string{file}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := tt.setupFunc(t)

			if len(args) == 0 {
				return
			}

			err := BuildAndRunExecute(cmd, args)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			}
		})
	}
}

func TestFindAndTranspileDatagenModelsEdgeCases(t *testing.T) {
	t.Run("model with mismatched filename", func(t *testing.T) {
		tmpDir := t.TempDir()
		outDir := t.TempDir()

		// Model name doesn't match filename
		file := filepath.Join(tmpDir, "Wrong.dg")
		content := []byte(`model DifferentName {
	fields {
		id() int
	}
	gens {
		func id() {
			return iter
		}
	}
}`)
		err := os.WriteFile(file, content, 0600)
		require.NoError(t, err)

		err = findAndTranspileDatagenModels(outDir, tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "should be in file named")
	})

	t.Run("multiple models with errors", func(t *testing.T) {
		tmpDir := t.TempDir()
		outDir := t.TempDir()

		// Create multiple files, one with error
		file1 := filepath.Join(tmpDir, "Good.dg")
		err := os.WriteFile(file1, []byte(`model Good { fields { id() int } gens { func id() { return iter } } }`), 0600)
		require.NoError(t, err)

		file2 := filepath.Join(tmpDir, "Bad.dg")
		err = os.WriteFile(file2, []byte(`model Bad { invalid`), 0600)
		require.NoError(t, err)

		err = findAndTranspileDatagenModels(outDir, tmpDir)
		assert.Error(t, err)
	})
}

func TestInvokeGenArguments(t *testing.T) {
	t.Run("with all optional flags", func(t *testing.T) {
		outDir := t.TempDir()
		inputPath := "/test/input.dg"

		// Test with all flags including verbose
		count := 100
		tags := "prod,test"
		output := "/output"
		format := "xml"
		seed := int64(999)
		verbose := true

		// Construct expected args
		args := []string{"gen", inputPath}
		args = append(args, "-n", "100")
		if tags != "" {
			args = append(args, "-t", tags)
		}
		if output != "" {
			args = append(args, "-o", output)
		}
		if format != "" {
			args = append(args, "-f", format)
		}
		if seed != 0 {
			args = append(args, "--seed", "999")
		}
		if verbose {
			args = append(args, "-v")
		}

		expectedArgs := []string{"gen", inputPath, "-n", "100", "-t", "prod,test", "-o", "/output", "-f", "xml", "--seed", "999", "-v"}
		assert.Equal(t, expectedArgs, args)

		_ = outDir
		_ = count
	})

	t.Run("without optional flags", func(t *testing.T) {
		inputPath := "/test/input.dg"

		// Test with minimal flags
		args := []string{"gen", inputPath}
		args = append(args, "-n", "1")
		// No tags, output, format, seed, or verbose

		expectedArgs := []string{"gen", inputPath, "-n", "1"}
		assert.Equal(t, expectedArgs, args)
	})
}

func TestInvokeExecuteArguments(t *testing.T) {
	t.Run("with verbose flag", func(t *testing.T) {
		inputPath := "/test/input.dg"
		config := "/test/config.yaml"
		output := "/output"
		verbose := true

		args := []string{"execute", inputPath}
		args = append(args, "-c", config)
		if output != "" {
			args = append(args, "-o", output)
		}
		if verbose {
			args = append(args, "-v")
		}

		expectedArgs := []string{"execute", inputPath, "-c", "/test/config.yaml", "-o", "/output", "-v"}
		assert.Equal(t, expectedArgs, args)
	})

	t.Run("without verbose flag", func(t *testing.T) {
		inputPath := "/test/input.dg"
		config := "/test/config.yaml"

		args := []string{"execute", inputPath}
		args = append(args, "-c", config)

		expectedArgs := []string{"execute", inputPath, "-c", "/test/config.yaml"}
		assert.Equal(t, expectedArgs, args)
	})
}
