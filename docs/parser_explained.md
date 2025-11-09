# Go-Reloaded Parser: File Reading and Memory Management Guide

## What Does the Parser Do?

The parser is Go-Reloaded's **file reading specialist**. It handles the complex task of reading files efficiently while maintaining constant memory usage, no matter how big the file is.

Think of it like a **smart librarian** who can read any book (file) by:
- Reading it **page by page** (chunks) instead of loading the whole book into memory
- Making sure **words aren't cut in half** between pages (UTF-8 safety)
- **Remembering the last few words** from each page to maintain context
- **Never using more than a small desk** (constant memory) regardless of book size

## Core Responsibilities

### 1. Chunked File Reading
```go
func ReadChunk(filepath string, offset int64) ([]byte, error)
```

**Reads files in small, manageable pieces** instead of loading everything at once.

### 2. UTF-8 Safety
```go
func AdjustToRuneBoundary(data []byte) []byte
```

**Prevents character corruption** by ensuring chunks end at complete Unicode characters.

### 3. Context Preservation
```go
func ExtractOverlapWords(text string) (overlap, remaining string)
```

**Maintains word context** between chunks so commands can work across chunk boundaries.

## Why Chunked Reading?

### The Problem: Large Files
```go
// BAD - Loads entire file into memory
content, err := os.ReadFile("huge_file.txt") // Could be 1GB!
// Your program now uses 1GB+ of RAM
```

### The Solution: Chunked Reading
```go
// GOOD - Reads only 4KB at a time
chunk, err := parser.ReadChunk("huge_file.txt", offset) // Always 4KB max
// Your program uses constant ~4KB regardless of file size
```

## Step-by-Step Process

### Step 1: ReadChunk() - Smart File Reading

```go
func ReadChunk(filepath string, offset int64) ([]byte, error) {
    // 1. Open file
    file, err := os.Open(filepath)
    
    // 2. Jump to specific position (for large files)
    file.Seek(offset, io.SeekStart)
    
    // 3. Read exactly CHUNK_BYTES (4KB by default)
    buffer := make([]byte, config.CHUNK_BYTES)
    n, err := file.Read(buffer)
    
    // 4. Return only bytes actually read
    chunk := buffer[:n]
    
    // 5. Make sure we don't cut characters in half
    return AdjustToRuneBoundary(chunk)
}
```

**Key Features:**
- **Offset-based reading**: Can start reading from any position in the file
- **Fixed buffer size**: Always allocates exactly `CHUNK_BYTES` (4KB)
- **Actual bytes returned**: Only returns bytes that were actually read
- **UTF-8 safe**: Ensures no character corruption

### Step 2: AdjustToRuneBoundary() - UTF-8 Safety

**The Problem:**
```
File content: "Hello 世界 World"
Chunk boundary cuts here: "Hello 世|界 World"
                              ↑
                         Cuts UTF-8 character in half!
```

**The Solution:**
```go
func AdjustToRuneBoundary(data []byte) []byte {
    // Check if the chunk is valid UTF-8
    if utf8.Valid(data) {
        return data // All good!
    }
    
    // Work backwards to find last complete character
    for i := len(data) - 1; i >= 0; i-- {
        if utf8.Valid(data[:i+1]) {
            return data[:i+1] // Return up to last complete character
        }
    }
    
    return []byte{} // Fallback: return empty if no valid UTF-8
}
```

**Example:**
```
Original chunk: [72, 101, 108, 108, 111, 32, 228, 184] // "Hello " + incomplete UTF-8
After adjustment: [72, 101, 108, 108, 111, 32]         // "Hello " (safe)
```

**Why This Matters:**
- **Prevents crashes**: Invalid UTF-8 can cause string operations to fail
- **Preserves characters**: International characters stay intact
- **Maintains data integrity**: No corrupted text in output

### Step 3: ExtractOverlapWords() - Context Preservation

**The Challenge:**
```
Chunk 1: "these three words should"
Chunk 2: "be uppercase (up, 4)"
```

