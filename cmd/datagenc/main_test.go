package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ds-horizon/datagen/utils"
)

const (
	expectedRootHelp = `Generate realistic test data from model definitions

Usage:
  datagenc [command]

Available Commands:
  execute     Generate data from .dg model files and load into configured data stores
  gen         Generate data from .dg model files and output to CSV, JSON, XML, or stdout
  help        Help about any command

Flags:
  -h, --help      help for datagenc
  -v, --verbose   enable verbose (debug level) logging
  -V, --version   show version information

Use "datagenc [command] --help" for more information about a command.
`

	expectedGenHelp = `Generate data from .dg model files and output to CSV, JSON, XML, or stdout

Usage:
  datagenc gen [file|directory] [flags]

Flags:
  -n, --count int       number of records per model (default -1)
  -f, --format string   csv|json|xml|stdout
  -h, --help            help for gen
      --noexec          skip building and executing generated binary
  -o, --output string   output directory or file path (default ".")
  -s, --seed int        deterministic seed for random data generation (default is 0 for random seed)
  -t, --tags string     comma-separated key=value tags to filter models

Global Flags:
  -v, --verbose   enable verbose (debug level) logging
  -V, --version   show version information
`

	expectedExecuteHelp = `Generate data from .dg model files and load into configured data stores

Usage:
  datagenc execute [file|directory] [flags]

Flags:
  -c, --config string   path to config file (specifies models, data stores, and record counts)
  -h, --help            help for execute
      --noexec          skip building and executing generated binary
  -o, --output string   output directory or file path (default ".")

Global Flags:
  -v, --verbose   enable verbose (debug level) logging
  -V, --version   show version information
`
)

func TestValidateSingleFileOrDir(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantErr   bool
		errSubstr string
	}{
		{
			name:      "valid single file",
			args:      []string{"test.dg"},
			wantErr:   false,
			errSubstr: "",
		},
		{
			name:      "valid single directory",
			args:      []string{"./models"},
			wantErr:   false,
			errSubstr: "",
		},
		{
			name:      "valid absolute path",
			args:      []string{"/path/to/file.dg"},
			wantErr:   false,
			errSubstr: "",
		},
		{
			name:      "empty args",
			args:      []string{},
			wantErr:   true,
			errSubstr: "must provide a file or directory",
		},
		{
			name:      "too many args - two files",
			args:      []string{"file1.dg", "file2.dg"},
			wantErr:   true,
			errSubstr: "requires exactly one file or directory path, received 2",
		},
		{
			name:      "too many args - three files",
			args:      []string{"file1.dg", "file2.dg", "file3.dg"},
			wantErr:   true,
			errSubstr: "requires exactly one file or directory path, received 3",
		},
		{
			name:      "too many args - mixed paths",
			args:      []string{"./models", "./other"},
			wantErr:   true,
			errSubstr: "requires exactly one file or directory path, received 2",
		},
		{
			name:      "single empty string",
			args:      []string{""},
			wantErr:   false,
			errSubstr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := validateSingleFileOrDir(cmd, tt.args)

			if tt.wantErr {
				require.Error(t, err, "expected error for args: %v", tt.args)
				if tt.errSubstr != "" {
					assert.Contains(t, err.Error(), tt.errSubstr,
						"error message should contain %q, got: %v", tt.errSubstr, err)
				}
				return
			}

			require.NoError(t, err, "unexpected error for args: %v", tt.args)
		})
	}
}

