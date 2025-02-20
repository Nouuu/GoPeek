# GoPeek

[![Go Report Card](https://goreportcard.com/badge/github.com/nouuu/gopeek)](https://goreportcard.com/report/github.com/nouuu/gopeek)
[![Go Reference](https://pkg.go.dev/badge/github.com/nouuu/gopeek.svg)](https://pkg.go.dev/github.com/nouuu/gopeek)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nouuu/gopeek)](https://golang.org/doc/devel/release.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Build Status](https://github.com/nouuu/gopeek/workflows/build/badge.svg)](https://github.com/nouuu/gopeek/actions?query=workflow%3Abuild)
[![Tests](https://github.com/nouuu/gopeek/workflows/tests/badge.svg)](https://github.com/nouuu/gopeek/actions?query=workflow%3Atests)
[![Lint](https://github.com/nouuu/gopeek/workflows/lint/badge.svg)](https://github.com/nouuu/gopeek/actions?query=workflow%3Alint)
[![Security](https://github.com/nouuu/gopeek/workflows/security/badge.svg)](https://github.com/nouuu/gopeek/actions?query=workflow%3Asecurity)


[![Release](https://img.shields.io/github/v/release/nouuu/gopeek)](https://github.com/nouuu/gopeek/releases)
[![Issues](https://img.shields.io/github/issues/nouuu/gopeek)](https://github.com/nouuu/gopeek/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/nouuu/gopeek)](https://github.com/nouuu/gopeek/pulls)
[![Contributors](https://img.shields.io/github/contributors/nouuu/gopeek)](https://github.com/nouuu/gopeek/graphs/contributors)
[![Lines of Code](https://tokei.rs/b1/github/nouuu/gopeek)](https://github.com/nouuu/gopeek)
[![Last Commit](https://img.shields.io/github/last-commit/nouuu/gopeek)](https://github.com/nouuu/gopeek/commits/main)

GoPeek is a lightweight command-line tool that generates comprehensive documentation of project structures. It recursively scans directories to create a navigable Markdown document containing both the project tree and file contents, making it ideal for project exploration, documentation and AI context providing.
## Features

- 🌳 Recursive directory scanning with an intuitive tree structure
- 📝 Automatic Markdown generation with file contents
- 🔍 Smart binary file detection
- ⚡ Efficient large file handling with size limits
- 🎯 Configurable ignore patterns (supports .gitignore)
- 🔗 Generated anchors for easy navigation

## Installation

```bash
go install github.com/nouuu/gopeek/cmd/gopeek@latest
```

Or build from source:

```bash
git clone https://github.com/nouuu/gopeek.git
cd gopeek
make build
```

You can also install it from source :

```bash
git clone https://github.com/nouuu/gopeek.git
cd gopeek
make install
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
  -v, --version              Show version
  --verbose                  Enable verbose output
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

## Roadmap 🗺️

### Core Features ✨

- [x] Basic directory scanning
- [x] Markdown output generation
- [x] Binary file detection
- [x] .gitignore support
- [x] CLI interface with cobra
- [x] File size limits
- [x] Github Actions CI/CD

### Next Steps 🚀

- [ ] Advanced Error Handling 🛡️
    - [ ] Custom error types
    - [ ] Error context and wrapping
    - [ ] Operation summaries
- [ ] Extended Output Options 📝
    - [ ] HTML with navigation
    - [ ] JSON output
    - [ ] Template customization
- [ ] Performance Features ⚡
    - [ ] Parallel file scanning
    - [ ] Memory usage optimization
    - [ ] Progress indicators

## Acknowledgments

- Inspired by the need for better project documentation tools
- Built with [Cobra](https://github.com/spf13/cobra)