The `(up, 4)` command needs to access "these three words should" from the previous chunk!

**The Solution - Overlap Context:**
```go
func ExtractOverlapWords(text string) (overlap, remaining string) {
    words := strings.Fields(text) // Split into words
    
    if len(words) <= config.OVERLAP_WORDS {
        return text, "" // All words become overlap
    }
    
    // Split: most words go to output, last N words become overlap
    remainingWords := words[:len(words)-config.OVERLAP_WORDS]
    overlapWords := words[len(words)-config.OVERLAP_WORDS:]
    
    remaining = strings.Join(remainingWords, " ")
    overlap = strings.Join(overlapWords, " ")
    
    return overlap, remaining
}
```

**Example with OVERLAP_WORDS = 3:**
```
Input text: "The quick brown fox jumps over the lazy dog"
Words: ["The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"]

Result:
- remaining: "The quick brown fox jumps over"  (written to file)
- overlap: "the lazy dog"                      (saved for next chunk)
```

**Why Overlap is Critical:**
- **Command context**: Commands can reference words from previous chunks
- **Transformation continuity**: Multi-word commands work across boundaries
- **Memory efficiency**: Only saves last N words, not entire chunk

## Memory Efficiency Deep Dive

### Constant Memory Usage

**No matter the file size, parser uses:**
- **Chunk buffer**: 4KB (configurable 1KB-8KB)
- **Overlap context**: ~400 bytes (20 words × ~20 chars average)
- **Working variables**: ~100 bytes
- **Total**: ~4.5KB constant memory usage

### Memory Comparison

```go
// Traditional approach (scales with file size)
func ProcessFileTraditional(filename string) {
    content, _ := os.ReadFile(filename)  // 1GB file = 1GB RAM
    result := transform(string(content))  // Another 1GB for string
    os.WriteFile("output.txt", []byte(result), 0644) // Another 1GB
    // Total: 3GB RAM for 1GB file!
}

// Go-Reloaded approach (constant memory)
func ProcessFileChunked(filename string) {
    for offset := 0; offset < fileSize; offset += chunkSize {
        chunk := parser.ReadChunk(filename, offset)     // 4KB
        processed := transformer.ProcessText(chunk)      // 4KB
        exporter.AppendChunk("output.txt", processed)   // Write immediately
        // Total: 8KB RAM for any size file!
    }
}
```

## File Size Handling Examples

### Small File (< 4KB)
```
File: "Hello (up) world!"
Process: Single chunk, no overlap needed
Memory: ~4KB
```

### Medium File (100KB)
```
File: 100KB text file
Process: ~25 chunks with overlap between each
Memory: ~4KB (constant)
Chunks: 1→2→3→...→25 (sequential processing)
```

### Large File (1GB)
```
File: 1GB text file  
Process: ~250,000 chunks with overlap
Memory: ~4KB (still constant!)
Time: Linear with file size, but memory stays constant
```

## UTF-8 and International Characters

### Why UTF-8 Safety Matters

**UTF-8 characters can be 1-4 bytes:**
- `A` = 1 byte: `[65]`
- `é` = 2 bytes: `[195, 169]`
- `世` = 3 bytes: `[228, 184, 150]`
- `🚀` = 4 bytes: `[240, 159, 154, 128]`

**If chunk boundary cuts a multi-byte character:**
```
Correct: "Hello 世界"
Broken:  "Hello 世|界" → "Hello 世�界" (corruption!)
```

### AdjustToRuneBoundary() in Action

```go
// Example: Chunk ends in middle of "世" character
chunk := []byte{72, 101, 108, 108, 111, 32, 228, 184} // "Hello " + incomplete "世"

// utf8.Valid(chunk) returns false - incomplete character detected

// Work backwards:
// utf8.Valid(chunk[:8]) = false (still incomplete)
// utf8.Valid(chunk[:7]) = false (still incomplete) 
// utf8.Valid(chunk[:6]) = true  (ends at space after "Hello")

// Return: []byte{72, 101, 108, 108, 111, 32} = "Hello "
```

