package main

import (
	"github.com/nouuu/gopeek/internal/scanner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gopeek [path]",
	Example: "gopeek ./",
	Short:   "Scan a project directory and output its structure and content",
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

func init() {
	rootCmd.Flags().StringP("output", "o", "project_knowledge.md", "Output file")
	rootCmd.Flags().StringSliceP("ignore", "i", []string{}, "Patterns to ignore")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
