# Go-Reloaded Exporter: File Writing and Output Management Guide

## What Does the Exporter Do?

The exporter is Go-Reloaded's **file writing specialist**. It handles the final step of the processing pipeline - taking transformed text and efficiently writing it to output files.

Think of it like a **smart secretary** who:
- **Creates folders** automatically when needed (directory management)
- **Writes documents** efficiently without wasting resources
- **Appends to existing documents** when building large files piece by piece
- **Handles file permissions** properly for security
- **Never loses data** through proper error handling

## Core Responsibilities

### 1. File Creation and Writing
```go
func WriteChunk(filePath, content string) error
```

**Creates new files** or **overwrites existing ones** with processed content.

### 2. File Appending
```go
func AppendChunk(filePath, content string) error
```

**Adds content to existing files** - essential for chunked processing of large files.

### 3. Directory Management
Both functions automatically **create directories** if they don't exist.

## Why Two Different Functions?

### WriteChunk() - "Start Fresh"
```go
// First chunk of processing - create new file
err := exporter.WriteChunk("output.txt", "First chunk of text")
// Creates new file, overwrites if exists
```

### AppendChunk() - "Add More"
```go
// Subsequent chunks - add to existing file
err := exporter.AppendChunk("output.txt", "Second chunk of text")
err := exporter.AppendChunk("output.txt", "Third chunk of text")
// Builds file incrementally
```

**Result**: `output.txt` contains "First chunk of textSecond chunk of textThird chunk of text"

## Step-by-Step Process

### Step 1: WriteChunk() - Create New File

```go
func WriteChunk(filePath, content string) error {
    // 1. Extract directory path from file path
    dir := filepath.Dir(filePath) // "/path/to/file.txt" → "/path/to"
    
    // 2. Create directory structure if needed
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory %s: %w", dir, err)
    }
    
    // 3. Write content to file (overwrites existing)
    err := os.WriteFile(filePath, []byte(content), 0644)
    if err != nil {
        return fmt.Errorf("failed to write file %s: %w", filePath, err)
    }
    
    return nil
}
```

**Key Features:**
- **Directory creation**: Automatically creates parent directories
- **File overwriting**: Replaces existing file content completely
- **Atomic operation**: Uses `os.WriteFile` for safe, atomic writes
- **Proper permissions**: Sets file permissions to 0644 (readable by all, writable by owner)

### Step 2: AppendChunk() - Add to Existing File

```go
func AppendChunk(filePath, content string) error {
    // 1. Create directory structure if needed
    dir := filepath.Dir(filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory %s: %w", dir, err)
    }
    
    // 2. Open file for appending (create if doesn't exist)
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open file %s for appending: %w", filePath, err)
    }
    defer file.Close() // Always close file when done
    
    // 3. Write content to end of file
    _, err = file.WriteString(content)
    if err != nil {
        return fmt.Errorf("failed to append to file %s: %w", filePath, err)
    }
    
    return nil
}
```

**Key Features:**
- **Append mode**: Adds content to end of existing file
- **Create if missing**: Creates file if it doesn't exist
- **Proper file handling**: Opens, writes, and closes file safely
- **Error context**: Detailed error messages for debugging

## File Permissions Explained

### Directory Permissions: 0755
```
Owner: rwx (read, write, execute) = 7
Group: r-x (read, execute)        = 5  
Other: r-x (read, execute)        = 5
```
**Meaning**: Owner can do anything, others can read and navigate the directory.

### File Permissions: 0644
```
Owner: rw- (read, write)     = 6
Group: r-- (read only)       = 4
Other: r-- (read only)       = 4
```
**Meaning**: Owner can read/write, others can only read the file.

**Why These Permissions?**
- **Security**: Prevents unauthorized modification
- **Accessibility**: Allows others to read the output
- **Standard practice**: Common permissions for text files

## Integration with Chunked Processing

