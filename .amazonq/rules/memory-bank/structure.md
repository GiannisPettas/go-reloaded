# Go-Reloaded Project Structure

## Directory Organization

### Root Level
```
go-reloaded/
├── cmd/go-reloaded/          # CLI application entry point
├── internal/                 # Core application packages
├── docs/                     # Documentation and examples
├── .amazonq/rules/          # Development rules and guidelines
├── go.mod                   # Go module definition
├── README.md               # Project documentation
└── run_*_tests.*           # Test execution scripts
```

### Core Components (`internal/`)

#### `internal/config/`
- **Purpose**: System configuration constants and settings
- **Key Files**: `config.go`, `config_test.go`
- **Responsibilities**: Buffer sizes, chunk limits, system constants

#### `internal/parser/`
- **Purpose**: File reading and chunking operations
- **Key Files**: `parser.go`, `parser_test.go`
- **Responsibilities**: File I/O, chunk management, overlap handling

#### `internal/transformer/`
- **Purpose**: Dual-FSM text transformation engine
- **Key Files**: `transformer.go`, `transformer_test.go`
- **Responsibilities**: Token processing, command execution, text transformations

#### `internal/exporter/`
- **Purpose**: File writing operations
- **Key Files**: `exporter.go`, `exporter_test.go`
- **Responsibilities**: Output file creation, write operations

#### `internal/controller/`
- **Purpose**: Workflow orchestration
- **Key Files**: `controller.go`, `controller_test.go`
- **Responsibilities**: Component coordination, error handling, process flow

#### `internal/testutils/`
- **Purpose**: Testing utilities and golden tests
- **Key Files**: `golden_test.go`, `golden.go`, `testutils.go`
- **Responsibilities**: Test infrastructure, golden test cases, test utilities

### Documentation (`docs/`)

#### Core Documentation
- `technical_architecture.md` - Complete technical overview
- `golden_tests.md` - All 29 test cases with examples
- Component-specific explanations (`*_explained.md`)

#### Interactive Demo
- `index.html`, `script.js`, `styles.css` - Web-based demo
- `samples/` - Example input files for testing

#### Archive
- `archive/` - Historical development documentation
- Agent reports and workflow documentation

### CLI Application (`cmd/go-reloaded/`)
- `main.go` - Application entry point
- `main_test.go` - CLI integration tests
- Example files for testing

## Architectural Patterns

### Dual Finite State Machine Architecture
- **Low-Level FSM**: Character-by-character parsing and tokenization
- **High-Level FSM**: Token processing and transformation application
- **Conveyor Belt System**: Fixed-size token buffer for memory efficiency

### Component Relationships
```
CLI (main.go) 
    ↓
Controller (orchestration)
    ↓
Parser → Transformer → Exporter
    ↑         ↓
Config ←→ TestUtils
```

### Data Flow
1. **Input**: Parser reads file in chunks
2. **Processing**: Transformer applies dual-FSM processing
3. **Output**: Exporter writes transformed text
4. **Coordination**: Controller manages the entire workflow

### Memory Management
- Fixed-size buffers prevent memory growth
- Chunked processing for large files
- Overlap handling maintains transformation context
- Constant memory usage regardless of file size