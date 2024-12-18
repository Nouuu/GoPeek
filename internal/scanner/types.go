package scanner

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Output struct {
	structure []string
	contents  []string
}

const maxFileSize = 10 * 1024 * 1024 // 10MB

func (o *Output) AddStructure(path string, info fs.FileInfo, depth int) {
	indent := strings.Repeat("  ", depth)
	if info.IsDir() {
		o.structure = append(o.structure, fmt.Sprintf("%s- 📁 %s", indent, info.Name()))
	} else {
		anchor := createAnchor(path)
		o.structure = append(o.structure, fmt.Sprintf("%s- 📄 [%s](#%s)", indent, info.Name(), anchor))
	}
}

func (o *Output) AddContent(path string, relPath string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting file stats: %w", err)
	}

	if info.Size() > maxFileSize {
		return fmt.Errorf("file too large (max %dMB): %s", maxFileSize/1024/1024, path)
	}

	isBinary, err := isBinaryFile(path)
	if err != nil {
		return fmt.Errorf("error checking if file is binary: %w", err)
	}

	anchor := createAnchor(path)

	if isBinary {
		o.contents = append(o.contents, fmt.Sprintf("\n<a id=\"%s\"></a>\n# 📄 %s\n```\n[binary file]\n```\n", anchor, relPath))
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	ext := filepath.Ext(path)
	if ext != "" {
		ext = ext[1:]
	}

	fileContent := fmt.Sprintf("\n<a id=\"%s\"></a>\n# 📄 %s\n```%s\n%s\n```\n",
		anchor, relPath, ext, string(content))
	o.contents = append(o.contents, fileContent)
	return nil
}

func (o *Output) Generate() string {
	return fmt.Sprintf("# Project Structure\n\n%s\n\n# Files Content\n%s",
		strings.Join(o.structure, "\n"),
		strings.Join(o.contents, "\n"))
}

func isBinaryFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read first 512 bytes
	buff := make([]byte, 512)
	n, err := file.Read(buff)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Check if content contains non-printable characters
	isbin := !utf8.Valid(buff[:n])

	if isbin {
		return true, nil
	}
	return false, nil
}

func createAnchor(path string) string {
	anchor := strings.Map(func(r rune) rune {
		if r == os.PathSeparator {
			return '-'
		}
		return r
	}, path)
	anchor = strings.ReplaceAll(anchor, ".", "-")
	anchor = strings.ToLower(strings.ReplaceAll(anchor, " ", "-"))

	// Supprime tous les caractères non alphanumériques (sauf les tirets)
	anchor = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, anchor)

	// Convertit en minuscules
	anchor = strings.ToLower(anchor)

	// Remplace les tirets multiples par un seul tiret
	for strings.Contains(anchor, "--") {
		anchor = strings.ReplaceAll(anchor, "--", "-")
	}

	// Supprime les tirets au début et à la fin
	anchor = strings.Trim(anchor, "-")

	// Si l'ancre est vide après tout ça, utilise un fallback
	if anchor == "" {
		return "file"
	}

	return anchor
}
