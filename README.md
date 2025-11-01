# Go-Reloaded

A high-performance text processing tool that transforms text files using various commands. Built with a dual finite state machine architecture for maximum efficiency and minimal memory usage.

## Features

- **Numeric Base Conversion**: Convert hexadecimal and binary numbers to decimal
- **Case Transformations**: Change text to uppercase, lowercase, or capitalize
- **Article Correction**: Automatically fix "a/an" usage based on vowel sounds  
- **Punctuation Spacing**: Fix spacing around punctuation marks
- **Command Chaining**: Apply multiple transformations to the same word
- **Memory Efficient**: Processes files of any size using only ~8KB of memory
- **Zero Dependencies**: Uses only Go standard library, no external packages

## Installation

### Prerequisites
- Go 1.19 or higher

### Build
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

```bash
./go-reloaded input.txt output.txt
```

- `input.txt`: Path to the input text file
- `output.txt`: Path where the processed output will be saved

## Commands

### Numeric Conversions

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

### Case Transformations

#### Single Word
```
Input:  "make this word (up)"
Output: "make this WORD"

Input:  "make this word (low)"  
Output: "make this word"

Input:  "make this word (cap)"
Output: "make this Word"
```
#### Multiple Words
```
Input:  "These three words (up, 3) should be uppercase"
Output: "THESE THREE WORDS should be uppercase"

Input:  "These two words (cap, 2) should be capitalized"  
Output: "These Two words should be capitalized"
```

### Article Corrections
```
Input:  "I need a apple and an car"
Output: "I need an apple and a car"

Input:  "It was a honor to meet an European"
Output: "It was an honor to meet a European"
```

### Punctuation Spacing
```
Input:  "Hello , world ! How are you ?"
Output: "Hello, world! How are you?"
```

### Command Chaining
```
Input:  "The number 1010 (bin) (hex) is interesting"
Output: "The number 16 is interesting"
```
*Explanation: 1010 (binary) → 10 (decimal) → 16 (hexadecimal)*

## Examples

### Sample Input (`sample.txt`)
```
Simply add 1010 (bin) (hex) , and check the total !

I need a apple and an car for the trip.

Convert FF (hex) to decimal and make it (up) .
```

### Expected Output (`result.txt`)
```
Simply add 16, and check the total!

I need an apple and a car for the trip.

Convert 255 to decimal and make it UP.
```

### Running the Example
```bash
./go-reloaded sample.txt result.txt
```

## Performance

### File Size Handling
- ✅ Small files (< 4KB): Single-pass processing
- ✅ Medium files (4KB - 100MB): Chunked processing with overlap
- ✅ Large files (100MB+): Constant memory usage (~8KB)
- ✅ Very large files (1GB+): No memory limitations

### System Requirements
- **RAM**: 16MB minimum, 64MB recommended
- **CPU**: Any modern processor
- **Disk**: Input file size + output file size

## Architecture Highlights

- **Dual FSM Design**: Two finite state machines working in tandem
- **Single-Pass Processing**: No multiple iterations over data
- **Memory Efficient**: Fixed-size buffers, constant memory usage
- **UTF-8 Safe**: Handles international characters without corruption
- **Chunked Processing**: Smart overlap handling for large files
- **Zero Dependencies**: Pure Go standard library implementation

## Testing

### Run All Tests (Recommended)
```bash
cd internal/testutils && go test -v -run TestAllProject
```
This single command runs all tests including the 27 golden test cases with nice formatting.

### Alternative Test Commands
```bash
# Run all tests manually
go test -count=1 ./...

# Run golden test suite only
cd internal/testutils && go test -v -run TestGoldenCases

# Run specific package tests
go test ./internal/transformer/
go test ./internal/config/
```

The project includes 27 comprehensive test cases covering all transformation scenarios.

## Project Structure

```
go-reloaded/
├── cmd/go-reloaded/          # CLI application entry point
├── internal/
│   ├── config/               # System configuration constants
│   ├── parser/               # File reading and chunking
│   ├── transformer/          # Dual-FSM text transformation engine
│   ├── exporter/             # File writing operations
│   ├── controller/           # Workflow orchestration
│   └── testutils/            # Testing utilities and golden tests
├── docs/                     # Technical documentation
└── README.md                 # This file
```

## Technical Details

For detailed technical documentation including FSM architecture, algorithms, and implementation details, see [`docs/technical_architecture.md`](docs/technical_architecture.md).

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass: `go test -count=1 ./...`
5. Submit a pull request

## Support

For technical issues, please create an issue in the repository with:
- Input text that causes the problem
- Expected vs actual output
- Go version and operating system