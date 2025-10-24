package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"log/slog"

	"github.com/ds-horizon/datagen/codegen"
	"github.com/ds-horizon/datagen/parser"
	"github.com/ds-horizon/datagen/utils"
	"github.com/ds-horizon/datagen/validators"
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
	seed, err := cmd.Flags().GetInt64("seed")
	if err != nil {
		return fmt.Errorf("invalid value for --seed: %w", err)
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

    verbose, err := cmd.Flags().GetBool("verbose")
    if err != nil {
        return fmt.Errorf("invalid value for --verbose: %w", err)
    }

	if !noexec {
		if err := invokeGen(outDir, count, tags, output, format, seed, inputPath, verbose); err != nil {
			return err
		}
	}

	return nil
}

func invokeGen(outDir string, count int, tags string, output string, format string, seed int64, inputPath string, verbose bool) error {
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
	if seed != 0 {
		args = append(args, "--seed", fmt.Sprintf("%d", seed))
	}
    if verbose {
        args = append(args, "-v")
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
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return fmt.Errorf("invalid value for --verbose: %w", err)
	}

	// Target directory where generated code and binary will live
	outDir := filepath.Join(output, "target")
	if err := findAndTranspileDatagenModels(outDir, inputPath); err != nil {
		return err
	}
	if !noexec {
		if err := invokeExecute(outDir, output, config, inputPath, verbose); err != nil {
			return err
		}
	}
	return nil
}

func invokeExecute(outDir string, output string, config string, inputPath string, verbose bool) error {
	binaryPath, err := buildTranspiledBinary(outDir)
	if err != nil {
		return nil
	}

	args := []string{"execute", inputPath}
	args = append(args, "-c", config)
	if strings.TrimSpace(output) != "" {
		args = append(args, "-o", output)
	}
	if verbose {
		args = append(args, "-v")
	}

	err = executeCmd(binaryPath, args)
	if err != nil {
		return err
	}
	return nil
}

func findAndTranspileDatagenModels(outDir string, inputPath string) error {
    slog.Debug(fmt.Sprintf("finding and transpiling datagen models from %s into %s", inputPath, outDir))

	dgDirData, err := GetDgDirStructure(inputPath, "")
	if err != nil {
		return fmt.Errorf("failed to read input\n  input_path: %s\n  cause: %w", inputPath, err)
	}

	dgModelsCount := dgDirData.ModelCount()
    slog.Debug(fmt.Sprintf("found %d datagen models in %s", dgModelsCount, inputPath))

	if dgModelsCount == 0 {
		slog.Warn("no .dg files found", "input_path", inputPath)
		return fmt.Errorf("no .dg files found in %s", inputPath)
	}

	parsedAll, err := processDgDirData(dgDirData, outDir, []*codegen.DatagenParsed{})
	if err != nil {
		return fmt.Errorf("failed to process directory data\n  input_path: %s\n  cause: %w", inputPath, err)
	}

    slog.Debug(fmt.Sprintf("generating code for %d models into %s", len(parsedAll), outDir))
	if err := codegen.Codegen(parsedAll, outDir, dgDirData); err != nil {
		return fmt.Errorf("code generation failed\n  output_dir: %s\n  cause: %w", outDir, err)
	}

    slog.Info(fmt.Sprintf("successfully transpiled %d datagen models into %s", len(parsedAll), outDir))
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
        slog.Debug(fmt.Sprintf("parsing model file %s", path))

		result, err := parser.Parse(src, path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse model\n  file: %s\n  cause: %w", path, err)
		}

		result.FullyQualifiedModelName = path

		if err := validators.Validate(result); err != nil {
			return nil, fmt.Errorf("model validation failed\n  model: %s\n  file: %s\n  cause: %w", result.ModelName, path, err)
		}

    slog.Debug(fmt.Sprintf("parsed and validated model %s from %s", result.ModelName, path))
		parsedResults = append(parsedResults, result)
	}

	return parsedResults, nil
}

func buildTranspiledBinary(outDir string) (string, error) {
	binaryName := "datagen"
	if runtime.GOOS == "windows" && filepath.Ext(binaryName) == "" {
		binaryName += ".exe"
	}

    slog.Debug(fmt.Sprintf("building transpiled binary %s in %s", binaryName, outDir))

	buildCmd := exec.Command("go", "build", "-C", outDir, "-o", binaryName)
	buildCmd.Stdout, buildCmd.Stderr, buildCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	if err := buildCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build transpiled binary: %s\n  cause: %w", outDir, err)
	}

	binaryPath := filepath.Join(outDir, binaryName)

	if !filepath.IsAbs(binaryPath) {
		binaryPath = "./" + binaryPath
	}

    slog.Debug(fmt.Sprintf("built transpiled binary at %s", binaryPath))
	return binaryPath, nil
}

func logCapturedOutput(stdout, stderr strings.Builder) {
	if stdout.Len() > 0 {
		logOutputLines(stdout.String(), "command stdout")
	} else {
        slog.Debug("no stdout captured")
	}

	if stderr.Len() > 0 {
		logOutputLines(stderr.String(), "command stderr")
	} else {
        slog.Debug("no stderr captured")
	}
}

func logOutputLines(output, logKey string) {
	if strings.TrimSpace(output) == "" {
		return
	}

	if logKey == "command stdout" {
		fmt.Println(output)
	} else {
		fmt.Fprintln(os.Stderr, output)
	}
}

func executeCmd(binaryPath string, args []string) error {
    slog.Debug(fmt.Sprintf("executing command: %s %s", binaryPath, strings.Join(args, " ")))

	runCmd := exec.CommandContext(context.Background(), binaryPath, args...)

	var stdout, stderr strings.Builder
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr
	runCmd.Stdin = os.Stdin

	if err := runCmd.Run(); err != nil {
		logCapturedOutput(stdout, stderr)
		return err
	}

	logCapturedOutput(stdout, stderr)

    slog.Debug(fmt.Sprintf("command executed successfully: %s", binaryPath))
	return nil
}