func TestBuildRootCommand(t *testing.T) {
	rootCmd := buildRootCommand()

	assert.Equal(t, utils.CompilerBinaryName, rootCmd.Use)
	assert.Equal(t, "Generate realistic test data from model definitions", rootCmd.Short)
	assert.Equal(t, version, rootCmd.Version)
	assert.True(t, rootCmd.CompletionOptions.DisableDefaultCmd)

	commands := rootCmd.Commands()
	assert.Len(t, commands, 2, "expected 2 subcommands")

	var genCmd, executeCmd *cobra.Command
	for _, cmd := range commands {
		switch cmd.Use {
		case "gen [file|directory]":
			genCmd = cmd
		case "execute [file|directory]":
			executeCmd = cmd
		default:
			t.Fatalf("unexpected command found: %q", cmd.Use)
		}
	}

	require.NotNil(t, genCmd, "gen command should exist")
	require.NotNil(t, executeCmd, "execute command should exist")

	assert.Equal(t, "Generate data from .dg model files and output to CSV, JSON, XML, or stdout", genCmd.Short)
	assert.NotNil(t, genCmd.Args, "gen command should have Args validator")
	assert.NotNil(t, genCmd.RunE, "gen command should have RunE handler")

	assert.Equal(t, "Generate data from .dg model files and load into configured data stores", executeCmd.Short)
	assert.NotNil(t, executeCmd.Args, "execute command should have Args validator")
	assert.NotNil(t, executeCmd.RunE, "execute command should have RunE handler")
}

func TestRootCommandFlags(t *testing.T) {
	rootCmd := buildRootCommand()

	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	require.NotNil(t, verboseFlag, "verbose flag should exist")
	assert.Equal(t, "v", verboseFlag.Shorthand)
	assert.Equal(t, "false", verboseFlag.DefValue)

	versionFlag := rootCmd.PersistentFlags().Lookup("version")
	require.NotNil(t, versionFlag, "version flag should exist")
	assert.Equal(t, "V", versionFlag.Shorthand)
	assert.Equal(t, "false", versionFlag.DefValue)
}

func TestGenCommandFlags(t *testing.T) {
	rootCmd := buildRootCommand()
	genCmd, _, err := rootCmd.Find([]string{"gen"})
	require.NoError(t, err)
	require.NotNil(t, genCmd)

	countFlag := genCmd.Flags().Lookup("count")
	require.NotNil(t, countFlag, "count flag should exist")
	assert.Equal(t, "n", countFlag.Shorthand)
	assert.Equal(t, "-1", countFlag.DefValue)

	tagsFlag := genCmd.Flags().Lookup("tags")
	require.NotNil(t, tagsFlag, "tags flag should exist")
	assert.Equal(t, "t", tagsFlag.Shorthand)
	assert.Equal(t, "", tagsFlag.DefValue)

	outputFlag := genCmd.Flags().Lookup("output")
	require.NotNil(t, outputFlag, "output flag should exist")
	assert.Equal(t, "o", outputFlag.Shorthand)
	assert.Equal(t, ".", outputFlag.DefValue)

	formatFlag := genCmd.Flags().Lookup("format")
	require.NotNil(t, formatFlag, "format flag should exist")
	assert.Equal(t, "f", formatFlag.Shorthand)
	assert.Equal(t, "", formatFlag.DefValue)

	seedFlag := genCmd.Flags().Lookup("seed")
	require.NotNil(t, seedFlag, "seed flag should exist")
	assert.Equal(t, "s", seedFlag.Shorthand)
	assert.Equal(t, "0", seedFlag.DefValue)

	noexecFlag := genCmd.Flags().Lookup("noexec")
	require.NotNil(t, noexecFlag, "noexec flag should exist")
	assert.Equal(t, "false", noexecFlag.DefValue)
}

func TestExecuteCommandFlags(t *testing.T) {
	rootCmd := buildRootCommand()
	executeCmd, _, err := rootCmd.Find([]string{"execute"})
	require.NoError(t, err)
	require.NotNil(t, executeCmd)

	configFlag := executeCmd.Flags().Lookup("config")
	require.NotNil(t, configFlag, "config flag should exist")
	assert.Equal(t, "c", configFlag.Shorthand)
	assert.Equal(t, "", configFlag.DefValue)
	err = executeCmd.ValidateRequiredFlags()
	if err != nil {
		assert.Contains(t, err.Error(), "config", "config flag should be required")
	}

	outputFlag := executeCmd.Flags().Lookup("output")
	require.NotNil(t, outputFlag, "output flag should exist")
	assert.Equal(t, "o", outputFlag.Shorthand)
	assert.Equal(t, ".", outputFlag.DefValue)

	noexecFlag := executeCmd.Flags().Lookup("noexec")
	require.NotNil(t, noexecFlag, "noexec flag should exist")
	assert.Equal(t, "false", noexecFlag.DefValue)
}

