# Go-Reloaded Product Overview

## Project Purpose
Go-Reloaded is a high-performance text processing tool that transforms text files using various commands. It's designed to handle files of any size with minimal memory usage through a dual finite state machine architecture.

## Key Features

### Text Transformations
- **Numeric Base Conversion**: Convert hexadecimal and binary numbers to decimal (supports negative numbers)
- **Case Transformations**: Change text to uppercase, lowercase, or capitalize words
- **Article Correction**: Automatically fix "a/an" usage based on vowel sounds
- **Punctuation Spacing**: Fix spacing around punctuation marks
- **Quote Repositioning**: Properly position single quotes around words
- **Command Chaining**: Apply multiple transformations to the same word

### Performance & Reliability
- **Memory Efficient**: Processes files of any size using only ~8KB of memory
- **Error Resilience**: Invalid commands are gracefully ignored
- **Zero Dependencies**: Uses only Go standard library, no external packages
- **UTF-8 Safe**: Handles international characters without corruption

### File Size Handling
- Small files (< 4KB): Single-pass processing
- Medium files (4KB - 100MB): Chunked processing with overlap
- Large files (100MB+): Constant memory usage (~8KB)
- Very large files (1GB+): No memory limitations

## Target Users
- Developers needing text preprocessing tools
- Data processing professionals working with large text files
- System administrators requiring memory-efficient text transformation
- Anyone needing reliable, fast text processing with specific formatting rules

## Use Cases
- Batch text formatting and normalization
- Large file processing with memory constraints
- Text preprocessing for data pipelines
- Document formatting automation
- Educational tool for text processing algorithms

## Value Proposition
- **Performance**: Dual FSM architecture enables single-pass processing
- **Scalability**: Handles files from KB to GB with constant memory usage
- **Reliability**: Comprehensive test suite with 29 golden test cases
- **Simplicity**: Zero dependencies, easy deployment and integration
- **Flexibility**: Command chaining allows complex transformations