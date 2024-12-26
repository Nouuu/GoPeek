package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nouuu/gopeek/internal/logger"
)

func TestNew(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test .gitignore file
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("*.tmp\n"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		config         Config
		expectPatterns int
	}{
		{
			name: "Default configuration",
			config: Config{
				Output:         "output.md",
				IgnorePatterns: DefaultIgnorePatterns,
			},
			expectPatterns: len(DefaultIgnorePatterns) + 1, // +1 for .gitignore pattern
		},
		{
			name: "Custom ignore patterns",
			config: Config{
				Output:         "output.md",
				IgnorePatterns: []string{"*.log", "temp/"},
			},
			expectPatterns: 3, // 2 custom patterns + 1 from .gitignore
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.Default()
			scanner := New(tmpDir, tt.config, log)

			if scanner == nil {
				t.Fatal("Expected scanner to be created, got nil")
			}

			if scanner.config.Output != tt.config.Output {
				t.Errorf("Expected output %s, got %s", tt.config.Output, scanner.config.Output)
			}

			if len(scanner.ignoreMatcher.Patterns()) != tt.expectPatterns {
				t.Errorf("Expected %d ignore patterns, got %d", tt.expectPatterns, len(scanner.ignoreMatcher.Patterns()))
			}
		})
	}
}

func TestScanner_Run(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directory structure
	testFiles := map[string]string{
		"file1.txt":           "Content 1",
		"dir1/file2.txt":      "Content 2",
		"dir1/dir2/file3.txt": "Content 3",
		".gitignore":          "*.ignore",
		"ignored.ignore":      "Should be ignored",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name              string
		config            Config
		expectError       bool
		expectInOutput    []string
		expectNotInOutput []string
	}{
		{
			name: "Basic scan",
			config: Config{
				Output:         filepath.Join(tmpDir, "output.md"),
				IgnorePatterns: []string{},
			},
			expectError: false,
			expectInOutput: []string{
				"file1.txt",
				"dir1/file2.txt",
				"dir1/dir2/file3.txt",
			},
		},
		{
			name: "With ignore patterns",
			config: Config{
				Output:         filepath.Join(tmpDir, "output.md"),
				IgnorePatterns: []string{"*.ignore"},
			},
			expectError: false,
			expectInOutput: []string{
				"file1.txt",
			},
			expectNotInOutput: []string{
				"ignored.ignore",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.Default()
			scanner := New(tmpDir, tt.config, log)
			err := scanner.Run()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if err == nil {
				// Check if output file exists
				content, err := os.ReadFile(tt.config.Output)
				if err != nil {
					t.Fatal(err)
				}

				// Check expected content
				for _, expect := range tt.expectInOutput {
					if !strings.Contains(string(content), expect) {
						t.Errorf("Expected output to contain %q", expect)
					}
				}

				// Check not expected content
				for _, notExpect := range tt.expectNotInOutput {
					if strings.Contains(string(content), notExpect) {
						t.Errorf("Expected output to not contain %q", notExpect)
					}
				}
			}
		})
	}
}

func TestScanner_shouldIgnore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		config   Config
		path     string
		expected bool
	}{
		{
			name: "Ignore output file",
			config: Config{
				Output:         filepath.Join(tmpDir, "output.md"),
				IgnorePatterns: []string{},
			},
			path:     filepath.Join(tmpDir, "output.md"),
			expected: true,
		},
		{
			name: "Ignore pattern match",
			config: Config{
				Output:         filepath.Join(tmpDir, "output.md"),
				IgnorePatterns: []string{"*.log"},
			},
			path:     filepath.Join(tmpDir, "test.log"),
			expected: true,
		},
		{
			name: "Don't ignore normal file",
			config: Config{
				Output:         filepath.Join(tmpDir, "output.md"),
				IgnorePatterns: []string{"*.log"},
			},
			path:     filepath.Join(tmpDir, "test.txt"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.Default()
			scanner := New(tmpDir, tt.config, log)
			result := scanner.shouldIgnore(tt.path)
			if result != tt.expected {
				t.Errorf("shouldIgnore(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestScanner_processPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scanner-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		setupPath   string
		testPath    string
		expectError bool
	}{
		{
			name:        "Process existing file",
			setupPath:   testFile,
			testPath:    filepath.Join(tmpDir, "test.txt"),
			expectError: false,
		},
		{
			name:        "Process non-existent file",
			setupPath:   testFile,
			testPath:    filepath.Join(tmpDir, "nonexistent.txt"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.Default()
			scanner := New(tmpDir, DefaultConfig(), log)

			var fileInfo os.FileInfo
			var statErr error

			fileInfo, statErr = os.Stat(tt.testPath)
			err = scanner.processPath(tt.testPath, fileInfo, statErr)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
