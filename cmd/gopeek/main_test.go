package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gopeek-cmd-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		args        []string
		expectError bool
		setup       func(*testing.T, *cobra.Command)
		validate    func(*testing.T, error)
	}{
		{
			name:        "No arguments",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "Valid path",
			args:        []string{tmpDir},
			expectError: false,
			validate: func(t *testing.T, err error) {
				outFile := "project_knowledge.md"
				if _, err := os.Stat(outFile); err != nil {
					t.Errorf("Expected output file to exist: %v", err)
				}
				defer os.Remove(outFile)
			},
		},
		{
			name:        "Custom output with verbose mode",
			args:        []string{tmpDir, "-o", filepath.Join(tmpDir, "custom.md"), "--verbose"},
			expectError: false,
			validate: func(t *testing.T, err error) {
				outFile := filepath.Join(tmpDir, "custom.md")
				if _, err := os.Stat(outFile); err != nil {
					t.Errorf("Expected output file to exist: %v", err)
				}
			},
		},
		{
			name:        "With ignore patterns",
			args:        []string{tmpDir, "-i", "*.log", "-i", "*.tmp"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			rootCmd.AddCommand(cmd)

			if tt.setup != nil {
				tt.setup(t, rootCmd)
			}

			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, err)
			}
		})
	}
}

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		commit   string
		date     string
		expected string
	}{
		{
			name:     "Development version",
			version:  "dev",
			commit:   "none",
			date:     "unknown",
			expected: "dev",
		},
		{
			name:     "Release version with details",
			version:  "v1.0.0",
			commit:   "abc123",
			date:     "2024-12-26",
			expected: "v1.0.0\ncommit: abc123\nbuilt at: 2024-12-26",
		},
		{
			name:     "Version containing commit",
			version:  "v1.0.0-abc123",
			commit:   "abc123",
			date:     "2024-12-26",
			expected: "v1.0.0-abc123\nbuilt at: 2024-12-26",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origVersion, origCommit, origDate := version, commit, date

			// Set test values
			version = tt.version
			commit = tt.commit
			date = tt.date

			result := formatVersion()

			// Restore original values
			version, commit, date = origVersion, origCommit, origDate

			if result != tt.expected {
				t.Errorf("formatVersion() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestVerboseFlag(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "gopeek-cmd-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file for testing
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Setup command with verbose flag and custom output path
	outputFile := filepath.Join(outputDir, "output.md")
	rootCmd.SetArgs([]string{tmpDir, "--verbose", "-o", outputFile})

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); err != nil {
		t.Errorf("Expected output file to exist: %v", err)
	}
}

func TestVersionFlag(t *testing.T) {
	// Save original values
	origVersion, origCommit, origDate := version, commit, date
	defer func() {
		version, commit, date = origVersion, origCommit, origDate
	}()

	// Set test values
	version = "v1.0.0"
	commit = "abc123"
	date = "2024-12-26"

	// Reset the command to ensure version is updated
	rootCmd.Version = formatVersion()

	// Setup command output capture
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	expectedVersion := formatVersion()
	if !strings.Contains(output, expectedVersion) {
		t.Errorf("Version mismatch:\nExpected to contain: %q\nGot: %q", expectedVersion, output)
	}
}
