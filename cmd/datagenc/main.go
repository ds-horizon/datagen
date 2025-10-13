package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dream-sports-labs/datagen/runner"
	"github.com/spf13/cobra"
)

var (
	flagCount  int
	flagTags   string
	flagOutput string
	flagFormat string
	flagNoExec bool
	flagConfig string
	version    = "dev"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "datagen",
		Short:   "Transpile .dg model files into a runnable data generator",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Silence usage for all commands by default
			// Usage will only be shown for flag-related errors
			cmd.SilenceUsage = true
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	genCmd := &cobra.Command{
		Use:   "gen [file|directory]",
		Short: "Transpile models, and optionally invoke datagen",
		Args:  validateSingleFileOrDir,
		RunE:  runner.BuildAndRunGen,
	}

	genCmd.Flags().IntVarP(&flagCount, "count", "n", -1, "number of records per model")
	genCmd.Flags().StringVarP(&flagTags, "tags", "t", "", "comma-separated key=value tags to filter models")
	genCmd.Flags().StringVarP(&flagOutput, "output", "o", ".", "output directory or file path")
	genCmd.Flags().StringVarP(&flagFormat, "format", "f", "", strings.Join([]string{"csv", "json", "xml", "stdout"}, "|"))
	genCmd.Flags().BoolVar(&flagNoExec, "noexec", false, "skip building and executing generated binary")

	rootCmd.AddCommand(genCmd)

	executeCmd := &cobra.Command{
		Use:   "execute [file|directory]",
		Short: "Transpile models, invoke datagen and load data into relevant sinks",
		Args:  validateSingleFileOrDir,
		RunE:  runner.BuildAndRunExecute,
	}
	executeCmd.Flags().StringVarP(&flagConfig, "config", "c", "", "Path of config.json file to be passed to datagen")
	_ = executeCmd.MarkFlagRequired("config")
	executeCmd.Flags().StringVarP(&flagOutput, "output", "o", ".", "output directory or file path")
	executeCmd.Flags().BoolVar(&flagNoExec, "noexec", false, "skip building and executing generated binary")

	rootCmd.AddCommand(executeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func validateSingleFileOrDir(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide a file or directory")
	}
	if len(args) > 1 {
		return fmt.Errorf("requires exactly one file or directory path, received %d", len(args))
	}
	return nil
}
