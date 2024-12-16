package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Output struct {
	structure []string
	contents  []string
}

func (o *Output) AddStructure(path string, info fs.FileInfo, depth int) {
	indent := strings.Repeat("  ", depth)
	if info.IsDir() {
		o.structure = append(o.structure, fmt.Sprintf("%s- 📁 %s", indent, info.Name()))
	} else {
		o.structure = append(o.structure, fmt.Sprintf("%s- 📄 %s", indent, info.Name()))
	}
}

func (o *Output) AddContent(path string, relPath string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	ext := filepath.Ext(path)
	if ext != "" {
		ext = ext[1:]
	}

	fileContent := fmt.Sprintf("\n# 📄 %s\n```%s\n%s\n```\n",
		relPath, ext, string(content))
	o.contents = append(o.contents, fileContent)
	return nil
}

func (o *Output) Generate() string {
	return fmt.Sprintf("# Project Structure\n\n%s\n\n# Files Content\n%s",
		strings.Join(o.structure, "\n"),
		strings.Join(o.contents, "\n"))
}
