# Go-Reloaded Technology Stack

## Programming Language
- **Go**: Version 1.24.9 (minimum 1.19 required)
- **Standard Library Only**: Zero external dependencies
- **Module**: `go-reloaded` (defined in go.mod)

## Build System

### Standard Build
```bash
go build -o go-reloaded cmd/go-reloaded/main.go
```

### Optimized Build (Smaller Binary)
```bash
go build -ldflags="-s -w" -o go-reloaded cmd/go-reloaded/main.go
```

### Cross-Platform Support
- Windows, Linux, macOS
- Platform-specific test scripts provided

## Development Commands

### Testing
```bash
# Run all tests
go test -count=1 ./...

# Run golden tests only
cd internal/testutils && go test -count=1 -v -run TestGoldenCases

# Run complete project tests
cd internal/testutils && go test -count=1 -v -run TestAllProject

# Run specific package tests
go test ./internal/transformer/
go test ./internal/config/
```

### Test Scripts
- **Windows**: `run_all_tests.bat`, `run_golden_tests.bat`
- **Unix/Linux/macOS**: `run_all_tests.sh`, `run_golden_tests.sh`

### Usage
```bash
./go-reloaded input.txt output.txt
```

## Dependencies

### Standard Library Packages Used
- `fmt` - Formatted I/O operations
- `os` - Operating system interface
- `io` - I/O primitives
- `strings` - String manipulation utilities
- `strconv` - String conversions
- `unicode` - Unicode character properties
- `testing` - Testing framework

### No External Dependencies
- Pure Go implementation
- No package management complexity
- Easy deployment and distribution
- Reduced security surface area

## System Requirements

### Minimum Requirements
- **RAM**: 16MB minimum, 64MB recommended
- **CPU**: Any modern processor
- **Disk**: Input file size + output file size
- **Go**: Version 1.19 or higher

### Performance Characteristics
- **Memory Usage**: Constant ~8KB regardless of file size
- **Processing**: Single-pass algorithm
- **Throughput**: Limited by disk I/O, not CPU or memory

## Development Environment

### Recommended Setup
- Go 1.24.9 or later
- IDE with Go support (VS Code, GoLand, etc.)
- Git for version control

### Code Quality Tools
- Built-in Go formatter (`go fmt`)
- Built-in Go linter (`go vet`)
- Comprehensive test suite (29 golden tests)

### Testing Philosophy
- Golden tests for end-to-end validation
- Unit tests for individual components
- No-cache testing (`-count=1`) for reliability
- Table-driven tests for multiple scenarios

## Deployment

### Binary Distribution
- Single executable file
- No runtime dependencies
- Cross-platform compilation support
- Small binary size with optimization flags

### Integration
- Command-line interface
- Standard input/output patterns
- Exit codes for error handling
- File-based I/O for pipeline integration