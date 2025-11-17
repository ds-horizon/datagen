package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ds-horizon/datagen/utils"
)

func TestGetDgFileStructure(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (string, string)
		cleanupFunc   func(t *testing.T, dir string)
		expectedError bool
		validate      func(t *testing.T, result *utils.DgDir)
	}{
		{
			name: "valid .dg file",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "test.dg")
				content := []byte("model Test { field() string }")
				err := os.WriteFile(filePath, content, 0600)
				require.NoError(t, err)
				return filePath, "test.dg"
			},
			cleanupFunc:   func(t *testing.T, dir string) {},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, "test.dg", result.Name)
				assert.Equal(t, 1, len(result.Models))
				assert.Contains(t, result.Models, "test")
				assert.Equal(t, []byte("model Test { field() string }"), result.Models["test"])
				assert.Equal(t, 0, len(result.Children))
			},
		},
		{
			name: "invalid file extension",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "test.txt")
				content := []byte("not a dg file")
				err := os.WriteFile(filePath, content, 0600)
				require.NoError(t, err)
				return filePath, "test.txt"
			},
			cleanupFunc:   func(t *testing.T, dir string) {},
			expectedError: true,
			validate:      func(t *testing.T, result *utils.DgDir) {},
		},
		{
			name: "non-existent file",
			setupFunc: func(t *testing.T) (string, string) {
				return "/non/existent/file.dg", "file.dg"
			},
			cleanupFunc:   func(t *testing.T, dir string) {},
			expectedError: true,
			validate:      func(t *testing.T, result *utils.DgDir) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cumulatedPath := tt.setupFunc(t)
			defer tt.cleanupFunc(t, filePath)

			result, err := GetDgFileStructure(filePath, cumulatedPath)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validate(t, result)
			}
		})
	}
}

func TestGetDgDirectoryStructure(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (string, string)
		expectedError bool
		validate      func(t *testing.T, result *utils.DgDir)
	}{
		{
			name: "directory with single .dg file",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "model.dg")
				content := []byte("model Example { id() int }")
				err := os.WriteFile(filePath, content, 0600)
				require.NoError(t, err)
				return tmpDir, ""
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, "", result.Name)
				assert.Equal(t, 1, len(result.Models))
				assert.Contains(t, result.Models, "model")
				assert.Equal(t, 0, len(result.Children))
			},
		},
		{
			name: "directory with multiple .dg files",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				file1 := filepath.Join(tmpDir, "model1.dg")
				file2 := filepath.Join(tmpDir, "model2.dg")
				err := os.WriteFile(file1, []byte("model One {}"), 0600)
				require.NoError(t, err)
				err = os.WriteFile(file2, []byte("model Two {}"), 0600)
				require.NoError(t, err)
				return tmpDir, ""
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, 2, len(result.Models))
				assert.Contains(t, result.Models, "model1")
				assert.Contains(t, result.Models, "model2")
			},
		},
		{
			name: "nested directories with .dg files",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				subDir := filepath.Join(tmpDir, "subdir")
				err := os.Mkdir(subDir, 0750)
				require.NoError(t, err)

				file1 := filepath.Join(tmpDir, "root.dg")
				file2 := filepath.Join(subDir, "nested.dg")
				err = os.WriteFile(file1, []byte("model Root {}"), 0600)
				require.NoError(t, err)
				err = os.WriteFile(file2, []byte("model Nested {}"), 0600)
				require.NoError(t, err)
				return tmpDir, ""
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, 1, len(result.Models))
				assert.Contains(t, result.Models, "root")
				assert.Equal(t, 1, len(result.Children))
				assert.Equal(t, "subdir", result.Children[0].Name)
				assert.Equal(t, 1, len(result.Children[0].Models))
			},
		},
		{
			name: "directory with hidden files (should ignore)",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				visibleFile := filepath.Join(tmpDir, "visible.dg")
				hiddenFile := filepath.Join(tmpDir, ".hidden.dg")
				err := os.WriteFile(visibleFile, []byte("model Visible {}"), 0600)
				require.NoError(t, err)
				err = os.WriteFile(hiddenFile, []byte("model Hidden {}"), 0600)
				require.NoError(t, err)
				return tmpDir, ""
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, 1, len(result.Models))
				assert.Contains(t, result.Models, "visible")
				assert.NotContains(t, result.Models, "hidden")
			},
		},
		{
			name: "directory with non-.dg files (should ignore)",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				dgFile := filepath.Join(tmpDir, "model.dg")
				txtFile := filepath.Join(tmpDir, "readme.txt")
				err := os.WriteFile(dgFile, []byte("model Test {}"), 0600)
				require.NoError(t, err)
				err = os.WriteFile(txtFile, []byte("README"), 0600)
				require.NoError(t, err)
				return tmpDir, ""
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, 1, len(result.Models))
				assert.Contains(t, result.Models, "model")
			},
		},
		{
			name: "empty directory",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				return tmpDir, ""
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, 0, len(result.Models))
				assert.Equal(t, 0, len(result.Children))
			},
		},
		{
			name: "directory name with whitespace (should error)",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				subDir := filepath.Join(tmpDir, "bad name")
				err := os.Mkdir(subDir, 0750)
				require.NoError(t, err)
				return tmpDir, ""
			},
			expectedError: true,
			validate:      func(t *testing.T, result *utils.DgDir) {},
		},
		{
			name: "file name with whitespace (should error)",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				badFile := filepath.Join(tmpDir, "bad file.dg")
				err := os.WriteFile(badFile, []byte("model Test {}"), 0600)
				require.NoError(t, err)
				return tmpDir, ""
			},
			expectedError: true,
			validate:      func(t *testing.T, result *utils.DgDir) {},
		},
		{
			name: "non-existent directory",
			setupFunc: func(t *testing.T) (string, string) {
				return "/non/existent/directory", ""
			},
			expectedError: true,
			validate:      func(t *testing.T, result *utils.DgDir) {},
		},
		{
			name: "cumulated path in nested directory",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				subDir := filepath.Join(tmpDir, "sub1", "sub2")
				err := os.MkdirAll(subDir, 0750)
				require.NoError(t, err)

				file := filepath.Join(subDir, "nested.dg")
				err = os.WriteFile(file, []byte("model Nested {}"), 0600)
				require.NoError(t, err)
				return tmpDir, filepath.Join("sub1", "sub2")
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, filepath.Join("sub1", "sub2"), result.Name)
				assert.Equal(t, 1, len(result.Models))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputDir, cumulatedPath := tt.setupFunc(t)

			result, err := GetDgDirectoryStructure(inputDir, cumulatedPath)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validate(t, result)
			}
		})
	}
}

