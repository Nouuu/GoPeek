package scanner

import (
	"os"
	"testing"
)

func TestIsBinaryFile(t *testing.T) {
	// Create temp files
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "Empty file",
			content:  []byte{},
			expected: false,
		},
		{
			name:     "Simple text",
			content:  []byte("Hello, World!\n"),
			expected: false,
		},
		{
			name:     "Text with emoji",
			content:  []byte("Hello! 👋 World 🌍"),
			expected: false,
		},
		{
			name:     "Invalid UTF-8 sequence",
			content:  []byte{0xFF, 0xFE, 0xFD},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpfile, err := os.CreateTemp("", "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write content
			if _, err := tmpfile.Write(tt.content); err != nil {
				t.Fatal(err)
			}
			tmpfile.Close()

			// Test the file
			result, err := isBinaryFile(tmpfile.Name())
			if err != nil {
				t.Fatal(err)
			}
			if result != tt.expected {
				t.Errorf("isBinaryFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}
