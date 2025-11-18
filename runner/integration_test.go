package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ds-horizon/datagen/codegen"
)

func TestIntegrationTranspileValidModels(t *testing.T) {
	validModelsDir := filepath.Join("testdata", "valid")
	goldenFilesDir := filepath.Join("testdata", "transpiledTestFiles")
	outputDir := t.TempDir()

	dgDirData, err := GetDgDirStructure(validModelsDir, "")
	require.NoError(t, err, "failed to read valid models directory")
	require.Greater(t, dgDirData.ModelCount(), 0, "no models found in valid directory")

	parsed, err := processDgDirData(dgDirData, outputDir, []*codegen.DatagenParsed{})
	require.NoError(t, err, "failed to transpile models")
	require.Greater(t, len(parsed), 0, "no parsed models returned")

	err = codegen.Codegen(parsed, outputDir, dgDirData)
	require.NoError(t, err, "code generation failed")

	generatedFiles, err := filepath.Glob(filepath.Join(outputDir, "*.go"))
	require.NoError(t, err, "failed to list generated files")
	require.Greater(t, len(generatedFiles), 0, "no Go files generated")

	for _, generatedFile := range generatedFiles {
		fileName := filepath.Base(generatedFile)
		goldenFile := filepath.Join(goldenFilesDir, fileName)

		t.Run(fileName, func(t *testing.T) {
			_, err := os.Stat(goldenFile)
			if os.IsNotExist(err) {
				t.Logf("Warning: golden file does not exist: %s", goldenFile)
				t.Skip("golden file not found")
				return
			}

			generatedContent, err := os.ReadFile(generatedFile) // #nosec G304 -- Test file path constructed from known test directory
			require.NoError(t, err, "failed to read generated file: %s", generatedFile)

			goldenContent, err := os.ReadFile(goldenFile) // #nosec G304 -- Golden file path constructed from known test directory
			require.NoError(t, err, "failed to read golden file: %s", goldenFile)

			assert.Equal(t, goldenContent, generatedContent,
				"Generated file %s does not match golden file %s", fileName, goldenFile)
		})
	}
}

func TestIntegrationTranspileSpecificModel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tests := []struct {
		name          string
		modelFile     string
		expectedFiles []string
	}{
		{
			name:      "simple model",
			modelFile: "simple.dg",
			expectedFiles: []string{
				"simple.go",
				"simple_mysql.go",
				"simple_init_mysql.go",
				"simple_sink_mysql.go",
			},
		},
		{
			name:      "minimal model",
			modelFile: "minimal.dg",
			expectedFiles: []string{
				"minimal.go",
				"minimal_mysql.go",
				"minimal_init_mysql.go",
				"minimal_sink_mysql.go",
			},
		},
		{
			name:      "nested model",
			modelFile: "nested.dg",
			expectedFiles: []string{
				"nested.go",
				"nested_mysql.go",
				"nested_init_mysql.go",
				"nested_sink_mysql.go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelPath := filepath.Join("testdata", "valid", tt.modelFile)
			goldenFilesDir := filepath.Join("testdata", "transpiledTestFiles")
			outputDir := t.TempDir()

			if _, err := os.Stat(modelPath); os.IsNotExist(err) {
				t.Skip("model file not found")
				return
			}

			err := findAndTranspileDatagenModels(outputDir, modelPath)
			require.NoError(t, err, "failed to transpile model")

			for _, expectedFile := range tt.expectedFiles {
				generatedFile := filepath.Join(outputDir, expectedFile)
				goldenFile := filepath.Join(goldenFilesDir, expectedFile)

				t.Run(expectedFile, func(t *testing.T) {
					_, err := os.Stat(generatedFile)
					if os.IsNotExist(err) {
						t.Logf("Generated file not found: %s", generatedFile)
						return
					}
					require.NoError(t, err, "error checking generated file")

					_, err = os.Stat(goldenFile)
					if os.IsNotExist(err) {
						t.Skip("golden file not found")
						return
					}

					generatedContent, err := os.ReadFile(generatedFile) // #nosec G304 -- Test file path constructed from known test directory
					require.NoError(t, err, "failed to read generated file")

					goldenContent, err := os.ReadFile(goldenFile) // #nosec G304 -- Golden file path constructed from known test directory
					require.NoError(t, err, "failed to read golden file")

					assert.Equal(t, goldenContent, generatedContent,
						"Generated file does not match golden file")
				})
			}
		})
	}
}

func TestIntegrationInvalidModels(t *testing.T) {
	tests := []struct {
		name          string
		modelFile     string
		expectedError string
	}{
		{
			name:          "invalid syntax",
			modelFile:     "invalid_syntax.dg",
			expectedError: "failed to parse",
		},
		{
			name:          "empty model",
			modelFile:     "empty_model.dg",
			expectedError: "model has no fields section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelPath := filepath.Join("testdata", "invalid", tt.modelFile)
			outputDir := t.TempDir()

			err := findAndTranspileDatagenModels(outputDir, modelPath)
			assert.Error(t, err, "expected transpilation to fail for invalid model")
			if err != nil {
				assert.Contains(t, err.Error(), tt.expectedError,
					"error message should contain expected text")
			}
		})
	}
}

func TestIntegrationUpdateGoldenFiles(t *testing.T) {
	updateGolden := false
	for _, arg := range os.Args {
		if arg == "-update-golden" {
			updateGolden = true
			break
		}
	}

	if !updateGolden {
		t.Skip("skipping golden file update (use -update-golden flag to update)")
	}

	validModelsDir := filepath.Join("testdata", "valid")
	goldenFilesDir := filepath.Join("testdata", "transpiledTestFiles")
	outputDir := t.TempDir()

	dgDirData, err := GetDgDirStructure(validModelsDir, "")
	require.NoError(t, err)

	parsed, err := processDgDirData(dgDirData, outputDir, []*codegen.DatagenParsed{})
	require.NoError(t, err)

	err = codegen.Codegen(parsed, outputDir, dgDirData)
	require.NoError(t, err)

	generatedFiles, err := filepath.Glob(filepath.Join(outputDir, "*.go"))
	require.NoError(t, err)

	err = os.MkdirAll(goldenFilesDir, 0o750)
	require.NoError(t, err)

	for _, generatedFile := range generatedFiles {
		fileName := filepath.Base(generatedFile)
		goldenFile := filepath.Join(goldenFilesDir, fileName)

		content, err := os.ReadFile(generatedFile) // #nosec G304 -- Test file path constructed from known test directory
		require.NoError(t, err)

		err = os.WriteFile(goldenFile, content, 0o600)
		require.NoError(t, err)

		t.Logf("Updated golden file: %s", fileName)
	}

	t.Logf("Successfully updated %d golden files", len(generatedFiles))
}
