# Go-Reloaded Controller: Workflow Orchestration Guide

## What Does the Controller Do?

The controller is the **orchestrator** of Go-Reloaded. It coordinates the entire processing pipeline:

```
Input File → Parser → Transformer → Exporter → Output File
```

It handles:
- **File size detection**: Single chunk vs. chunked processing
- **Memory management**: Configurable constant memory usage (1KB-8KB) regardless of file size
- **Overlap handling**: Maintains context between chunks for large files (10-20 words, configurable)
- **Error propagation**: Wraps errors with context throughout the pipeline
- **Workflow coordination**: Ensures components work together seamlessly

## Core Architecture

### Single Entry Point

```go
func ProcessFile(inputPath, outputPath string) error
```

**This is the ONLY public function** - clean, simple API that hides all complexity.

### Two Processing Strategies

The controller automatically chooses the optimal strategy based on file size:

#### 1. Single Chunk Processing (Small Files ≤ CHUNK_BYTES)
```go
if fileInfo.Size() <= int64(config.CHUNK_BYTES) {
    return processSingleChunk(inputPath, outputPath)
}
```

**Default**: Files ≤ 4KB, but **configurable** from 1KB to 8KB via `config.CHUNK_BYTES`

#### 2. Chunked Processing (Large Files > CHUNK_BYTES)
```go
return processChunkedFile(inputPath, outputPath)
```

**Default**: Files > 4KB, but **configurable** based on `config.CHUNK_BYTES` setting

**Why this split?** 
- Small files: Maximum performance, no overhead
- Large files: Constant memory usage, no size limits

### Why Two Separate Functions?

**The functions handle completely different complexity levels:**

**`processSingleChunk` (Simple):**
```go
// Simple: Read → Transform → Write
data := parser.ReadChunk(inputPath, 0)
result := transformer.ProcessText(string(data))
exporter.WriteChunk(outputPath, result)
```

**`processChunkedFile` (Complex):**
```go
// Complex: Loop with overlap management
for {
    data := parser.ReadChunk(inputPath, offset)
    
    // Merge with previous overlap
    textToProcess = overlapContext + " " + chunkText
    
    // Process and handle overlap removal
    processedChunk := transformer.ProcessText(textToProcess)
    
    // Extract overlap for next iteration
    newOverlap, remaining := parser.ExtractOverlapWords(processedChunk)
    
    // Write/append logic, update offset, check boundaries, etc.
}
```

**Key differences:**
- **Overlap handling**: Large files need word context between chunks
- **Write strategy**: Single write vs append operations  
- **Loop complexity**: Simple call vs complex iteration with state management
- **Memory management**: Different patterns for small vs large files

**Benefits of separation:**
- **Performance**: Small files avoid unnecessary overhead
- **Simplicity**: Easy debugging for common case (most files are small)
- **Maintainability**: Clear separation of concerns
- **Code clarity**: Each function has a single, focused purpose

You could combine them, but it would make the code more complex for the simple case.

## Single Chunk Processing - `processSingleChunk()`

**For files ≤ CHUNK_BYTES - Simple and fast:**

```go
func processSingleChunk(inputPath, outputPath string) error {
    // 1. Read entire file at once
    data, err := parser.ReadChunk(inputPath, 0)
    
    // 2. Transform in single pass
    result := transformer.ProcessText(string(data))
    
    // 3. Write complete result
    err = exporter.WriteChunk(outputPath, result)
    
    return nil
}
```

**Benefits:**
- **Fastest possible processing** - no chunking overhead
- **Simplest logic** - straight pipeline
- **Perfect for most use cases** - most text files are small

## Chunked Processing - `processChunkedFile()`

**For files > CHUNK_BYTES - Constant memory usage:**

### The Challenge: Context Preservation

**Problem:** Commands can reference words from previous chunks:
```
Chunk 1: "these three words should"
Chunk 2: "be uppercase (up, 4)"
```

The `(up, 4)` command needs to access "these three words should" from the previous chunk!

### The Solution: Overlap Context

```go
var overlapContext string  // Maintains context between chunks

for each chunk {
    // 1. Merge with previous context
    textToProcess = overlapContext + " " + chunkText
    
    // 2. Process merged text
    processedChunk = transformer.ProcessText(textToProcess)
    
    // 3. Remove overlap from result (avoid duplication)
    // 4. Extract new overlap for next chunk
    newOverlap, remaining = parser.ExtractOverlapWords(processedChunk)
    
    // 5. Write remaining text
    // 6. Update context for next iteration
    overlapContext = newOverlap
}
```

### Step-by-Step Chunked Processing

#### Step 1: Read Chunk with Offset
```go
data, err := parser.ReadChunk(inputPath, offset)
```

#### Step 2: Merge with Overlap Context
```go
var textToProcess string
if overlapContext != "" {
    textToProcess = overlapContext + " " + chunkText
} else {
    textToProcess = chunkText
}
```

**Example:**
```
Chunk 1: "word1 word2 word3 word4"
Overlap: "word3 word4" (last 2 words)

Chunk 2: "word5 word6 word7 word8"
Merged:  "word3 word4 word5 word6 word7 word8"
```

#### Step 3: Transform Merged Text
```go
processedChunk := transformer.ProcessText(textToProcess)
```

**The transformer sees the complete context** and can apply commands correctly.

#### Step 4: Remove Overlap Duplication
```go
if overlapContext != "" {
    overlapWordCount := len(strings.Fields(overlapContext))
    processedWords := strings.Fields(processedChunk)
    if len(processedWords) > overlapWordCount {
        processedChunk = strings.Join(processedWords[overlapWordCount:], " ")
    }
}
```