func TestGetDgDirStructure(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) (string, string)
		expectedError bool
		validate      func(t *testing.T, result *utils.DgDir)
	}{
		{
			name: "input is a single file",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "single.dg")
				content := []byte("model Single { id() int }")
				err := os.WriteFile(filePath, content, 0600)
				require.NoError(t, err)
				return filePath, ""
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, 1, len(result.Models))
				assert.Contains(t, result.Models, "single")
			},
		},
		{
			name: "input is a directory",
			setupFunc: func(t *testing.T) (string, string) {
				tmpDir := t.TempDir()
				file := filepath.Join(tmpDir, "model.dg")
				err := os.WriteFile(file, []byte("model Test {}"), 0600)
				require.NoError(t, err)
				return tmpDir, ""
			},
			expectedError: false,
			validate: func(t *testing.T, result *utils.DgDir) {
				assert.NotNil(t, result)
				assert.Equal(t, 1, len(result.Models))
			},
		},
		{
			name: "input is non-existent",
			setupFunc: func(t *testing.T) (string, string) {
				return "/non/existent/path", ""
			},
			expectedError: true,
			validate:      func(t *testing.T, result *utils.DgDir) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputDir, cumulatedPath := tt.setupFunc(t)

			result, err := GetDgDirStructure(inputDir, cumulatedPath)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validate(t, result)
			}
		})
	}
}

