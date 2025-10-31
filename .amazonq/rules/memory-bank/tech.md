# Go-Reloaded Technology Stack

## Programming Language
- **Go 1.25.1**: Primary development language
- **Standard Library Only**: No external dependencies required

## Build System
- **Go Modules**: Dependency management with go.mod
- **Native Go Build**: Standard `go build` command
- **Optimized Builds**: Support for `-ldflags="-s -w"` for smaller binaries

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
# Run all tests
go test -count=1 ./...

# Run golden test suite
cd internal/testutils && go test -v

# Run specific package tests
go test ./internal/transformer/
```

### Usage
```bash
./go-reloaded input.txt output.txt
```

## Architecture Technologies

### Finite State Machines
- **Dual FSM Design**: Two state machines working in parallel
- **Token-based Processing**: Efficient text parsing and transformation
- **State Management**: Clean state transitions for complex transformations

### Memory Management
- **Fixed Buffer Sizes**: Constant memory usage regardless of file size
- **Chunked Processing**: Smart overlap handling for large files
- **Stream Processing**: Single-pass data processing

### File I/O
- **Buffered Reading**: Efficient file reading with configurable buffer sizes
- **UTF-8 Support**: Proper handling of international characters
- **Error Handling**: Comprehensive error management throughout the pipeline

## System Requirements
- **RAM**: 16MB minimum, 64MB recommended
- **CPU**: Any modern processor
- **Disk**: Input file size + output file size
- **OS**: Cross-platform (Windows, Linux, macOS)

## Performance Characteristics
- **Memory Usage**: Constant ~8KB regardless of file size
- **Processing Speed**: Single-pass, linear time complexity
- **File Size Support**: No practical limits (tested with 1GB+ files)
- **Concurrency**: Single-threaded design for predictable behavior