package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Pattern struct {
	raw      string
	segments []string
	negate   bool
	matchAll bool // For ** pattern
}

type Matcher struct {
	patterns []Pattern
	cache    sync.Map
}

func NewMatcher() *Matcher {
	return &Matcher{
		patterns: make([]Pattern, 0),
	}
}

func (m *Matcher) AddPattern(raw string) {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.HasPrefix(raw, "#") {
		return
	}
	pattern := Pattern{raw: raw, segments: make([]string, 0)}
	if strings.HasPrefix(raw, "!") {
		pattern.negate = true
		raw = raw[1:]
	}
	raw = filepath.ToSlash(raw)

	raw = strings.Trim(raw, "/")
	if strings.Contains(raw, "**") {
		pattern.matchAll = true
		raw = strings.ReplaceAll(raw, "**/", "*")
		raw = strings.ReplaceAll(raw, "**", "*")
	}

	if raw != "" {
		pattern.segments = strings.Split(raw, "/")
	}
	m.patterns = append(m.patterns, pattern)
}

func (m *Matcher) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m.AddPattern(scanner.Text())
	}
	return scanner.Err()
}

func (m *Matcher) ShouldIgnore(path string) bool {
	if res, ok := m.cache.Load(path); ok {
		return res.(bool)
	}

	normalizedPath := filepath.ToSlash(path)
	segments := strings.Split(normalizedPath, "/")

	ignored := false
	for _, pattern := range m.patterns {
		if match := pattern.matches(segments); match {
			ignored = !pattern.negate
		}
	}

	m.cache.Store(path, ignored)
	return ignored
}

func (p *Pattern) matches(pathSegments []string) bool {
	if p.matchAll {
		return matchWildcard(p.segments, pathSegments)
	}
	return matchExact(p.segments, pathSegments)
}

func matchExact(pattern, path []string) bool {
	if len(pattern) == 0 {
		return len(path) == 0
	}

	if len(path) == 0 {
		return false
	}

	if pattern[0] == "*" {
		return matchExact(pattern[1:], path) || matchExact(pattern[1:], path[1:]) || matchExact(pattern, path[1:])
	}

	if matched, _ := filepath.Match(pattern[0], path[0]); matched {
		return matchExact(pattern[1:], path[1:])
	}

	return false
}

func matchWildcard(pattern, path []string) bool {
	for i := 0; i < len(path); i++ {
		if matchExact(pattern, path[i:]) {
			return true
		}
	}
	return false
}
