package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Scanner struct {
	rootDir string
	config  Config
	output  Output
}

func New(rootDir string, config Config) *Scanner {
	return &Scanner{
		rootDir: rootDir,
		config:  config,
		output:  Output{},
	}
}

func (s *Scanner) Run() error {
	if err := s.scan(); err != nil {
		return fmt.Errorf("scanning error: %w", err)
	}

	output := s.output.Generate()
	return os.WriteFile(s.config.Output, []byte(output), 0644)
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
	for _, pattern := range s.config.IgnorePatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}
