package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nouuu/gopeek/internal/ignore"
	"github.com/nouuu/gopeek/internal/logger"
)

type Scanner struct {
	rootDir       string
	config        Config
	output        Output
	ignoreMatcher *ignore.Matcher
	log           *logger.Logger
}

func New(rootDir string, config Config, log *logger.Logger) *Scanner {
	ignoreList := ignore.NewMatcher()

	for _, pattern := range config.IgnorePatterns {
		ignoreList.AddPattern(pattern)
	}

	gitignorePath := filepath.Join(rootDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		log.Debug("loading .gitignore file", "path", gitignorePath)
		if err = ignoreList.LoadFile(gitignorePath); err != nil {
			log.Warn("error loading .gitignore file", "error", err)
		}
	}

	return &Scanner{
		rootDir: rootDir,
		config:  config,
		output: Output{
			log: log,
		},
		ignoreMatcher: ignoreList,
		log:           log,
	}
}

func (s *Scanner) Run() error {
	if err := s.scan(); err != nil {
		return fmt.Errorf("scanning error: %w", err)
	}

	output := s.output.Generate()
	s.log.Info("writing output", "file", s.config.Output)
	return os.WriteFile(s.config.Output, []byte(output), 0o644)
}

func (s *Scanner) scan() error {
	return filepath.Walk(s.rootDir, s.processPath)
}

func (s *Scanner) processPath(path string, info fs.FileInfo, err error) error {
	if err != nil {
		s.log.Debug("error accessing path", "path", path, "error", err)
		return fmt.Errorf("error accessing path %s: %w", path, err)
	}

	s.log.Debug("processing path", "path", path)

	if s.shouldIgnore(path) {
		s.log.Debug("ignoring path", "path", path)
		return skipIfDir(info)
	}

	relPath, err := filepath.Rel(s.rootDir, path)
	if err != nil {
		s.log.Debug("error getting relative path", "path", path, "error", err)
		return fmt.Errorf("error getting relative path for %s: %w", path, err)
	}

	if relPath == "." {
		return nil
	}

	depth := strings.Count(relPath, string(os.PathSeparator))
	s.output.AddStructure(path, info, depth)

	if !info.IsDir() {
		if err := s.output.AddContent(path, relPath); err != nil {
			s.log.Warn("error adding content", "path", path, "error", err)
		}
	}

	return nil
}

func skipIfDir(info fs.FileInfo) error {
	if info.IsDir() {
		return fs.SkipDir
	}
	return nil
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