func TestDgDirStructureModelPaths(t *testing.T) {
	t.Run("model path delimiter conversion", func(t *testing.T) {
		tmpDir := t.TempDir()
		subDir := filepath.Join(tmpDir, "models", "users")
		err := os.MkdirAll(subDir, 0750)
		require.NoError(t, err)

		file := filepath.Join(subDir, "User.dg")
		err = os.WriteFile(file, []byte("model User {}"), 0600)
		require.NoError(t, err)

		result, err := GetDgDirectoryStructure(tmpDir, "")
		require.NoError(t, err)

		// Check that the file system path is properly converted to use DgDirDelimeter
		found := false
		expectedKey := "models" + utils.DgDirDelimeter + "users" + utils.DgDirDelimeter + "User"
		for key := range result.Children[0].Children[0].Models {
			if key == expectedKey {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find model with path using DgDirDelimeter")
	})
}

func TestDgDirModelCount(t *testing.T) {
	t.Run("count models in nested structure", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create root level file
		rootFile := filepath.Join(tmpDir, "root.dg")
		err := os.WriteFile(rootFile, []byte("model Root {}"), 0600)
		require.NoError(t, err)

		// Create nested directory with files
		subDir1 := filepath.Join(tmpDir, "sub1")
		subDir2 := filepath.Join(tmpDir, "sub2")
		err = os.Mkdir(subDir1, 0755)
		require.NoError(t, err)
		err = os.Mkdir(subDir2, 0755)
		require.NoError(t, err)

		file1 := filepath.Join(subDir1, "model1.dg")
		file2 := filepath.Join(subDir2, "model2.dg")
		file3 := filepath.Join(subDir2, "model3.dg")
		err = os.WriteFile(file1, []byte("model M1 {}"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(file2, []byte("model M2 {}"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(file3, []byte("model M3 {}"), 0600)
		require.NoError(t, err)

		result, err := GetDgDirectoryStructure(tmpDir, "")
		require.NoError(t, err)

		// Should count 1 root + 1 in sub1 + 2 in sub2 = 4 total
		assert.Equal(t, 4, result.ModelCount())
	})
}

func TestGetDgDirectoryStructureDeepNesting(t *testing.T) {
	t.Run("deeply nested directory structure", func(t *testing.T) {
		tmpDir := t.TempDir()

		deepPath := filepath.Join(tmpDir, "level1", "level2", "level3", "level4")
		err := os.MkdirAll(deepPath, 0750)
		require.NoError(t, err)

		for i, level := range []string{
			tmpDir,
			filepath.Join(tmpDir, "level1"),
			filepath.Join(tmpDir, "level1", "level2"),
			filepath.Join(tmpDir, "level1", "level2", "level3"),
			deepPath,
		} {
			file := filepath.Join(level, "model.dg")
			content := []byte("model Level {}")
			err = os.WriteFile(file, content, 0600)
			require.NoError(t, err)
			_ = i
		}

		result, err := GetDgDirectoryStructure(tmpDir, "")
		require.NoError(t, err)
		assert.Equal(t, 5, result.ModelCount())
	})
}

func TestGetDgDirectoryStructureComplexPaths(t *testing.T) {
	t.Run("complex nested paths with multiple models per directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create structure: api/v1/, api/v2/, internal/models/
		apiV1 := filepath.Join(tmpDir, "api", "v1")
		apiV2 := filepath.Join(tmpDir, "api", "v2")
		internal := filepath.Join(tmpDir, "internal", "models")

		err := os.MkdirAll(apiV1, 0750)
		require.NoError(t, err)
		err = os.MkdirAll(apiV2, 0750)
		require.NoError(t, err)
		err = os.MkdirAll(internal, 0750)
		require.NoError(t, err)

		// Add multiple models in api/v1
		err = os.WriteFile(filepath.Join(apiV1, "User.dg"), []byte("model User {}"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(apiV1, "Order.dg"), []byte("model Order {}"), 0600)
		require.NoError(t, err)

		// Add models in api/v2
		err = os.WriteFile(filepath.Join(apiV2, "Product.dg"), []byte("model Product {}"), 0600)
		require.NoError(t, err)

		// Add models in internal/models
		err = os.WriteFile(filepath.Join(internal, "Config.dg"), []byte("model Config {}"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(internal, "Settings.dg"), []byte("model Settings {}"), 0600)
		require.NoError(t, err)

		result, err := GetDgDirectoryStructure(tmpDir, "")
		require.NoError(t, err)
		assert.Equal(t, 5, result.ModelCount())
		assert.Equal(t, 2, len(result.Children)) // api and internal
	})

	t.Run("single file in deeply nested cumulated path", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create nested structure
		nestedPath := filepath.Join(tmpDir, "a", "b", "c")
		err := os.MkdirAll(nestedPath, 0750)
		require.NoError(t, err)

		// Add a file in the nested path
		err = os.WriteFile(filepath.Join(nestedPath, "Deep.dg"), []byte("model Deep {}"), 0600)
		require.NoError(t, err)

		// Test with cumulated path
		result, err := GetDgDirectoryStructure(tmpDir, filepath.Join("a", "b", "c"))
		require.NoError(t, err)
		assert.Equal(t, 1, result.ModelCount())
		assert.Equal(t, filepath.Join("a", "b", "c"), result.Name)
	})
}

func TestGetDgFileStructureEdgeCases(t *testing.T) {
	t.Run("file with special characters in name", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Test with file containing underscore and numbers
		file := filepath.Join(tmpDir, "Model_v2_test.dg")
		err := os.WriteFile(file, []byte("model Test {}"), 0600)
		require.NoError(t, err)

		result, err := GetDgFileStructure(file, "Model_v2_test.dg")
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result.Models, "Model_v2_test")
	})

	t.Run("empty .dg file", func(t *testing.T) {
		tmpDir := t.TempDir()

		file := filepath.Join(tmpDir, "Empty.dg")
		err := os.WriteFile(file, []byte(""), 0600)
		require.NoError(t, err)

		result, err := GetDgFileStructure(file, "Empty.dg")
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []byte(""), result.Models["Empty"])
	})

	t.Run("large .dg file", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create a large model with many fields
		largeModel := "model LargeModel {\n  fields {\n"
		for i := 0; i < 100; i++ {
			largeModel += "    field" + string(rune(i)) + "() string\n"
		}
		largeModel += "  }\n  gens {\n"
		for i := 0; i < 100; i++ {
			largeModel += "    func field" + string(rune(i)) + "() { return \"value\" }\n"
		}
		largeModel += "  }\n}"

		file := filepath.Join(tmpDir, "LargeModel.dg")
		err := os.WriteFile(file, []byte(largeModel), 0600)
		require.NoError(t, err)

		result, err := GetDgFileStructure(file, "LargeModel.dg")
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, len(result.Models["LargeModel"]) > 1000)
	})
}