### Single File Processing
```go
// Small file - single chunk
content := transformer.ProcessText(inputText)
err := exporter.WriteChunk("output.txt", content)
```

### Large File Processing
```go
// Large file - multiple chunks
isFirstChunk := true

for each chunk {
    processedChunk := transformer.ProcessText(chunk)
    
    if isFirstChunk {
        err = exporter.WriteChunk("output.txt", processedChunk)
        isFirstChunk = false
    } else {
        err = exporter.AppendChunk("output.txt", processedChunk)
    }
}
```

**Why This Pattern?**
- **Clean start**: First chunk creates fresh file
- **Incremental building**: Subsequent chunks append
- **Memory efficiency**: Never holds entire output in memory
- **Streaming output**: Results written immediately

## Directory Management Deep Dive

### Automatic Directory Creation

**The Problem:**
```go
// This fails if "/deep/nested/path/" doesn't exist
err := os.WriteFile("/deep/nested/path/file.txt", data, 0644)
// Error: no such file or directory
```

**The Solution:**
```go
// Exporter automatically creates the path
dir := filepath.Dir("/deep/nested/path/file.txt") // "/deep/nested/path"
err := os.MkdirAll(dir, 0755) // Creates entire path
// Now file writing succeeds
```

### MkdirAll() Behavior
```go
os.MkdirAll("/path/to/deep/directory", 0755)
```

**Creates entire path:**
```
/path/           (if doesn't exist)
/path/to/        (if doesn't exist)  
/path/to/deep/   (if doesn't exist)
/path/to/deep/directory/ (if doesn't exist)
```

**Safe to call multiple times** - no error if directories already exist.

## Error Handling Strategy

### Comprehensive Error Context

**Directory Creation Errors:**
```go
if err := os.MkdirAll(dir, 0755); err != nil {
    return fmt.Errorf("failed to create directory %s: %w", dir, err)
}
```

**File Writing Errors:**
```go
err := os.WriteFile(filePath, []byte(content), 0644)
if err != nil {
    return fmt.Errorf("failed to write file %s: %w", filePath, err)
}
```

**File Opening Errors:**
```go
file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    return fmt.Errorf("failed to open file %s for appending: %w", filePath, err)
}
```

**Benefits:**
- **Clear context**: Know exactly which operation failed
- **File path included**: Easy to identify problematic files
- **Original error preserved**: Full error chain available for debugging

### Common Error Scenarios

**Permission Denied:**
```
failed to write file /root/protected.txt: permission denied
```
**Solution**: Check file/directory permissions or run with appropriate privileges.

**Disk Full:**
```
failed to append to file /tmp/large.txt: no space left on device
```
**Solution**: Free up disk space or write to different location.

**Invalid Path:**
```
failed to create directory /invalid\path: invalid argument
```
**Solution**: Use valid file path characters for the operating system.

## Memory Efficiency

### Streaming Writes
```go
// BAD - Builds entire output in memory first
var result strings.Builder
for each chunk {
    processed := transformer.ProcessText(chunk)
    result.WriteString(processed) // Accumulates in memory
}
os.WriteFile("output.txt", []byte(result.String()), 0644) // Write all at once
```

```go
// GOOD - Writes immediately, constant memory
isFirst := true
for each chunk {
    processed := transformer.ProcessText(chunk)
    if isFirst {
        exporter.WriteChunk("output.txt", processed) // Write immediately
        isFirst = false
    } else {
        exporter.AppendChunk("output.txt", processed) // Append immediately
    }
}
```

**Memory Usage:**
- **Bad approach**: Memory grows with output size (could be GB)
- **Good approach**: Constant memory (~KB) regardless of output size

### File Handle Management

**AppendChunk() properly manages file handles:**
```go
file, err := os.OpenFile(filePath, flags, 0644)
if err != nil {
    return err
}
defer file.Close() // Ensures file is closed even if error occurs

_, err = file.WriteString(content)
return err
```

