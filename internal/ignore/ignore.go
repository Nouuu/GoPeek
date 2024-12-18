package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Pattern struct {
	pattern string
	negate  bool
}

type IgnoreList struct {
	patterns []Pattern
}

func New() *IgnoreList {
	return &IgnoreList{}
}

func (i *IgnoreList) AddPattern(pattern string) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" || pattern[0] == '#' {
		return
	}

	negate := false
	if pattern[0] == '!' {
		negate = true
		pattern = pattern[1:]
	}

	if slices.ContainsFunc(i.patterns, func(p Pattern) bool {
		return p.pattern == pattern && p.negate == negate
	}) {
		return
	}

	i.patterns = append(i.patterns, Pattern{pattern: pattern, negate: negate})
}

func (i *IgnoreList) LoadGitignore(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i.AddPattern(scanner.Text())
	}
	return scanner.Err()
}

func (i *IgnoreList) ShouldIgnore(path string) bool {
	path = filepath.ToSlash(path)
	ignored := false

	for _, pattern := range i.patterns {
		if match, _ := filepath.Match(pattern.pattern, path); match {
			ignored = !pattern.negate
		}
		// Support directory wilcards
		dirPattern := strings.ReplaceAll(pattern.pattern, "**/", "")
		if match, _ := filepath.Match(dirPattern, filepath.Base(path)); match {
			ignored = !pattern.negate
		}
	}
	return ignored
}