**Result**: The incomplete "世" character will be read in the next chunk, preventing corruption.

## Integration with Other Components

### Controller Integration
```go
// Controller orchestrates chunked reading
for offset < fileSize {
    chunk, err := parser.ReadChunk(inputPath, offset)     // Parser reads
    processed := transformer.ProcessText(string(chunk))   // Transformer processes
    err = exporter.AppendChunk(outputPath, processed)     // Exporter writes
    offset += int64(len(chunk))
}
```

### Transformer Integration
```go
// Parser provides clean, UTF-8 safe text to transformer
chunk := parser.ReadChunk(file, offset)
text := string(chunk) // Safe conversion - no UTF-8 corruption
result := transformer.ProcessText(text)
```

### Config Integration
```go
// Parser respects user-configurable settings
buffer := make([]byte, config.CHUNK_BYTES)    // 1KB-8KB configurable
overlapWords := config.OVERLAP_WORDS          // 1-50 words configurable
```

## Error Handling Strategy

### File Operation Errors
```go
file, err := os.Open(filepath)
if err != nil {
    return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
}
```

### Seek Operation Errors
```go
_, err = file.Seek(offset, io.SeekStart)
if err != nil {
    return nil, fmt.Errorf("failed to seek to offset %d: %w", offset, err)
}
```

### Read Operation Errors
```go
n, err := file.Read(buffer)
if err != nil && err != io.EOF {
    return nil, fmt.Errorf("failed to read from file: %w", err)
}
```

**Key Principles:**
- **Context preservation**: Every error includes operation details
- **Error wrapping**: Original errors preserved with `%w`
- **Graceful EOF handling**: End-of-file is expected, not an error

## Performance Characteristics

### Time Complexity
- **File reading**: O(n) where n = file size
- **UTF-8 adjustment**: O(k) where k = bytes to check (typically < 10)
- **Overlap extraction**: O(w) where w = number of words in chunk
- **Overall**: Linear time, constant space

### Space Complexity
- **Chunk buffer**: O(1) - fixed size regardless of file size
- **Overlap storage**: O(1) - fixed number of words
- **Working variables**: O(1) - constant overhead
- **Total**: O(1) space complexity

### Scalability
- **1KB file**: ~4KB memory usage
- **1MB file**: ~4KB memory usage  
- **1GB file**: ~4KB memory usage
- **100GB file**: ~4KB memory usage

**The parser scales to any file size with constant memory usage!**

## Common Use Cases

### Processing Log Files
```go
// Can process multi-GB log files with constant memory
for offset := 0; offset < logFileSize; offset += chunkSize {
    chunk := parser.ReadChunk("app.log", offset)
    // Process each chunk independently
}
```

### Text Transformation
```go
// Transform large documents without loading into memory
chunk := parser.ReadChunk("document.txt", offset)
transformed := transformer.ProcessText(string(chunk))
exporter.AppendChunk("output.txt", transformed)
```

### Data Migration
```go
// Convert file formats for large datasets
for each chunk {
    data := parser.ReadChunk(inputFile, offset)
    converted := convertFormat(data)
    writeToNewFormat(converted)
}
```

## Design Principles

### Memory Efficiency First
- **Fixed buffers**: No growing data structures
- **Immediate processing**: Process and discard chunks
- **Minimal overhead**: Only essential data kept in memory

### UTF-8 Safety Always
- **Character integrity**: Never corrupt international characters
- **Boundary awareness**: Always end chunks at character boundaries
- **Validation**: Check UTF-8 validity before processing

### Context Preservation
- **Command continuity**: Maintain context for multi-word commands
- **Overlap management**: Smart word boundary detection
- **State maintenance**: Preserve processing state between chunks

### Error Transparency
- **Detailed errors**: Include file paths, offsets, and operation context
- **Error chaining**: Preserve original error information
- **Graceful degradation**: Handle edge cases without crashing

This makes the parser a robust, efficient foundation for Go-Reloaded's file processing capabilities, enabling it to handle files of any size with predictable performance and memory usage.