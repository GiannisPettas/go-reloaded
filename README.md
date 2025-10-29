# Go-Reloaded

A high-performance text processing application that transforms text files using various rules including numeric conversions, case modifications, punctuation corrections, and article adjustments.

## Features

### Text Transformations
- **Numeric Conversions**: Convert hexadecimal and binary numbers to decimal
- **Case Modifications**: Transform text to uppercase, lowercase, or capitalize
- **Punctuation Corrections**: Automatically fix spacing around punctuation marks
- **Quote Repositioning**: Move punctuation inside quotes where appropriate
- **Article Corrections**: Fix "a/an" usage based on vowel sounds
- **Command Chaining**: Apply multiple transformations to the same word

### Performance & Scalability
- **Memory Efficient**: Processes files of any size with constant memory usage (~8KB)
- **UTF-8 Safe**: Handles international characters without corruption
- **Stream Processing**: No temporary files created during processing
- **Error Resilient**: Invalid commands are ignored, processing continues

## Installation

### Prerequisites
- Go 1.19 or higher

### Build from Source
```bash
git clone <repository-url>
cd go-reloaded
go build -o go-reloaded cmd/go-reloaded/main.go
```

### Optimized Build (Smaller Binary)
```bash
go build -ldflags="-s -w" -o go-reloaded cmd/go-reloaded/main.go
```

## Usage

### Basic Usage
```bash
./go-reloaded input.txt output.txt
```

### Command Line Arguments
- `input.txt`: Path to the input text file
- `output.txt`: Path where the processed output will be saved

### Example
```bash
./go-reloaded sample.txt result.txt
```

## Transformation Rules

### 1. Numeric Conversions

#### Hexadecimal to Decimal
```
Input:  "The value is 1E (hex)"
Output: "The value is 30"
```

#### Binary to Decimal
```
Input:  "Binary 1010 (bin) equals decimal"
Output: "Binary 10 equals decimal"
```

### 2. Case Transformations

#### Single Word
```
Input:  "make this WORD (low) lowercase"
Output: "make this word lowercase"

Input:  "make this word (up) uppercase"
Output: "make this WORD uppercase"

Input:  "make this word (cap) capitalized"
Output: "make this Word capitalized"
```

#### Multiple Words
```
Input:  "These three words (up, 3) should be uppercase"
Output: "THESE THREE WORDS should be uppercase"

Input:  "These two words (cap, 2) should be capitalized"
Output: "These Two words should be capitalized"
```

### 3. Punctuation Corrections

#### Spacing Rules
```
Input:  "Hello , world !"
Output: "Hello, world!"

Input:  "What is this? It is great."
Output: "What is this? It is great."
```

### 4. Quote Repositioning

#### Punctuation Inside Quotes
```
Input:  "Hello there" , she said.
Output: "Hello there," she said.

Input:  "Are you sure" ?
Output: "Are you sure?"
```

### 5. Article Corrections

#### A/An Usage
```
Input:  "I need a apple and an car"
Output: "I need an apple and a car"

Input:  "It was a honor to meet an European"
Output: "It was an honor to meet a European"
```

### 6. Command Chaining

#### Multiple Transformations
```
Input:  "The number 1010 (bin) (hex) is interesting"
Output: "The number 16 is interesting"
```
*Explanation: 1010 (binary) → 10 (decimal) → 16 (hexadecimal)*

## Advanced Features

### Memory Efficiency
- Processes files of any size using only ~8KB of memory
- Uses chunked processing with smart overlap to maintain context
- No file size limitations

### UTF-8 Support
- Safely handles international characters and emojis
- Prevents character corruption at chunk boundaries
- Maintains text encoding integrity

### Error Handling
- Invalid commands are gracefully ignored
- Malformed syntax doesn't break processing
- Clear error messages for file operations

## Technical Implementation

### Architecture
- **Finite State Machine**: Robust text processing with state tracking
- **Chunked Processing**: Memory-efficient streaming for large files
- **Component-Based**: Modular design with Parser, Transformer, and Exporter
- **UTF-8 Boundary Alignment**: Prevents character corruption

### Performance Characteristics
- **Time Complexity**: O(n) where n is file size
- **Memory Usage**: O(1) constant memory regardless of file size
- **I/O Optimized**: Sequential reads and buffered writes

## Configuration

### System Constants
- **Chunk Size**: 4096 bytes (configurable in `internal/config/config.go`)
- **Word Overlap**: 20 words between chunks for context preservation

### Tuning Guidelines
- Increase chunk size for better I/O performance on large files
- Adjust overlap for better context preservation vs. performance trade-off

## Examples

### Sample Input File (`sample.txt`)
```
Simply add 1010 (bin) (hex) , and check the total !

I need a apple and an car for the trip.

"Hello there" , she said quietly .

Convert FF (hex) to decimal and make it (up) .
```

### Expected Output (`result.txt`)
```
Simply add 16, and check the total!

I need an apple and a car for the trip.

"Hello there," she said quietly.

Convert 255 to decimal and make it UP.
```

### Running the Example
```bash
./go-reloaded sample.txt result.txt
echo "Processing complete! Check result.txt"
```

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Specific Component Tests
```bash
go test ./internal/parser
go test ./internal/transformer
go test ./internal/exporter
go test ./internal/controller
```

### Golden Test Suite
The project includes 22 comprehensive test cases covering all transformation scenarios. See `docs/golden_tests.md` for detailed test specifications.

## Project Structure

```
go-reloaded/
├── cmd/go-reloaded/          # CLI application entry point
├── internal/
│   ├── config/               # System configuration
│   ├── parser/               # File reading and chunking
│   ├── transformer/          # Text transformation engine
│   ├── exporter/             # File writing
│   ├── controller/           # Workflow orchestration
│   └── testutils/            # Testing utilities
├── docs/                     # Documentation and analysis
└── reports/                  # Development progress reports
```

## Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Implement the feature
5. Ensure all tests pass
6. Submit a pull request

### Code Standards
- Follow Go conventions and best practices
- Maintain test coverage above 90%
- Use meaningful variable and function names
- Add comments for complex algorithms

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For technical documentation, see `docs/technical_architecture.md`.

For issues and questions, please create an issue in the repository.

## Performance Benchmarks

### File Size Handling
- ✅ Small files (< 1KB): Instant processing
- ✅ Medium files (1MB - 100MB): Seconds
- ✅ Large files (100MB - 1GB): Minutes
- ✅ Very large files (> 1GB): Constant memory usage

### System Requirements
- **Minimum RAM**: 16MB
- **Recommended RAM**: 64MB
- **Disk Space**: Input file size + output file size
- **CPU**: Any modern processor

## Changelog

### Version 1.0.0
- Initial release with all core features
- Memory-efficient chunked processing
- Complete transformation rule set
- UTF-8 safety and error resilience
- Comprehensive test suite