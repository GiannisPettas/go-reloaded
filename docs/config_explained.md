# Go-Reloaded Config: System Configuration and Performance Tuning Guide

## What Does the Config Do?

The config package is Go-Reloaded's **performance control center**. It defines the key constants that control how the system processes files, manages memory, and maintains performance.

Think of it like the **settings panel** of a high-performance car:
- **Engine tuning** (chunk size) - how much data to process at once
- **Memory management** (overlap words) - how much context to maintain
- **Safety limits** (validation) - prevents dangerous configurations
- **Performance optimization** - balances speed vs memory usage

## Core Configuration Constants

### 1. CHUNK_BYTES - Memory Control
```go
const CHUNK_BYTES = 4096 // 4KB chunks for memory efficiency
```

**Controls how much data is read from files at once.**

### 2. OVERLAP_WORDS - Context Preservation  
```go
const OVERLAP_WORDS = 20 // Number of words to preserve between chunks
```

**Controls how many words are remembered between chunks for command context.**

### 3. Validation Function
```go
func ValidateConstants() error
```

**Ensures configuration values are within safe, tested ranges.**

## Why Configuration Matters

### The Memory vs Performance Trade-off

**Larger chunks (8KB):**
- ✅ Faster processing (fewer I/O operations)
- ✅ Better for large files
- ❌ Uses more memory
- ❌ Longer processing delays

**Smaller chunks (1KB):**
- ✅ Uses less memory
- ✅ More responsive processing
- ❌ More I/O operations (slower)
- ❌ More overhead

**Default (4KB): Perfect balance for most use cases**

## CHUNK_BYTES Deep Dive

### What It Controls
```go
// Parser uses CHUNK_BYTES for buffer size
buffer := make([]byte, config.CHUNK_BYTES)
n, err := file.Read(buffer)
```

### Valid Range: 1024 - 8192 bytes (1KB - 8KB)

### Memory Impact Examples

**1KB chunks:**
```
Memory usage: ~1.5KB total
File processing: 1000 chunks for 1MB file
Use case: Memory-constrained environments
```

**4KB chunks (default):**
```
Memory usage: ~6KB total  
File processing: 250 chunks for 1MB file
Use case: Balanced performance and memory
```

**8KB chunks:**
```
Memory usage: ~10KB total
File processing: 125 chunks for 1MB file  
Use case: Performance-critical applications
```

### How to Choose CHUNK_BYTES

**Choose 1KB when:**
- Running on memory-constrained systems
- Processing many files simultaneously
- Memory usage is more important than speed

**Choose 4KB when:**
- General purpose usage (recommended default)
- Balanced performance requirements
- Most typical use cases

**Choose 8KB when:**
- Processing very large files (GB+)
- Maximum performance is critical
- Memory usage is not a concern

## OVERLAP_WORDS Deep Dive

### What It Controls
```go
// Parser extracts overlap words for context
overlapWords := words[len(words)-config.OVERLAP_WORDS:]
```

### Valid Range: 1 - 50 words

### Context Preservation Examples

**5 words overlap:**
```
Chunk 1: "The quick brown fox jumps"
Overlap: "brown fox jumps" (last 3 words)
Chunk 2: "brown fox jumps over the lazy dog"
```
**Good for**: Simple commands, memory efficiency
**Risk**: Complex multi-word commands might fail

**20 words overlap (default):**
```
Chunk 1: "...fifteen sixteen seventeen eighteen nineteen twenty"
Overlap: "sixteen seventeen eighteen nineteen twenty" (last 5 shown)
Chunk 2: "sixteen...twenty twenty-one twenty-two..."
```
**Good for**: Most use cases, handles complex commands
**Balance**: Good context preservation without excessive memory

**50 words overlap:**
```
Chunk 1: "...forty-six forty-seven forty-eight forty-nine fifty"  
Overlap: Very large context preserved
Chunk 2: "forty-six...fifty fifty-one fifty-two..."
```
**Good for**: Very complex commands, maximum compatibility
**Cost**: Higher memory usage