**Why `defer file.Close()` is important:**
- **Resource cleanup**: Prevents file handle leaks
- **Data integrity**: Ensures data is flushed to disk
- **Error safety**: Closes file even if write operation fails

## Integration with Other Components

### Controller Integration
```go
// Controller orchestrates the write process
if isFirstChunk {
    err = exporter.WriteChunk(outputPath, processedText)
    isFirstChunk = false
} else {
    err = exporter.AppendChunk(outputPath, processedText)
}
```

### Transformer Integration
```go
// Exporter receives processed text from transformer
transformedText := transformer.ProcessText(inputChunk)
err := exporter.WriteChunk("output.txt", transformedText)
```

### Parser Integration (Indirect)
```go
// Complete pipeline: Parser → Transformer → Exporter
chunk := parser.ReadChunk(inputFile, offset)
processed := transformer.ProcessText(string(chunk))
err := exporter.AppendChunk(outputFile, processed)
```

## File Operations Comparison

### WriteChunk vs AppendChunk

| Feature | WriteChunk | AppendChunk |
|---------|------------|-------------|
| **Purpose** | Create new file | Add to existing file |
| **Existing file** | Overwrites | Preserves and extends |
| **Use case** | First chunk | Subsequent chunks |
| **Performance** | Fast (single operation) | Fast (single append) |
| **Safety** | Atomic write | Incremental build |

### When to Use Each

**WriteChunk:**
- Starting new output file
- Small files (single chunk)
- Replacing existing file completely
- Atomic operations needed

**AppendChunk:**
- Building large files incrementally
- Continuing chunked processing
- Preserving existing content
- Streaming output

## Performance Characteristics

### Time Complexity
- **WriteChunk**: O(n) where n = content size
- **AppendChunk**: O(n) where n = content size
- **Directory creation**: O(d) where d = directory depth

### Space Complexity
- **Memory usage**: O(1) - constant regardless of file size
- **Disk usage**: O(n) where n = total content written

### Scalability
- **Small files**: Microsecond writes
- **Large files**: Linear time, constant memory
- **Many files**: Parallel processing possible

## Real-World Examples

### Processing Log Files
```go
// Process large log file in chunks
for offset := 0; offset < logSize; offset += chunkSize {
    chunk := parser.ReadChunk("app.log", offset)
    processed := transformer.ProcessText(string(chunk))
    
    if offset == 0 {
        exporter.WriteChunk("processed.log", processed)
    } else {
        exporter.AppendChunk("processed.log", processed)
    }
}
```

### Batch Document Processing
```go
// Process multiple documents into single output
isFirst := true
for _, inputFile := range documents {
    content := readAndProcess(inputFile)
    
    if isFirst {
        exporter.WriteChunk("combined.txt", content)
        isFirst = false
    } else {
        exporter.AppendChunk("combined.txt", content)
    }
}
```

### Nested Directory Output
```go
// Automatically creates deep directory structure
err := exporter.WriteChunk("reports/2024/january/summary.txt", reportData)
// Creates: reports/ → reports/2024/ → reports/2024/january/ → summary.txt
```

## Design Principles

### Simplicity First
- **Two functions**: Clear separation of create vs append operations
- **Automatic setup**: Directory creation handled transparently
- **Minimal API**: Easy to understand and use correctly

### Safety and Reliability
- **Error handling**: Comprehensive error reporting with context
- **Resource management**: Proper file handle cleanup
- **Atomic operations**: Safe file writing practices

### Performance Optimization
- **Streaming writes**: No memory accumulation
- **Efficient I/O**: Direct file operations without buffering overhead
- **Minimal overhead**: Simple, fast implementations

### Integration Friendly
- **Controller compatibility**: Works seamlessly with chunked processing
- **Error propagation**: Consistent error handling throughout pipeline
- **Flexible usage**: Supports both single-file and multi-chunk scenarios

This makes the exporter a reliable, efficient foundation for Go-Reloaded's output generation, ensuring processed text is written safely and efficiently regardless of file size or complexity.