# Go-Reloaded Technical Architecture

## Overview

Go-Reloaded is a high-performance text processing application built in Go that transforms text files using a finite state machine (FSM) architecture with chunked processing for memory efficiency. The application processes large files with constant memory usage while applying various text transformations including numeric conversions, case modifications, punctuation corrections, and article adjustments.

## Core Architecture

### 1. Modular Component Design

The application follows a clean architecture pattern with distinct layers:

```
┌─────────────────┐
│   CLI Layer     │ ← main.go (Command-line interface)
├─────────────────┤
│ Controller Layer│ ← controller.go (Workflow orchestration)
├─────────────────┤
│ Processing Layer│ ← parser.go, transformer.go, exporter.go
├─────────────────┤
│ Config Layer    │ ← config.go (System constants)
└─────────────────┘
```

### 2. Memory-Efficient Chunked Processing

**Problem Solved**: Processing large files (GB+) without loading entire content into memory.

**Solution**: Stream processing with overlapping chunks:

```
File: [AAAAA|BBBBB|CCCCC|DDDDD]
Chunk 1: [AAAAA|BB...]  (4096 bytes + overlap)
Chunk 2: [...BB|BBBBB|CC...]  (overlap + 4096 bytes + overlap)
Chunk 3: [...CC|CCCCC|DD...]  (overlap + 4096 bytes + overlap)
```

**Key Parameters**:
- `CHUNK_BYTES = 4096`: Base chunk size for processing
- `OVERLAP_WORDS = 20`: Word overlap between chunks for context preservation

### 3. UTF-8 Boundary Alignment

**Challenge**: Chunk boundaries might split multi-byte UTF-8 characters.

**Solution**: Rune-boundary alignment in parser:

```go
// Find last complete rune boundary
for i := len(chunk) - 1; i >= 0; i-- {
    if utf8.RuneStart(chunk[i]) {
        return chunk[:i], nil
    }
}
```

This ensures chunks always end at complete character boundaries, preventing corruption.

## Component Deep Dive

### Parser Component (`internal/parser/parser.go`)

**Responsibility**: File reading with chunked processing and UTF-8 safety.

**Key Functions**:
- `ReadChunk()`: Reads next chunk with rune boundary alignment
- `HasMoreChunks()`: Checks if more data available
- `Close()`: Cleanup file resources

**Algorithm**:
1. Read raw bytes (CHUNK_BYTES + buffer for UTF-8)
2. Find last complete rune boundary
3. Split into words for overlap calculation
4. Return chunk with overlap context

### Transformer Component (`internal/transformer/transformer.go`)

**Responsibility**: Text transformation using finite state machine.

**FSM States**:
- `StateNormal`: Default processing state
- `StateInQuotes`: Inside quoted text (different punctuation rules)

**Transformation Rules**:

1. **Numeric Conversions**:
   - `(hex)` → Convert previous word from hexadecimal to decimal
   - `(bin)` → Convert previous word from binary to decimal
   - `(up)` → Convert previous word to uppercase
   - `(low)` → Convert previous word to lowercase
   - `(cap)` → Capitalize previous word

2. **Quantified Transformations**:
   - `(up, N)` → Apply uppercase to previous N words
   - `(low, N)` → Apply lowercase to previous N words
   - `(cap, N)` → Capitalize previous N words

3. **Punctuation Corrections**:
   - Remove spaces before: `. , ! ? : ;`
   - Add spaces after: `. , ! ? : ;`

4. **Quote Repositioning**:
   - Move punctuation inside quotes: `word" ,` → `word," `

5. **Article Corrections**:
   - `a` + vowel sound → `an`
   - `an` + consonant sound → `a`

**Command Chaining**: Multiple commands can be applied left-to-right:
```
"1010 (bin) (hex)" → "1010" → "10" (binary) → "16" (hex)
```

### Exporter Component (`internal/exporter/exporter.go`)

**Responsibility**: Progressive file writing with chunk coordination.

**Key Functions**:
- `WriteChunk()`: Write first chunk (creates/overwrites file)
- `AppendChunk()`: Append subsequent chunks
- `Close()`: Ensure all data is written

**Overlap Handling**: Removes overlap words from chunks 2+ to prevent duplication.

### Controller Component (`internal/controller/controller.go`)

**Responsibility**: Orchestrate the complete processing workflow.

**Processing Pipeline**:
```
Input File → Parser → Transformer → Exporter → Output File
     ↓           ↓          ↓           ↓
   Chunks    Transform   Export    Final File
```

**Algorithm**:
1. Initialize all components
2. For each chunk:
   - Parse chunk with overlap
   - Transform text using FSM
   - Export to output file
