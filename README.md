# GoPeek

[![Go Report Card](https://goreportcard.com/badge/github.com/nouuu/gopeek)](https://goreportcard.com/report/github.com/nouuu/gopeek)
[![GoDoc](https://godoc.org/github.com/nouuu/gopeek?status.svg)](https://godoc.org/github.com/nouuu/gopeek)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

GoPeek is a lightweight command-line tool that generates comprehensive documentation of project structures. It recursively scans directories to create a navigable Markdown document containing both the project tree and file contents, making it ideal for project exploration and documentation.

## Features

- 🌳 Recursive directory scanning with an intuitive tree structure
- 📝 Automatic Markdown generation with file contents
- 🔍 Smart binary file detection
- ⚡ Efficient large file handling with size limits
- 🎯 Configurable ignore patterns (supports .gitignore)
- 🔗 Generated anchors for easy navigation

## Installation

```bash
go install github.com/yourusername/gopeek@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/gopeek.git
cd gopeek
make build
```

## Usage

Basic usage:
```bash
gopeek [path] [flags]
```

Available flags:
```bash
Flags:
  -o, --output string        Output file path (default "project_knowledge.md")
  -i, --ignore stringSlice   Patterns to ignore
```

Example:
```bash
# Scan current directory
gopeek .

# Scan specific directory with custom output
gopeek /path/to/project -o documentation.md

# Scan with custom ignore patterns
gopeek . -i "*.log" -i "build/*"
```

## Output Format

GoPeek generates a structured Markdown document with two main sections:

1. Project Structure: A tree view of your project with clickable links
2. File Contents: The content of each file with syntax highlighting

Example output:
```markdown
# Project Structure
- 📁 project
  - 📄 [main.go](#main-go)
  - 📁 internal
    - 📄 [types.go](#internal-types-go)

# Files Content
# 📄 main.go
```go
package main
// ... file content
```
```

## Development

### Prerequisites

- Go 1.22 or higher
- Make (for build automation)

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] HTML output format
- [ ] Project analytics (LOC, file types, etc.)
- [ ] Custom templating support
- [ ] VS Code extension
- [ ] Symlink handling
- [ ] Code syntax highlighting improvements

## Acknowledgments

- Inspired by the need for better project documentation tools
- Built with [Cobra](https://github.com/spf13/cobra)
