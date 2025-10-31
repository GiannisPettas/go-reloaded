# Go-Reloaded Development Guidelines

## Code Quality Standards

### Naming Conventions
- **Constants**: ALL_CAPS with underscores (e.g., `CHUNK_BYTES`, `OVERLAP_WORDS`, `STATE_TEXT`)
- **Types**: PascalCase (e.g., `Token`, `TokenProcessor`, `StreamProcessor`)
- **Functions**: PascalCase for exported, camelCase for private (e.g., `ProcessText`, `fixArticles`)
- **Variables**: camelCase (e.g., `wordBuilder`, `cmdBuilder`, `lastWordIdx`)
- **Packages**: lowercase single words (e.g., `parser`, `transformer`, `controller`)

### Error Handling Patterns
- Always use `fmt.Errorf` with context and error wrapping using `%w` verb
- Check errors immediately after operations that can fail
- Provide meaningful error messages with operation context
```go
if err != nil {
    return fmt.Errorf("failed to read chunk at offset %d: %w", offset, err)
}
```

### Function Structure
- Keep functions focused on single responsibility
- Use early returns to reduce nesting
- Validate inputs at function start
- Handle edge cases explicitly (empty strings, nil values)

## Architectural Patterns

### State Machine Implementation
- Use integer constants for states (iota pattern)
- Maintain clear state transitions in switch statements
- Separate low-level and high-level FSM concerns
- Buffer management with fixed-size arrays for memory efficiency

### Memory Management
- Use fixed-size buffers to ensure constant memory usage
- Implement buffer flushing when capacity is reached
- Reset builders after use to prevent memory leaks
- Prefer `strings.Builder` over string concatenation

### File Processing Patterns
- Implement chunked processing for large files
- Handle UTF-8 boundaries properly with `AdjustToRuneBoundary`
- Use overlap context to maintain processing continuity
- Separate single-chunk and multi-chunk processing logic

## Internal API Usage

### Transformer Package
```go
// Single-pass processing with dual FSM
result := transformer.ProcessText(text)

// Stream processing for large files
processor := transformer.NewStreamProcessor()
output := processor.ProcessChunk(data)
finalOutput := processor.Flush()
```

### Parser Package
```go
// Read chunk with UTF-8 safety
data, err := parser.ReadChunk(filepath, offset)

// Handle word overlap between chunks
overlap, remaining := parser.ExtractOverlapWords(text)
merged := parser.PrependOverlapWords(overlap, newChunk)
```

### Controller Package
```go
// Main processing entry point
err := controller.ProcessFile(inputPath, outputPath)
```

## Code Idioms and Patterns

### String Processing
- Use `strings.Fields` for word tokenization
- Use `strings.Builder` for efficient string construction
- Reset builders after use: `builder.Reset()`
- Check string suffixes/prefixes before operations

### Buffer Management
```go
// Fixed-size token buffer with overflow handling
tokens [50]Token
if tp.tokenIdx < len(tp.tokens) {
    tp.tokens[tp.tokenIdx] = token
    tp.tokenIdx++
} else {
    // Flush and shift logic
}
```

### Rune Processing
- Always work with `[]rune` for Unicode safety
- Use `strings.ToLower()` for case-insensitive comparisons
- Handle punctuation as separate tokens

### File Operations
- Always defer file closure: `defer file.Close()`
- Use `os.Stat` to check file existence and get size
- Handle `io.EOF` as expected condition, not error

## Testing Standards

### Test Organization
- Place tests in same package with `_test.go` suffix
- Use table-driven tests for multiple scenarios
- Implement golden tests for comprehensive validation
- Test edge cases (empty input, large files, UTF-8 boundaries)

### Test Utilities
- Centralize test utilities in `internal/testutils`
- Use golden files for expected output validation
- Test both single-chunk and multi-chunk processing paths

## Documentation Standards

### Function Comments
- Document exported functions with purpose and behavior
- Include parameter descriptions for complex functions
- Document return values and error conditions
- Use complete sentences starting with function name

### Package Comments
- Each package should have a package-level comment
- Describe the package's primary responsibility
- Include usage examples for complex packages

## Performance Considerations

### Memory Efficiency
- Use fixed-size buffers to prevent unbounded growth
- Implement proper buffer flushing strategies
- Avoid string concatenation in loops
- Reset reusable objects after use

### Processing Efficiency
- Single-pass algorithms preferred over multi-pass
- Minimize memory allocations in hot paths
- Use appropriate data structures (arrays vs slices)
- Implement streaming for large data processing

## Configuration Management

### Constants Organization
- Centralize configuration in `internal/config` package
- Use meaningful constant names with units where applicable
- Group related constants together
- Document the purpose of each configuration value