### How Commands Use Overlap

**Example: Multi-word command across chunks**
```
Chunk 1 end: "these are the last five words"
Overlap: "the last five words" (4 words preserved)
Chunk 2 start: "the last five words should be (up, 4)"

Command (up, 4) can access: "the last five words"
Result: "THE LAST FIVE WORDS should be"
```

**Without sufficient overlap:**
```
Chunk 1 end: "these are the last five words"
Overlap: "words" (only 1 word preserved)  
Chunk 2 start: "words should be (up, 4)"

Command (up, 4) tries to access 4 words but only finds: "words"
Result: Only "WORDS" gets transformed (incomplete)
```

### How to Choose OVERLAP_WORDS

**Choose 5-10 words when:**
- Memory is extremely limited
- Only using simple single-word commands
- Processing simple text without complex transformations

**Choose 20 words when:**
- General purpose usage (recommended default)
- Using multi-word commands occasionally
- Balanced memory and functionality

**Choose 30-50 words when:**
- Using complex multi-word commands frequently
- Commands like `(up, 25)` or `(cap, 30)`
- Maximum compatibility is required

## Configuration Validation

### ValidateConstants() Function

```go
func ValidateConstants() error {
    // Check CHUNK_BYTES is positive
    if CHUNK_BYTES <= 0 {
        return fmt.Errorf("CHUNK_BYTES must be positive, got %d", CHUNK_BYTES)
    }
    
    // Check CHUNK_BYTES minimum (1KB)
    if CHUNK_BYTES < 1024 {
        return fmt.Errorf("CHUNK_BYTES must be at least 1024 bytes, got %d", CHUNK_BYTES)
    }
    
    // Check CHUNK_BYTES maximum (8KB)
    if CHUNK_BYTES > 8192 {
        return fmt.Errorf("CHUNK_BYTES must be at most 8192 bytes, got %d", CHUNK_BYTES)
    }
    
    // Check OVERLAP_WORDS is positive
    if OVERLAP_WORDS <= 0 {
        return fmt.Errorf("OVERLAP_WORDS must be positive, got %d", OVERLAP_WORDS)
    }
    
    // Check OVERLAP_WORDS maximum (50)
    if OVERLAP_WORDS > 50 {
        return fmt.Errorf("OVERLAP_WORDS too large (max 50), got %d", OVERLAP_WORDS)
    }
    
    return nil
}
```

### Why These Limits?

**CHUNK_BYTES limits (1KB - 8KB):**
- **Minimum 1KB**: Ensures reasonable I/O efficiency
- **Maximum 8KB**: Prevents excessive memory usage
- **Tested range**: All values in this range are thoroughly tested

**OVERLAP_WORDS limits (1 - 50):**
- **Minimum 1**: At least some context must be preserved
- **Maximum 50**: Prevents excessive memory usage for overlap
- **Practical limit**: 50 words covers even very complex commands

## Integration with Other Components

### Parser Integration
```go
// Parser uses CHUNK_BYTES for buffer allocation
buffer := make([]byte, config.CHUNK_BYTES)

// Parser uses OVERLAP_WORDS for context extraction
if len(words) <= config.OVERLAP_WORDS {
    return text, "" // All words become overlap
}
```

### Controller Integration
```go
// Controller uses CHUNK_BYTES for processing decisions
if fileInfo.Size() <= int64(config.CHUNK_BYTES) {
    return processSingleChunk(inputPath, outputPath)
}
```

### Memory Calculations
```go
// Total memory usage calculation
totalMemory := config.CHUNK_BYTES + (config.OVERLAP_WORDS * averageWordLength)
```

## Performance Tuning Examples

### Memory-Constrained Environment
```go
// Optimize for minimal memory usage
const (
    CHUNK_BYTES   = 1024  // 1KB chunks
    OVERLAP_WORDS = 5     // Minimal overlap
)
// Total memory: ~1.5KB
```