func TestCommandExecution(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantErr   bool
		errSubstr string
	}{
		{
			name:      "gen command with valid file",
			args:      []string{"gen", "test.dg"},
			wantErr:   true,
			errSubstr: "",
		},
		{
			name:      "unknown command",
			args:      []string{"unknown"},
			wantErr:   true,
			errSubstr: "unknown",
		},
		{
			name:      "no command provided",
			args:      []string{},
			wantErr:   false,
			errSubstr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := buildRootCommand()
			rootCmd.SetArgs(tt.args)
			rootCmd.SetOut(&bytes.Buffer{})
			rootCmd.SetErr(&bytes.Buffer{})

			err := rootCmd.Execute()

			if tt.wantErr {
				if tt.errSubstr != "" {
					require.Error(t, err, "expected error for args: %v", tt.args)
					assert.Contains(t, err.Error(), tt.errSubstr,
						"error message should contain %q, got: %v", tt.errSubstr, err)
				}
				return
			}
		})
	}
}

func TestVersionFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	rootCmd := buildRootCommand()
	versionFlag := rootCmd.PersistentFlags().Lookup("version")
	require.NotNil(t, versionFlag)

	err := rootCmd.PersistentFlags().Set("version", "true")
	require.NoError(t, err)

	val, err := rootCmd.PersistentFlags().GetBool("version")
	require.NoError(t, err)
	assert.True(t, val)
}

func TestCommandHelp(t *testing.T) {
	rootCmd := buildRootCommand()
	rootCmd.SetArgs([]string{"--help"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	assert.Equal(t, expectedRootHelp, out.String())
}

func TestGenCommandHelp(t *testing.T) {
	rootCmd := buildRootCommand()
	rootCmd.SetArgs([]string{"gen", "--help"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	assert.Equal(t, expectedGenHelp, out.String())
}

func TestExecuteCommandHelp(t *testing.T) {
	rootCmd := buildRootCommand()
	rootCmd.SetArgs([]string{"execute", "--help"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	assert.Equal(t, expectedExecuteHelp, out.String())
}
func TestHelpCommand(t *testing.T) {
	rootCmd := buildRootCommand()
	rootCmd.SetArgs([]string{"help"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	assert.Equal(t, expectedRootHelp, out.String())
}

func TestHelpGenCommand(t *testing.T) {
	rootCmd := buildRootCommand()
	rootCmd.SetArgs([]string{"help", "gen"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	assert.Equal(t, expectedGenHelp, out.String())
}

func TestHelpExecuteCommand(t *testing.T) {
	rootCmd := buildRootCommand()
	rootCmd.SetArgs([]string{"help", "execute"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&bytes.Buffer{})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	assert.Equal(t, expectedExecuteHelp, out.String())
}

func TestPersistentPreRun(t *testing.T) {
	rootCmd := buildRootCommand()

	assert.NotNil(t, rootCmd.PersistentPreRun)

	rootCmd.SetArgs([]string{"gen", "test.dg"})
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})

	_ = rootCmd.Execute()

	genCmd, _, err := rootCmd.Find([]string{"gen"})
	require.NoError(t, err)

	assert.NotNil(t, rootCmd.PersistentPreRun, "PersistentPreRun should be defined")
	assert.NotNil(t, genCmd, "gen command should exist")
}

func TestFlagDefaults(t *testing.T) {
	rootCmd := buildRootCommand()
	genCmd, _, err := rootCmd.Find([]string{"gen"})
	require.NoError(t, err)

	count, err := genCmd.Flags().GetInt("count")
	require.NoError(t, err)
	assert.Equal(t, -1, count)

	output, err := genCmd.Flags().GetString("output")
	require.NoError(t, err)
	assert.Equal(t, ".", output)

	seed, err := genCmd.Flags().GetInt64("seed")
	require.NoError(t, err)
	assert.Equal(t, int64(0), seed)

	noexec, err := genCmd.Flags().GetBool("noexec")
	require.NoError(t, err)
	assert.False(t, noexec)
}
