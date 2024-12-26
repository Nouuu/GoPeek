package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nouuu/gopeek/internal/ignore"
)

type Scanner struct {
	rootDir       string
	config        Config
	output        Output
	ignoreMatcher *ignore.Matcher
}

func New(rootDir string, config Config) *Scanner {
	ignoreList := ignore.NewMatcher()

	for _, pattern := range config.IgnorePatterns {
		ignoreList.AddPattern(pattern)
	}

	gitignorePath := filepath.Join(rootDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		if err = ignoreList.LoadFile(gitignorePath); err != nil {
			fmt.Printf("Warning: %v\n", err)
		}
	}

	return &Scanner{
		rootDir:       rootDir,
		config:        config,
		output:        Output{},
		ignoreMatcher: ignoreList,
	}
}

func (s *Scanner) Run() error {
	if err := s.scan(); err != nil {
		return fmt.Errorf("scanning error: %w", err)
	}

	output := s.output.Generate()
	return os.WriteFile(s.config.Output, []byte(output), 0o644)
}

func (s *Scanner) scan() error {
	return filepath.Walk(s.rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if s.shouldIgnore(path) {
			if info.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(s.rootDir, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %w", path, err)
		}

		if relPath == "." {
			return nil
		}

		depth := len(strings.Split(relPath, string(os.PathSeparator))) - 1
		s.output.AddStructure(path, info, depth)

		if !info.IsDir() {
			if err := s.output.AddContent(path, relPath); err != nil {
				fmt.Printf("Warning: %v\n", err)
			}
		}

		return nil
	})
}

func (s *Scanner) shouldIgnore(path string) bool {
	// Ignore output file
	if filepath.Clean(path) == filepath.Clean(s.config.Output) {
		return true
	}

	relPath, err := filepath.Rel(s.rootDir, path)
	if err != nil {
		return false
	}

	return s.ignoreMatcher.ShouldIgnore(relPath)
}
