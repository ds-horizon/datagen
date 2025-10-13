package utils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CompareDirectories(t *testing.T, expectedDir, actualDir string) {
	err := filepath.Walk(expectedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(expectedDir, path)
		if err != nil {
			t.Errorf("Failed to get relative path: %v", err)
			return err
		}

		actualPath := filepath.Join(actualDir, relPath)

		if _, err := os.Stat(actualPath); errors.Is(err, os.ErrNotExist) {
			t.Errorf("Expected file %s not found in generated output", relPath)
			return nil
		}

		expectedContent, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read expected file %s: %v", relPath, err)
			return err
		}

		actualContent, err := os.ReadFile(actualPath)
		if err != nil {
			t.Errorf("Failed to read actual file %s: %v", relPath, err)
			return err
		}

		assert.Equal(t, string(expectedContent), string(actualContent), "Content mismatch in file %s", relPath)

		return nil
	})

	if err != nil {
		t.Errorf("Error walking expected directory: %v", err)
	}

	err = filepath.Walk(actualDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(actualDir, path)
		if err != nil {
			t.Errorf("Failed to get relative path: %v", err)
			return err
		}

		expectedPath := filepath.Join(expectedDir, relPath)
		if _, err := os.Stat(expectedPath); errors.Is(err, os.ErrNotExist) {
			t.Errorf("Unexpected file %s found in generated output", relPath)
		}

		return nil
	})

	if err != nil {
		t.Errorf("Error walking actual directory: %v", err)
	}
}