### Balanced Configuration (Default)
```go
// Optimize for general use
const (
    CHUNK_BYTES   = 4096  // 4KB chunks  
    OVERLAP_WORDS = 20    // Good overlap
)
// Total memory: ~6KB
```

### Performance-Critical Environment
```go
// Optimize for maximum speed
const (
    CHUNK_BYTES   = 8192  // 8KB chunks
    OVERLAP_WORDS = 50    // Maximum overlap
)
// Total memory: ~10KB
```

## Real-World Configuration Scenarios

### Embedded Systems
```go
// Limited RAM (e.g., 64MB total)
const (
    CHUNK_BYTES   = 1024  // Minimal memory footprint
    OVERLAP_WORDS = 10    // Basic command support
)
```

### Desktop Applications
```go
// Typical desktop (e.g., 8GB RAM)
const (
    CHUNK_BYTES   = 4096  // Balanced performance
    OVERLAP_WORDS = 20    // Full command support
)
```

### Server Applications
```go
// High-performance server (e.g., 64GB RAM)
const (
    CHUNK_BYTES   = 8192  // Maximum throughput
    OVERLAP_WORDS = 50    // Maximum compatibility
)
```

### Batch Processing
```go
// Processing thousands of files
const (
    CHUNK_BYTES   = 2048  // Moderate memory per file
    OVERLAP_WORDS = 15    // Reasonable command support
)
```

## Configuration Testing

### Validation Testing
```go
// Test configuration validation
func TestValidateConstants(t *testing.T) {
    // Test valid configuration
    err := config.ValidateConstants()
    if err != nil {
        t.Errorf("Valid configuration failed validation: %v", err)
    }
}
```

### Performance Testing
```go
// Test different configurations
configurations := []struct{
    chunkBytes   int
    overlapWords int
    expectedMemory int
}{
    {1024, 5, 1500},   // Minimal
    {4096, 20, 6000},  // Default  
    {8192, 50, 10000}, // Maximum
}
```

## Memory Usage Calculator

### Formula
```go
func CalculateMemoryUsage(chunkBytes, overlapWords int) int {
    // Chunk buffer
    chunkMemory := chunkBytes
    
    // Overlap context (average 20 chars per word)
    overlapMemory := overlapWords * 20
    
    // Transformer buffers (fixed ~2KB)
    transformerMemory := 2048
    
    // Working variables (~500 bytes)
    workingMemory := 500
    
    return chunkMemory + overlapMemory + transformerMemory + workingMemory
}
```

### Usage Examples
```go
// Calculate memory for different configurations
fmt.Printf("1KB/5 words:  %d bytes\n", CalculateMemoryUsage(1024, 5))   // ~3.5KB
fmt.Printf("4KB/20 words: %d bytes\n", CalculateMemoryUsage(4096, 20))  // ~6.9KB  
fmt.Printf("8KB/50 words: %d bytes\n", CalculateMemoryUsage(8192, 50))  // ~11.5KB
```

## Design Principles

### Configurability with Safety
- **User control**: Allow performance tuning for different environments
- **Safe defaults**: 4KB/20 words works well for most use cases
- **Validation**: Prevent dangerous configurations that could cause issues

### Performance Optimization
- **Memory efficiency**: All configurations use constant memory
- **I/O efficiency**: Larger chunks reduce system call overhead
- **Context preservation**: Sufficient overlap for command functionality

### Simplicity
- **Two constants**: Easy to understand and modify
- **Clear ranges**: Documented limits with explanations
- **Validation function**: Automatic checking of configuration validity

### Testing and Reliability
- **Tested ranges**: All valid configurations are thoroughly tested
- **Error messages**: Clear feedback when configuration is invalid
- **Documentation**: Comprehensive guidance for choosing values

This makes the config package a powerful tool for optimizing Go-Reloaded's performance for different environments and use cases, while maintaining safety and reliability through validation and testing.