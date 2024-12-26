package ignore

import (
	"os"
	"testing"
)

func TestMatcher_AddPattern(t *testing.T) {
	tests := []struct {
		name          string
		pattern       string
		expectPattern bool
		expected      Pattern
	}{
		{
			name:          "Empty pattern",
			pattern:       "",
			expectPattern: false,
		},
		{
			name:          "Comment pattern",
			pattern:       "# comment",
			expectPattern: false,
		},
		{
			name:          "Simple pattern",
			pattern:       "build",
			expectPattern: true,
			expected: Pattern{
				raw:      "build",
				segments: []string{"build"},
			},
		},
		{
			name:          "Pattern with slashes",
			pattern:       "/build/",
			expectPattern: true,
			expected: Pattern{
				raw:      "/build/",
				segments: []string{"build"},
			},
		},
		{
			name:          "Pattern with **",
			pattern:       "**/test",
			expectPattern: true,
			expected: Pattern{
				raw:      "**/test",
				segments: []string{"test"},
				matchAll: true,
			},
		},
		{
			name:          "Negated pattern",
			pattern:       "!build",
			expectPattern: true,
			expected: Pattern{
				raw:      "!build",
				segments: []string{"build"},
				negate:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatcher()
			m.AddPattern(tt.pattern)

			if tt.expectPattern {
				if len(m.patterns) != 1 {
					t.Fatal("Expected one pattern")
				}
				pattern := m.patterns[0]
				if pattern.negate != tt.expected.negate {
					t.Errorf("negate = %v, want %v", pattern.negate, tt.expected.negate)
				}
				if pattern.matchAll != tt.expected.matchAll {
					t.Errorf("matchAll = %v, want %v", pattern.matchAll, tt.expected.matchAll)
				}
				if len(pattern.segments) != len(tt.expected.segments) {
					t.Errorf("segments length = %v, want %v", len(pattern.segments), len(tt.expected.segments))
				}
			} else {
				if len(m.patterns) != 0 {
					t.Error("Expected no patterns")
				}
			}
		})
	}
}

func TestMatcher_LoadFile(t *testing.T) {
	content := `
# Comment
*.go
/build/
!important.go
**/test
`
	tmpfile, err := os.CreateTemp("", "ignore")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	m := NewMatcher()
	if err := m.LoadFile(tmpfile.Name()); err != nil {
		t.Fatal(err)
	}

	expectedPatterns := 4 // Exclusion du commentaire
	if len(m.patterns) != expectedPatterns {
		t.Errorf("Expected %d patterns, got %d", expectedPatterns, len(m.patterns))
	}
}

func TestMatcher_ShouldIgnore(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		{
			name:     "Simple match",
			patterns: []string{"build"},
			path:     "build",
			want:     true,
		},
		{
			name:     "Directory match",
			patterns: []string{"build/*"},
			path:     "build/output",
			want:     true,
		},
		{
			name:     "Wildcard match",
			patterns: []string{"*.go"},
			path:     "main.go",
			want:     true,
		},
		{
			name:     "No match",
			patterns: []string{"build"},
			path:     "src",
			want:     false,
		},
		{
			name:     "Negated pattern",
			patterns: []string{"*.go", "!main.go"},
			path:     "main.go",
			want:     false,
		},
		{
			name:     "Double wildcard",
			patterns: []string{"**/test"},
			path:     "a/b/test",
			want:     true,
		},
		{
			name:     "Double wildcard with extension",
			patterns: []string{"**/*.go"},
			path:     "src/pkg/main.go",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatcher()
			for _, pattern := range tt.patterns {
				m.AddPattern(pattern)
			}

			if got := m.ShouldIgnore(tt.path); got != tt.want {
				t.Errorf("ShouldIgnore() = %v, want %v", got, tt.want)
			}

			// Test cache
			if got := m.ShouldIgnore(tt.path); got != tt.want {
				t.Errorf("ShouldIgnore() (cached) = %v, want %v", got, tt.want)
			}
		})
	}
}
