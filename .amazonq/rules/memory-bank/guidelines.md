# Go-Reloaded Development Guidelines

## Code Quality Standards

### Package Structure
- **Single Responsibility**: Each package has one focused purpose (parser, transformer, exporter, controller)
- **Internal Organization**: All core logic in `internal/` to prevent external imports
- **Test Co-location**: Test files alongside source files with `_test.go` suffix
- **Utility Separation**: Testing utilities in dedicated `testutils/` package

### Naming Conventions
- **Constants**: ALL_CAPS with underscores (e.g., `CHUNK_BYTES`, `OVERLAP_WORDS`, `STATE_TEXT`)
- **Types**: PascalCase for exported types (`Token`, `TokenProcessor`, `GoldenTest`)
- **Functions**: PascalCase for exported, camelCase for internal (`ProcessText`, `addToken`)
- **Variables**: camelCase throughout (`tokenIdx`, `wordBuilder`, `cmdBuilder`)

### Error Handling Patterns
- **Wrapped Errors**: Use `fmt.Errorf` with `%w` verb for error chaining
- **Context Preservation**: Include operation context in error messages
- **Early Returns**: Check errors immediately and return early
```go
if err != nil {
    return fmt.Errorf("failed to read chunk at offset %d: %w", offset, err)
}
```

### Memory Management
- **Fixed Buffers**: Use arrays instead of slices for predictable memory usage
- **String Builders**: Use `strings.Builder` for efficient string concatenation
- **Buffer Reuse**: Reset and reuse builders instead of creating new ones
- **Chunked Processing**: Process large files in fixed-size chunks to maintain constant memory

## Architectural Patterns

### Finite State Machine Implementation
- **State Constants**: Define states as `iota` constants
- **State Switching**: Use switch statements for state transitions
- **Dual FSM**: Separate low-level (parsing) and high-level (token processing) state machines
```go
const (
    STATE_TEXT = iota
    STATE_COMMAND
)
```

### Token-Based Processing
- **Token Types**: Define token types as constants (`WORD`, `PUNCTUATION`, `SPACE`, `NEWLINE`)
- **Token Struct**: Simple struct with `Type` and `Value` fields
- **Buffer Management**: Fixed-size token arrays with overflow handling

### Configuration Management
- **Centralized Config**: All constants in `config/` package
- **Import Pattern**: Import config package in all modules that need constants
- **Naming**: Descriptive constant names (`CHUNK_BYTES`, `OVERLAP_WORDS`)

## Testing Standards

### Test Organization
- **Golden Tests**: Use markdown files for comprehensive test cases
- **Test Utilities**: Centralized test helpers in `testutils/` package
- **File Management**: Create and cleanup test files using utilities
```go
filepath, err := testutils.CreateTestFile(content)
if err != nil {
    t.Fatalf("Failed to create test file: %v", err)
}
defer testutils.CleanupTestFile(filepath)
```

### Test Patterns
- **Table-Driven Tests**: Use structs for test cases with multiple scenarios
- **Error Testing**: Verify both success and failure cases
- **Boundary Testing**: Test edge cases (empty files, exact chunk sizes, UTF-8 boundaries)
- **Integration Testing**: Test complete workflows through controller

### Test Naming
- **Descriptive Names**: `TestReadChunkExactSize`, `TestAdjustToRuneBoundaryIncomplete`
- **Scenario Coverage**: Test normal, edge, and error cases
- **Helper Functions**: Use helper functions for common validation patterns

## Performance Optimization

### Memory Efficiency
- **Constant Memory**: Maintain ~8KB memory usage regardless of file size
- **Single Pass**: Process text in one pass through dual FSM
- **Buffer Recycling**: Reuse buffers and builders instead of allocating new ones
- **UTF-8 Safety**: Use `AdjustToRuneBoundary` to prevent character corruption

### Processing Patterns
- **Chunked Reading**: Read files in `CHUNK_BYTES` sized chunks
- **Overlap Handling**: Maintain word context between chunks using `OVERLAP_WORDS`
- **State Preservation**: Maintain FSM state across chunk boundaries
- **Efficient String Operations**: Use `strings.Builder` for concatenation

## Code Documentation

### Comment Standards
- **Package Comments**: Brief description of package purpose
- **Function Comments**: Describe what function does, not how
- **Complex Logic**: Explain non-obvious algorithms and state transitions
- **Constants**: Document purpose and usage of configuration constants

### API Design
- **Exported Functions**: Clear, descriptive names with proper error handling
- **Parameter Validation**: Check inputs and return meaningful errors
- **Return Values**: Consistent error handling patterns
- **Interface Simplicity**: Minimal, focused function signatures

## Development Workflow

### File Processing Pipeline
1. **Controller**: Orchestrates the workflow
2. **Parser**: Reads and chunks input files
3. **Transformer**: Applies FSM-based transformations
4. **Exporter**: Writes processed output

### Error Propagation
- **Consistent Wrapping**: Use `fmt.Errorf` with `%w` for error chains
- **Context Information**: Include relevant details (file paths, offsets)
- **Early Termination**: Return errors immediately, don't continue processing

### Testing Integration
- **Comprehensive Coverage**: Test all transformation scenarios
- **Golden Test Suite**: 27 test cases covering all features
- **Automated Validation**: Use `go test -count=1 ./...` for full test runs