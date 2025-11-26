package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func RemoveDirIfExists(dirPath string) error {
	if _, statErr := os.Stat(dirPath); statErr == nil {
		slog.Debug(fmt.Sprintf("removing existing generated directory %s", dirPath))
		if rmErr := os.RemoveAll(dirPath); rmErr != nil {
			return fmt.Errorf("failed to remove existing generated directory\n  path: %s\n  cause: %w", dirPath, rmErr)
		}
	}
	return nil
}
