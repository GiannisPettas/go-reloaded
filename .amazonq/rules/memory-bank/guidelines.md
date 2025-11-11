# Go-Reloaded Development Guidelines

## Code Quality Standards

### Error Handling Patterns (100% of Go files)
- **Consistent Error Wrapping**: Always use `fmt.Errorf` with `%w` verb for error chaining
- **Context-Rich Messages**: Include operation details in error messages
```go
return fmt.Errorf("failed to read chunk at offset %d: %w", offset, err)
return fmt.Errorf("failed to open file %s: %w", filepath, err)
```
- **Early Return Pattern**: Check errors immediately, avoid nested error handling
- **Preserve Error Chain**: Never lose original error information

### Naming Conventions (Observed across codebase)
- **Package Names**: Single word, lowercase (`parser`, `transformer`, `controller`)
- **Function Names**: PascalCase for exported, camelCase for internal
- **Constants**: ALL_CAPS with underscores (`STATE_TEXT`, `CHUNK_BYTES`)
- **Variables**: Descriptive camelCase (`lastWordIdx`, `overlapContext`)

### Code Structure Patterns
- **Single Responsibility**: Each function has one clear purpose
- **Fixed-Size Buffers**: Use arrays for memory control (`[50]Token`)
- **State Machines**: Explicit state constants and switch statements
- **Builder Pattern**: Use `strings.Builder` for efficient string construction

## Architectural Patterns

### Finite State Machine Implementation (transformer.go)
```go
// State constants
const (
    STATE_TEXT = iota
    STATE_COMMAND
)

// State switching with explicit conditions
switch state {
case STATE_TEXT:
    // Handle text parsing
case STATE_COMMAND:
    // Handle command processing
}
```

### Memory Management Patterns
- **Fixed Buffers**: `[80]Token` array prevents unbounded growth
- **Chunked Processing**: Process large files in fixed-size chunks
- **Overlap Handling**: Maintain context between chunks without memory growth
```go
buffer := make([]byte, config.CHUNK_BYTES)  // Fixed size allocation
```

### Component Orchestration (controller.go)
- **Single Entry Point**: Controller orchestrates all components
- **Error Propagation**: Wrap errors with context at each layer
- **File Size Adaptation**: Different strategies for small vs large files
```go
if fileInfo.Size() <= int64(config.CHUNK_BYTES) {
    return processSingleChunk(inputPath, outputPath)
}
return processChunkedFile(inputPath, outputPath)
```

## Implementation Standards

### String Processing Algorithms
- **Rune-Safe Operations**: Always use `[]rune(text)` for Unicode safety
- **Builder Usage**: Use `strings.Builder` for efficient concatenation
- **Field Splitting**: Use `strings.Fields()` for word extraction
```go
runes := []rune(text)  // Unicode-safe character access
words := strings.Fields(line)  // Whitespace-aware splitting
```

### Token Processing System
- **Token Types**: Enum-style constants for different token categories
- **Conveyor Belt**: Fixed-size token buffer with overflow handling
- **Real-Time Processing**: Process tokens as they arrive, not batch
```go
type Token struct {
    Type  int     // WORD, COMMAND, PUNCTUATION, etc.
    Value string  // Actual content
}

// TokenProcessor with 80-token belt
type TokenProcessor struct {
    tokens [80]Token  // Fixed-size conveyor belt
    // ... other fields
}
```

### Configuration Management
- **Centralized Constants**: All system limits in config package
- **Validation**: Validate configuration at startup
- **Immutable Values**: Constants prevent accidental modification
```go
const (
    CHUNK_BYTES = 4096
    OVERLAP_WORDS = 20  // Range: 10-20 words for chunk overlap
)
```

## Testing Patterns

### Test Organization (Observed in test files)
- **Golden Tests**: End-to-end validation with expected outputs
- **Unit Tests**: Component-specific functionality testing
- **Table-Driven Tests**: Multiple test cases in structured format
- **No-Cache Testing**: Use `-count=1` flag for reliable results

### Error Testing Standards
- **Both Paths**: Test success and failure scenarios
- **Error Content**: Verify error message content and wrapped errors
- **Edge Cases**: Test boundary conditions and invalid inputs

## Performance Optimization

### Memory Efficiency Patterns
- **Constant Memory**: Fixed-size data structures regardless of input size
- **Single-Pass Processing**: Never re-read or re-process data
- **Streaming**: Process and output data continuously
- **UTF-8 Boundary Handling**: Prevent character corruption in chunks

### Algorithm Design
- **Dual FSM Architecture**: Low-level parsing + high-level processing
- **Lookahead Validation**: Check command validity before state changes
- **Overlap Management**: Maintain transformation context across chunks

## Code Documentation

### Comment Standards
- **Algorithm Explanation**: Document complex logic and state machines
- **Function Purpose**: Clear description of what each function does
- **Edge Cases**: Document special handling and boundary conditions
- **Performance Notes**: Explain memory and efficiency considerations

### Package Documentation
- **Clear Purpose**: Each package has well-defined responsibility
- **Usage Examples**: Show how components work together
- **Error Handling**: Document error conditions and recovery

## Development Workflow

### Build and Test Standards
- **Zero Dependencies**: Use only Go standard library
- **Cross-Platform**: Support Windows, Linux, macOS
- **Optimized Builds**: Provide size-optimized compilation flags
- **Comprehensive Testing**: 29 golden test cases for validation

### Code Review Checklist
- Error handling follows wrapping patterns
- Memory usage is bounded and predictable
- Unicode safety is maintained throughout
- State machines have clear transitions
- Tests cover both success and failure paths