func TestGetDgDirectoryStructureWithMixedContent(t *testing.T) {
	t.Run("directory with .dg files and subdirectories", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create root level .dg file
		err := os.WriteFile(filepath.Join(tmpDir, "Root.dg"), []byte("model Root {}"), 0600)
		require.NoError(t, err)

		// Create subdirectory with files
		subDir := filepath.Join(tmpDir, "subdir")
		err = os.Mkdir(subDir, 0755)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(subDir, "Sub.dg"), []byte("model Sub {}"), 0600)
		require.NoError(t, err)

		// Create nested subdirectory
		nestedDir := filepath.Join(subDir, "nested")
		err = os.Mkdir(nestedDir, 0755)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(nestedDir, "Nested.dg"), []byte("model Nested {}"), 0600)
		require.NoError(t, err)

		result, err := GetDgDirectoryStructure(tmpDir, "")
		require.NoError(t, err)
		assert.Equal(t, 3, result.ModelCount())
		assert.Equal(t, 1, len(result.Models))               // Root level has 1 model
		assert.Equal(t, 1, len(result.Children))             // One subdirectory
		assert.Equal(t, 1, len(result.Children[0].Children)) // Nested subdirectory
	})

	t.Run("hidden directories should be ignored", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create visible file
		err := os.WriteFile(filepath.Join(tmpDir, "Visible.dg"), []byte("model Visible {}"), 0600)
		require.NoError(t, err)

		// Create hidden directory with files
		hiddenDir := filepath.Join(tmpDir, ".hidden")
		err = os.Mkdir(hiddenDir, 0755)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(hiddenDir, "Hidden.dg"), []byte("model Hidden {}"), 0600)
		require.NoError(t, err)

		result, err := GetDgDirectoryStructure(tmpDir, "")
		require.NoError(t, err)
		assert.Equal(t, 1, result.ModelCount())  // Should only count visible file
		assert.Equal(t, 0, len(result.Children)) // Hidden directory should be ignored
	})
}

func TestGetDgDirStructureWithBothFileAndDirectory(t *testing.T) {
	t.Run("input that is both file and directory parent", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create a single .dg file
		file := filepath.Join(tmpDir, "Single.dg")
		err := os.WriteFile(file, []byte("model Single {}"), 0600)
		require.NoError(t, err)

		// Test with file path
		resultFile, err := GetDgDirStructure(file, "")
		require.NoError(t, err)
		assert.Equal(t, 1, resultFile.ModelCount())

		// Test with directory path
		resultDir, err := GetDgDirStructure(tmpDir, "")
		require.NoError(t, err)
		assert.Equal(t, 1, resultDir.ModelCount())
	})
}
