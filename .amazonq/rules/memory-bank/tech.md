# Go-Reloaded Technology Stack

## Programming Language
- **Go**: Version 1.24.9 (minimum 1.19 required)
- **Module**: `go-reloaded` (no external dependencies)

## Build System
- **Native Go Build**: Uses standard `go build` command
- **Module System**: Go modules for dependency management
- **Zero Dependencies**: Pure Go standard library implementation

## Development Commands

### Building
```bash
# Standard build
go build -o go-reloaded cmd/go-reloaded/main.go

# Optimized build (smaller binary)
go build -ldflags="-s -w" -o go-reloaded cmd/go-reloaded/main.go
```

### Testing
```bash
# Recommended: Run all tests with formatting
cd internal/testutils && go test -v -run TestAllProject

# Run all tests manually
go test -count=1 ./...

# Run golden test suite only
cd internal/testutils && go test -v -run TestGoldenCases

# Run specific package tests
go test ./internal/transformer/
go test ./internal/config/
```

### Running
```bash
./go-reloaded input.txt output.txt
```

## Architecture Technologies

### Core Patterns
- **Finite State Machines (FSM)**: Dual FSM architecture for text processing
- **Chunked Processing**: Memory-efficient file handling
- **Single-Pass Processing**: No multiple iterations over data

### Standard Library Usage
- **`os`**: File operations and command-line arguments
- **`bufio`**: Buffered I/O for efficient file reading/writing
- **`strings`**: String manipulation and processing
- **`strconv`**: String to numeric conversions
- **`unicode`**: Character classification and case conversion
- **`testing`**: Comprehensive test framework

## Performance Characteristics
- **Memory Usage**: Constant ~8KB regardless of file size
- **Processing**: Single-pass algorithm
- **File Support**: Handles files from KB to GB+ sizes
- **UTF-8**: Full Unicode support without corruption

## System Requirements
- **RAM**: 16MB minimum, 64MB recommended
- **CPU**: Any modern processor
- **Disk**: Input file size + output file size
- **OS**: Cross-platform (Windows, Linux, macOS)

## Development Environment
- **Go Toolchain**: Standard Go development tools
- **Testing**: Built-in Go testing framework
- **Documentation**: Godoc-compatible comments
- **Version Control**: Git-based workflow

## Quality Assurance
- **Test Coverage**: 27 comprehensive golden test cases
- **Integration Tests**: End-to-end testing scenarios
- **Unit Tests**: Individual component testing
- **Performance Tests**: Memory and speed validation