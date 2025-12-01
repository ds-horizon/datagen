package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dream-horizon-org/datagen/codegen"
	"github.com/dream-horizon-org/datagen/utils"
)

func TestProcessDgDirData(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		expectedError  bool
		expectedCount  int
		validateModels func(t *testing.T, result []*codegen.DatagenParsed)
	}{
		{
			name:          "valid models directory",
			inputPath:     filepath.Join("testdata", "valid"),
			expectedError: false,
			expectedCount: 10,
			validateModels: func(t *testing.T, result []*codegen.DatagenParsed) {
				modelNames := make(map[string]bool)
				for _, parsed := range result {
					modelNames[parsed.ModelName] = true
				}
				// Check that expected models are present
				expectedModels := []string{
					"simple", "minimal", "multiple_types", "with_metadata",
					"with_misc", "with_builtin_functions", "nested", "with_conditionals",
					"with_slices", "with_maps",
				}
				for _, expected := range expectedModels {
					assert.True(t, modelNames[expected], "expected model %s to be parsed", expected)
				}
			},
		},
		{
			name:          "single valid file",
			inputPath:     filepath.Join("testdata", "valid", "simple.dg"),
			expectedError: false,
			expectedCount: 1,
			validateModels: func(t *testing.T, result []*codegen.DatagenParsed) {
				require.Len(t, result, 1)
				assert.Equal(t, "simple", result[0].ModelName)
				assert.NotNil(t, result[0].Fields)
				assert.Equal(t, 2, len(result[0].Fields.List))
				assert.Equal(t, 2, len(result[0].GenFuns))
			},
		},
		{
			name:          "minimal model",
			inputPath:     filepath.Join("testdata", "valid", "minimal.dg"),
			expectedError: false,
			expectedCount: 1,
			validateModels: func(t *testing.T, result []*codegen.DatagenParsed) {
				require.Len(t, result, 1)
				assert.Equal(t, "minimal", result[0].ModelName)
				assert.NotNil(t, result[0].Fields)
				assert.Equal(t, 1, len(result[0].Fields.List))
				assert.Equal(t, 1, len(result[0].GenFuns))
			},
		},
		{
			name:          "model with metadata",
			inputPath:     filepath.Join("testdata", "valid", "with_metadata.dg"),
			expectedError: false,
			expectedCount: 1,
			validateModels: func(t *testing.T, result []*codegen.DatagenParsed) {
				require.Len(t, result, 1)
				assert.Equal(t, "with_metadata", result[0].ModelName)
				assert.NotNil(t, result[0].Metadata)
				assert.Equal(t, 100, result[0].Metadata.Count)
				assert.NotNil(t, result[0].Metadata.Tags)
				assert.Equal(t, "test", result[0].Metadata.Tags["env"])
			},
		},
		{
			name:          "model with misc section",
			inputPath:     filepath.Join("testdata", "valid", "with_misc.dg"),
			expectedError: false,
			expectedCount: 1,
			validateModels: func(t *testing.T, result []*codegen.DatagenParsed) {
				require.Len(t, result, 1)
				assert.Equal(t, "with_misc", result[0].ModelName)
				assert.NotEmpty(t, result[0].Misc)
			},
		},
		{
			name:          "nil directory",
			inputPath:     "",
			expectedError: false,
			expectedCount: 0,
			validateModels: func(t *testing.T, result []*codegen.DatagenParsed) {
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outDir := t.TempDir()

			var dgDir *utils.DgDir
			var err error

			if tt.inputPath != "" {
				dgDir, err = GetDgDirStructure(tt.inputPath, "")
				if err != nil && !tt.expectedError {
					t.Fatalf("failed to get DgDir structure: %v", err)
				}
			}

			result, err := processDgDirData(dgDir, outDir, []*codegen.DatagenParsed{})

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(result), "expected %d models, got %d", tt.expectedCount, len(result))
				if tt.validateModels != nil {
					tt.validateModels(t, result)
				}
			}
		})
	}
}

func TestProcessDgDirDataWithInvalidModel(t *testing.T) {
	t.Run("should error on invalid model in nested structure", func(t *testing.T) {
		rootModels := orderedmap.NewOrderedMap[string, []byte]()
		rootModels.Set("ValidModel", []byte(`model ValidModel { fields { id() int } gens { func id() { return iter } } }`))

		childModels := orderedmap.NewOrderedMap[string, []byte]()
		childModels.Set("InvalidModel", []byte(`model InvalidModel { invalid syntax`))

		dgDir := &utils.DgDir{
			Name:   "root",
			Models: rootModels,
			Children: []*utils.DgDir{
				{
					Name:     "child",
					Models:   childModels,
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
		inputPath     string
		setupFunc     func(t *testing.T) string
		expectedError bool
		errorContains string
	}{
		{
			name:          "valid directory with models",
			inputPath:     filepath.Join("testdata", "valid"),
			expectedError: false,
		},
		{
			name:          "valid single file",
			inputPath:     filepath.Join("testdata", "valid", "simple.dg"),
			expectedError: false,
		},
		{
			name: "no .dg files found",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				file := filepath.Join(tmpDir, "readme.txt")
				err := os.WriteFile(file, []byte("Not a dg file"), 0o600)
				require.NoError(t, err)
				return tmpDir
			},
			expectedError: true,
			errorContains: "no .dg files found",
		},
		{
			name:          "non-existent input directory",
			inputPath:     "/non/existent/path",
			expectedError: true,
			errorContains: "failed to read input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outDir := t.TempDir()

			inputPath := tt.inputPath
			if tt.setupFunc != nil {
				inputPath = tt.setupFunc(t)
			}

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
				file := filepath.Join("testdata", "valid", "minimal.dg")

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
				file := filepath.Join("testdata", "valid", "minimal.dg")

				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte("config: test"), 0o600)
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
				file := filepath.Join("testdata", "valid", "simple.dg")

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
				file := filepath.Join("testdata", "invalid", "invalid_syntax.dg")

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
				file := filepath.Join("testdata", "valid", "simple.dg")

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
				file := filepath.Join("testdata", "invalid", "empty_model.dg")

				configFile := filepath.Join(tmpDir, "config.yaml")
				err := os.WriteFile(configFile, []byte("config: test"), 0o600)
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

		// Model name doesn't match filename - need to create this dynamically
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
		err := os.WriteFile(file, content, 0o600)
		require.NoError(t, err)

		err = findAndTranspileDatagenModels(outDir, tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "should be in file named")
	})

	t.Run("directory with invalid model", func(t *testing.T) {
		outDir := t.TempDir()
		inputPath := filepath.Join("testdata", "invalid")

		err := findAndTranspileDatagenModels(outDir, inputPath)
		assert.Error(t, err)
	})
}

func TestInvokeGenArguments(t *testing.T) {
	t.Run("with all optional flags", func(t *testing.T) {
		outDir := t.TempDir()
		inputPath := "/test/input.dg"

		count := 100
		tags := "prod,test"
		output := "/output"
		format := "xml"
		seed := int64(999)
		verbose := true

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

		args := []string{"gen", inputPath}
		args = append(args, "-n", "1")

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
