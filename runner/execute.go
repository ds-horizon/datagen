package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dream-sports-labs/datagen/codegen"
	"github.com/dream-sports-labs/datagen/parser"
	"github.com/dream-sports-labs/datagen/utils"
	"github.com/dream-sports-labs/datagen/validators"
	"github.com/spf13/cobra"
)

func BuildAndRunGen(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		return fmt.Errorf("invalid value for --count: %w", err)
	}
	tags, err := cmd.Flags().GetString("tags")
	if err != nil {
		return fmt.Errorf("invalid value for --tags: %w", err)
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("invalid value for --output: %w", err)
	}
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return fmt.Errorf("invalid value for --format: %w", err)
	}
	noexec, err := cmd.Flags().GetBool("noexec")
	if err != nil {
		return fmt.Errorf("invalid value for --noexec: %w", err)
	}

	// Target directory where generated code and binary will live
	outDir := filepath.Join(output, "target")
	if err := findAndTranspileDatagenModels(outDir, inputPath); err != nil {
		return err
	}
	if !noexec {
		if err := invokeGen(outDir, count, tags, output, format, inputPath); err != nil {
			return err
		}
	}

	return nil
}

func invokeGen(outDir string, count int, tags string, output string, format string, inputPath string) error {
	binaryPath, _ := buildTranspiledBinary(outDir)
	args := []string{"gen", inputPath}
	args = append(args, "-n", fmt.Sprintf("%d", count))
	if strings.TrimSpace(tags) != "" {
		args = append(args, "-t", tags)
	}
	if strings.TrimSpace(output) != "" {
		args = append(args, "-o", output)
	}
	if strings.TrimSpace(format) != "" {
		args = append(args, "-f", format)
	}
	err := executeCmd(binaryPath, args)
	if err != nil {
		return err
	}
	return nil
}

func BuildAndRunExecute(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	config, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("invalid value for --config: %w", err)
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("invalid value for --output: %w", err)
	}
	noexec, err := cmd.Flags().GetBool("noexec")
	if err != nil {
		return fmt.Errorf("invalid value for --noexec: %w", err)
	}

	// Target directory where generated code and binary will live
	outDir := filepath.Join(output, "target")
	if err := findAndTranspileDatagenModels(outDir, inputPath); err != nil {
		return err
	}
	if !noexec {
		if err := invokeExecute(outDir, output, config, inputPath); err != nil {
			return err
		}
	}
	return nil
}

func invokeExecute(outDir string, output string, config string, inputPath string) error {
	binaryPath, err := buildTranspiledBinary(outDir)
	if err != nil {
		return nil
	}

	args := []string{"execute", inputPath}
	args = append(args, "-c", config)
	if strings.TrimSpace(output) != "" {
		args = append(args, "-o", output)
	}
	err = executeCmd(binaryPath, args)
	if err != nil {
		return err
	}
	return nil
}

func findAndTranspileDatagenModels(outDir string, inputPath string) error {
	dgDirData, err := GetDgDirStructure(inputPath, "")
	if err != nil {
		return err
	}

	dgModelsCount := dgDirData.ModelCount()

	if dgModelsCount == 0 {
		return errors.New("No .dg files found in " + inputPath)
	}

	fmt.Printf("Found %d .dg files to process\n", dgModelsCount)
	parsedAll, err := processDgDirData(dgDirData, outDir, []*codegen.DatagenParsed{})

	if err != nil {
		return err
	}

	if err := codegen.Codegen(parsedAll, outDir, dgDirData); err != nil {
		return err
	}
	return nil
}

func processDgDirData(d *utils.DgDir, outDir string, accumulatedParsed []*codegen.DatagenParsed) ([]*codegen.DatagenParsed, error) {
	if d == nil {
		return accumulatedParsed, nil
	}

	parsed, err := transpile(d)
	if err != nil {
		return nil, err
	}

	for _, child := range d.Children {
		accumulatedChildren, err := processDgDirData(child, outDir, []*codegen.DatagenParsed{})
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, accumulatedChildren...)
	}

	return parsed, nil
}

func transpile(dgDirData *utils.DgDir) ([]*codegen.DatagenParsed, error) {
	var parsedResults []*codegen.DatagenParsed

	for path, src := range dgDirData.Models {
		result, err := parser.Parse(src, path)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", path, err)
			return nil, err
		}

		result.FullyQualifiedModelName = path

		if err := validators.Validate(result); err != nil {
			fmt.Printf("Validation failed for model %s: %v\n", result.ModelName, err)
			return nil, err
		}

		parsedResults = append(parsedResults, result)
	}

	return parsedResults, nil
}

func buildTranspiledBinary(outDir string) (string, error) {
	binaryName := "datagen"
	if runtime.GOOS == "windows" && filepath.Ext(binaryName) == "" {
		binaryName += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-C", outDir, "-o", binaryName)
	buildCmd.Stdout, buildCmd.Stderr, buildCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	if err := buildCmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %w", err)
	}

	binaryPath := filepath.Join(outDir, binaryName)

	if !filepath.IsAbs(binaryPath) {
		binaryPath = "./" + binaryPath
	}
	return binaryPath, nil
}

func executeCmd(binaryPath string, args []string) error {
	runCmd := exec.CommandContext(context.Background(), binaryPath, args...)
	runCmd.Stdout, runCmd.Stderr, runCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return runCmd.Run()
}
