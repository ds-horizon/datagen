package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/elliotchance/orderedmap/v3"

	"github.com/dream-horizon-org/datagen/utils"
)

func GetDgDirStructure(inputDir, cumulatedPath string) (*utils.DgDir, error) {
	fullPath := filepath.Join(inputDir, cumulatedPath)

	if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
		return GetDgFileStructure(fullPath, cumulatedPath)
	}

	return GetDgDirectoryStructure(inputDir, cumulatedPath)
}

func GetDgFileStructure(filePath, cumulatedPath string) (*utils.DgDir, error) {
	if filepath.Ext(filePath) != ".dg" {
		return nil, fmt.Errorf("file %s is not a .dg file", filePath)
	}
	cleanPath := filepath.Clean(filePath)
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	fileName := strings.TrimSuffix(filepath.Base(filePath), ".dg")
	orderedMap := orderedmap.NewOrderedMap[string, []byte]()
	orderedMap.Set(fileName, content)
	dgDir := &utils.DgDir{
		Name:     cumulatedPath,
		Models:   orderedMap,
		Children: []*utils.DgDir{},
	}
	return dgDir, nil
}

func GetDgDirectoryStructure(inputDir, cumulatedPath string) (*utils.DgDir, error) {
	dgDir := &utils.DgDir{
		Name:     cumulatedPath,
		Models:   orderedmap.NewOrderedMap[string, []byte](),
		Children: []*utils.DgDir{},
	}

	entries, err := os.ReadDir(filepath.Join(inputDir, cumulatedPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", filepath.Join(inputDir, cumulatedPath), err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if strings.IndexFunc(name, unicode.IsSpace) != -1 {
			return nil, fmt.Errorf("directory/file name contains whitespace: %q", name)
		}

		if entry.IsDir() {
			dgDirNested, err := GetDgDirectoryStructure(inputDir, filepath.Join(cumulatedPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			dgDir.Children = append(dgDir.Children, dgDirNested)
			continue
		}
		if filepath.Ext(entry.Name()) == ".dg" {
			cleanPath := filepath.Join(inputDir, cumulatedPath, entry.Name())
			content, readErr := os.ReadFile(filepath.Clean(cleanPath))
			if readErr != nil {
				return nil, readErr
			}

			// converting the file system path to dg delimited string for the map
			// eg: spends/foobar/Alert.dg ==> spends___DG_DIR_DELIMITER___foobar___DG_DIR_DELIMITER___Alert
			fsPath := strings.ReplaceAll(cumulatedPath, string(filepath.Separator), utils.DgDirDelimeter)
			ent := strings.Join([]string{fsPath, entry.Name()}, utils.DgDirDelimeter)
			ent = strings.TrimSuffix(ent, ".dg")
			if cumulatedPath == "" {
				ent = strings.TrimSuffix(entry.Name(), ".dg")
			}
			dgDir.Models.Set(ent, content)
		}
	}
	return dgDir, nil
}