3. Cleanup resources

## Advanced Features

### 1. Cross-Chunk Context Preservation

**Problem**: Commands might reference words from previous chunks.

**Solution**: Overlap mechanism ensures context availability:
- Last 20 words from previous chunk included in next chunk
- Transformer can access previous words for commands
- Exporter removes overlap to prevent duplication

### 2. Error Resilience

**Philosophy**: Invalid commands are ignored, processing continues.

**Examples**:
- `(invalid)` → Ignored, word remains unchanged
- `(up, abc)` → Invalid number, command ignored
- Malformed syntax → Gracefully handled

### 3. State Management

**Quote State Tracking**:
- FSM tracks whether currently inside quotes
- Different punctuation rules apply inside vs outside quotes
- State persists across chunk boundaries

## Performance Characteristics

### Memory Usage
- **Constant**: O(CHUNK_BYTES + OVERLAP_WORDS) regardless of file size
- **Typical**: ~8KB per processing cycle
- **Scalability**: Can process GB+ files on minimal memory systems

### Time Complexity
- **Overall**: O(n) where n = file size in bytes
- **Per Chunk**: O(m) where m = chunk size in words
- **Transformations**: O(1) per word (hash map lookups)

### I/O Efficiency
- **Sequential Reads**: Optimized for disk/SSD access patterns
- **Buffered Writes**: Reduces system call overhead
- **Stream Processing**: No temporary file creation

## Configuration System

### System Constants (`internal/config/config.go`)

```go
const (
    CHUNK_BYTES   = 4096  // Base chunk size for memory efficiency
    OVERLAP_WORDS = 20    // Context preservation between chunks
)
```

**Tuning Guidelines**:
- **CHUNK_BYTES**: Larger = fewer I/O operations, more memory usage
- **OVERLAP_WORDS**: Larger = better context, more processing overhead

### Validation
- Constants validated at startup
- Ensures CHUNK_BYTES > 0 and OVERLAP_WORDS ≥ 0
- Prevents runtime configuration errors

## Error Handling Strategy

### Graceful Degradation
1. **File Errors**: Clear error messages, early termination
2. **UTF-8 Errors**: Skip invalid sequences, continue processing
3. **Command Errors**: Ignore invalid commands, preserve original text
4. **Memory Errors**: Fail fast with descriptive messages

### Error Propagation
- Errors bubble up through component layers
- Context preserved in error messages
- Clean resource cleanup on failures

## Testing Architecture

### Test-Driven Development
- **Golden Tests**: 22 comprehensive test cases (T1-T22)
- **Unit Tests**: Each component thoroughly tested
- **Integration Tests**: End-to-end workflow validation

### Test Categories
1. **Transformation Tests**: Verify all FSM rules
2. **Chunking Tests**: Validate overlap and boundary handling
3. **UTF-8 Tests**: Ensure character integrity
4. **Error Tests**: Confirm graceful error handling
5. **Performance Tests**: Memory and speed validation

## Deployment Considerations

### System Requirements
- **Go Version**: 1.19+ (for UTF-8 improvements)
- **Memory**: Minimum 16MB available
- **Disk**: Space for input + output files
- **OS**: Cross-platform (Linux, macOS, Windows)

### Build Optimization
```bash
go build -ldflags="-s -w" -o go-reloaded cmd/go-reloaded/main.go
```
- `-s -w`: Strip debug symbols for smaller binary
- Static linking: No external dependencies

### Performance Tuning
- **Large Files**: Consider increasing CHUNK_BYTES
- **Memory Constrained**: Decrease OVERLAP_WORDS
- **I/O Bound**: Use SSD storage for better performance

## Future Enhancements

### Potential Improvements
1. **Parallel Processing**: Multi-threaded chunk processing
2. **Compression**: On-the-fly compression for large outputs
3. **Streaming**: Real-time processing for continuous input
4. **Plugin System**: Custom transformation rules
5. **Configuration Files**: Runtime parameter adjustment

### Scalability Considerations
- **Horizontal**: Multiple instances for different files
- **Vertical**: Larger chunks for high-memory systems
- **Distributed**: Network-based processing for massive files

## Conclusion

Go-Reloaded demonstrates sophisticated text processing with:
- **Memory Efficiency**: Constant memory usage regardless of file size
- **Robustness**: Graceful error handling and UTF-8 safety
- **Performance**: Linear time complexity with optimized I/O
- **Maintainability**: Clean architecture with comprehensive testing

The FSM-based approach with chunked processing provides a scalable foundation for text transformation tasks while maintaining code clarity and system reliability.