# Go-Reloaded Product Overview

## Purpose
Go-Reloaded is a high-performance text processing tool that transforms text files using various commands. It's designed for efficient text manipulation with minimal memory usage, capable of processing files of any size using only ~8KB of memory.

## Key Features

### Core Transformations
- **Numeric Base Conversion**: Convert hexadecimal and binary numbers to decimal format
- **Case Transformations**: Change text to uppercase, lowercase, or capitalize words
- **Article Correction**: Automatically fix "a/an" usage based on vowel sounds
- **Punctuation Spacing**: Fix spacing around punctuation marks
- **Command Chaining**: Apply multiple transformations to the same word sequentially

### Performance Characteristics
- **Memory Efficient**: Processes files of any size using only ~8KB of memory
- **Single-Pass Processing**: No multiple iterations over data
- **Dual FSM Architecture**: Two finite state machines working in tandem for maximum efficiency
- **UTF-8 Safe**: Handles international characters without corruption
- **Zero Dependencies**: Uses only Go standard library, no external packages

## Target Users
- Developers needing text preprocessing for data pipelines
- Content creators requiring bulk text transformations
- System administrators processing log files or configuration files
- Anyone needing efficient, memory-conscious text processing

## Use Cases
- Converting numeric formats in technical documentation
- Standardizing text case in large datasets
- Correcting grammar issues in bulk text processing
- Formatting text files for consistent punctuation
- Processing large files without memory constraints

## Value Proposition
- **Efficiency**: Constant memory usage regardless of file size
- **Reliability**: Comprehensive test suite with 27 golden test cases
- **Simplicity**: Single binary with straightforward command-line interface
- **Performance**: Optimized for speed with dual FSM architecture
- **Portability**: Pure Go implementation with no external dependencies