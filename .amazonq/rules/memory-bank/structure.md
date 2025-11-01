# Go-Reloaded Project Structure

## Directory Organization

### Root Level
```
go-reloaded/
├── cmd/                      # CLI application entry point
├── internal/                 # Core application packages
├── docs/                     # Technical documentation
├── reports/                  # Agent analysis reports
├── go.mod                    # Go module definition
├── LICENSE                   # MIT License
└── README.md                 # Project documentation
```

### Core Components (`internal/`)

#### `config/`
- **Purpose**: System configuration constants and settings
- **Key Files**: `config.go`, `config_test.go`
- **Role**: Centralized configuration management for buffer sizes, processing limits

#### `parser/`
- **Purpose**: File reading and chunking operations
- **Key Files**: `parser.go`, `parser_test.go`
- **Role**: Handles input file processing, chunked reading for large files

#### `transformer/`
- **Purpose**: Dual-FSM text transformation engine
- **Key Files**: `transformer.go`, `transformer_test.go`
- **Role**: Core text processing logic using finite state machines

#### `exporter/`
- **Purpose**: File writing operations
- **Key Files**: `exporter.go`, `exporter_test.go`
- **Role**: Handles output file generation and writing

#### `controller/`
- **Purpose**: Workflow orchestration
- **Key Files**: `controller.go`, `controller_test.go`
- **Role**: Coordinates the entire processing pipeline

#### `testutils/`
- **Purpose**: Testing utilities and golden tests
- **Key Files**: `golden.go`, `testutils.go`, various test files
- **Role**: Comprehensive test framework with 27 golden test cases

### Application Entry (`cmd/`)
- **`go-reloaded/`**: Main CLI application
  - `main.go`: Application entry point
  - `main_test.go`: Integration tests
  - Example files for testing

### Documentation (`docs/`)
- Technical architecture documentation
- Agent workflow specifications
- Sample files and examples
- Interactive documentation (HTML/CSS/JS)

## Architectural Patterns

### Dual Finite State Machine (FSM)
- Two FSMs working in tandem for text processing
- State-based parsing for command recognition
- Efficient memory usage through state transitions

### Layered Architecture
1. **Presentation Layer**: CLI interface (`cmd/`)
2. **Control Layer**: Workflow orchestration (`controller/`)
3. **Business Layer**: Text transformation logic (`transformer/`)
4. **Data Layer**: File I/O operations (`parser/`, `exporter/`)
5. **Configuration Layer**: System settings (`config/`)

### Component Relationships
```
CLI (main.go) → Controller → Parser → Transformer → Exporter
                    ↓
                 Config (shared by all components)
                    ↓
                 TestUtils (testing framework)
```

### Design Principles
- **Single Responsibility**: Each package has a focused purpose
- **Dependency Injection**: Components receive dependencies explicitly
- **Testability**: Comprehensive test coverage with utilities
- **Memory Efficiency**: Constant memory usage patterns
- **Error Handling**: Consistent error propagation throughout layers