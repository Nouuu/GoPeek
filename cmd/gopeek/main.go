package main

import (
	"fmt"
	"strings"

	"github.com/nouuu/gopeek/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "gopeek [path]",
	Example: "gopeek ./",
	Short:   "Scan a project directory and output its structure and content",
	Version: formatVersion(),
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := scanner.DefaultConfig()
		if ignore, _ := cmd.Flags().GetStringSlice("ignore"); len(ignore) > 0 {
			cfg.IgnorePatterns = ignore
		}
		outputFile, _ := cmd.Flags().GetString("output")
		cfg.Output = outputFile

		s := scanner.New(args[0], cfg)
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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
