package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/nouuu/gopeek/internal/logger"
	"github.com/nouuu/gopeek/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	version  = "dev"
	commit   = "none"
	date     = "unknown"
	exitFunc = os.Exit // Makes the exit function replaceable for testing
)

var rootCmd = &cobra.Command{
	Use:     "gopeek [path]",
	Example: "gopeek ./",
	Short:   "Scan a project directory and output its structure and content",
	Version: formatVersion(),
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.Default()
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			log = log.WithLevel(slog.LevelDebug)
		}

		cfg := scanner.DefaultConfig()
		if ignore, _ := cmd.Flags().GetStringSlice("ignore"); len(ignore) > 0 {
			cfg.IgnorePatterns = ignore
			log.Debug("ignore patterns", "patterns", ignore)
		}
		outputFile, _ := cmd.Flags().GetString("output")
		cfg.Output = outputFile
		log.Debug("output file", "file", outputFile)

		s := scanner.New(args[0], cfg, log)
		return s.Run()
	},
}

func formatVersion() string {
	result := version
	if commit != "none" && date != "unknown" {
		if !strings.Contains(version, commit) {
			result += fmt.Sprintf("\ncommit: %s", commit)
		}
		result += fmt.Sprintf("\nbuilt at: %s", date)
	}
	return result
}

func init() {
	rootCmd.Flags().StringP("output", "o", "project_knowledge.md", "Output file")
	rootCmd.Flags().StringSliceP("ignore", "i", []string{}, "Patterns to ignore")
	rootCmd.Flags().Bool("verbose", false, "Verbose output")
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitFunc(1)
	}
}