**Why?** We don't want to write the same words twice to the output file.

#### Step 5: Extract New Overlap
```go
newOverlap, remaining := parser.ExtractOverlapWords(processedChunk)
```

**Extracts the last `OVERLAP_WORDS` (20) words** to maintain context for the next chunk.

#### Step 6: Write Remaining Text
```go
if isFirstChunk {
    err = exporter.WriteChunk(outputPath, remaining)
    isFirstChunk = false
} else {
    err = exporter.AppendChunk(outputPath, remaining)
}
```

**First chunk creates file, subsequent chunks append.**

### Memory Efficiency in Chunked Processing

**Configurable Constant Memory Usage:**
- **Chunk size**: `config.CHUNK_BYTES` (1KB-8KB, default 4KB)
- **Overlap context**: `config.OVERLAP_WORDS` × ~20 chars (10-20 words, default 20 = ~400 bytes)
- **Processing buffers**: Transformer uses ~2.5KB (80 tokens × ~32 bytes)

**Memory is predictable and constant** regardless of file size.

**No growing data structures** - memory usage is predictable and constant.

## Error Handling Strategy

### Error Wrapping with Context

**Every error includes operation context:**

```go
// File operations
return fmt.Errorf("failed to read chunk at offset %d: %w", offset, err)
return fmt.Errorf("failed to write chunk: %w", err)

// File validation
return fmt.Errorf("input file does not exist: %s", inputPath)
return fmt.Errorf("failed to get file info: %w", err)
```

### Error Propagation Chain

```
Controller → Parser → Error
Controller → Transformer → Error  
Controller → Exporter → Error
```

**Controller wraps all errors** with context about which operation failed and where.

### Early Return Pattern

```go
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
// Continue only if no error
```

**No nested error handling** - clean, readable code.

## Complete Processing Example

### Small File (≤ 4KB)
```
Input: "hello (up) world"

1. processSingleChunk()
2. ReadChunk() → "hello (up) world"
3. ProcessText() → "HELLO world"  
4. WriteChunk() → File created
```

### Large File (> 4KB)
```
Input: 10KB file with "these words (up, 2)" command

Chunk 1: "word1 word2 ... word100"
- Process: "word1 word2 ... word100"
- Overlap: "word99 word100"
- Write: "word1 word2 ... word98"

Chunk 2: "word101 word102 ... these words (up, 2)"  
- Merge: "word99 word100 word101 ... these words (up, 2)"
- Process: "word99 word100 word101 ... THESE WORDS"
- Remove overlap: "word101 ... THESE WORDS"
- Write: "word101 ... THESE WORDS"
```

## Performance Characteristics

### File Size Handling (Configurable Thresholds)
- **Small files (< CHUNK_BYTES)**: Single-pass, maximum speed
- **Medium files (CHUNK_BYTES - 100MB)**: Chunked with overlap, constant memory
- **Large files (100MB+)**: Same constant memory, no performance degradation
- **Very large files (1GB+)**: Linear time complexity, constant space complexity

**Default threshold**: 4KB, **configurable**: 1KB-8KB

### Memory Usage (Configurable)
- **Small files**: ~2KB (transformer only)
- **Large files**: ~2KB to ~11KB (depending on config settings)
- **Default configuration**: ~7KB total
- **Scalability**: O(1) space complexity

**User can adjust memory usage** by changing `config.CHUNK_BYTES` (1KB-8KB) and `config.OVERLAP_WORDS` (10-20).

### Time Complexity
- **Processing**: O(n) where n = file size
- **Memory access**: O(1) - no growing data structures
- **I/O operations**: Minimized through efficient chunking

## Integration with Other Components

### Parser Integration
```go
data, err := parser.ReadChunk(inputPath, offset)
newOverlap, remaining := parser.ExtractOverlapWords(processedChunk)
```

**Controller uses parser for:**
- File reading with offset management
- Overlap word extraction
- UTF-8 safe chunk boundaries

### Transformer Integration
```go
result := transformer.ProcessText(textToProcess)
```

**Controller provides transformer with:**
- Complete context (original + overlap)
- Single string to process
- No knowledge of chunking (transformer is stateless)

### Exporter Integration
```go
err = exporter.WriteChunk(outputPath, result)      // First chunk
err = exporter.AppendChunk(outputPath, remaining)  // Subsequent chunks
```

**Controller manages file writing:**
- Creates file on first write
- Appends for subsequent writes
- Handles final overlap context

### Config Integration
```go
if fileInfo.Size() <= int64(config.CHUNK_BYTES) {
    // Use single chunk processing
}
```

**Controller uses config for:**
- Chunk size decisions
- Processing strategy selection
- Memory usage optimization

## Design Principles

### Single Responsibility
**Controller only orchestrates** - it doesn't:
- Parse file content (parser's job)
- Transform text (transformer's job)  
- Handle file I/O details (exporter's job)rter's job)

### Dependency Injection
**Controller receives all dependencies explicitly:**
- No global state
- Easy to test with mocks
- Clear component boundaries

### Error Transparency
**All errors bubble up with context:**
- Caller knows exactly what failed
- Original error preserved with `%w`
- Operation context added at each level

### Memory Efficiency
**Constant memory usage regardless of file size:**
- Fixed chunk sizes
- Overlap management
- No growing buffers

This makes the controller the perfect orchestrator - simple API, complex internal logic, optimal performance for all file sizes.