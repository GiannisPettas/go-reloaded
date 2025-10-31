# Go-Reloaded Project Structure

## Directory Organization

```
go-reloaded/
├── cmd/go-reloaded/          # CLI application entry point
├── internal/                 # Private application packages
│   ├── config/              # System configuration constants
│   ├── parser/              # File reading and chunking logic
│   ├── transformer/         # Dual-FSM text transformation engine
│   ├── exporter/            # File writing operations
│   ├── controller/          # Workflow orchestration
│   └── testutils/           # Testing utilities and golden tests
├── docs/                    # Technical documentation and web assets
├── reports/                 # Agent analysis reports
└── README.md               # Project documentation
```

## Core Components

### Entry Point
- **cmd/go-reloaded/main.go**: CLI interface handling command-line arguments and error reporting

### Internal Packages
- **config/**: System-wide constants and configuration values
- **parser/**: File reading with chunked processing for memory efficiency
- **transformer/**: Dual finite state machine implementation for text transformations
- **exporter/**: File writing operations with proper error handling
- **controller/**: Orchestrates the entire processing pipeline
- **testutils/**: Comprehensive testing framework with golden test suite

### Documentation
- **docs/**: Contains technical architecture documentation and web presentation
- **reports/**: Agent-generated analysis reports for each component

## Architectural Patterns

### Layered Architecture
- **Presentation Layer**: CLI interface (cmd/)
- **Business Logic Layer**: Transformation engine (internal/transformer/)
- **Data Access Layer**: File I/O operations (internal/parser/, internal/exporter/)
- **Orchestration Layer**: Workflow coordination (internal/controller/)

### Separation of Concerns
- Each internal package has a single, well-defined responsibility
- Clear interfaces between components
- Minimal coupling between layers

### Modular Design
- Independent packages that can be tested in isolation
- Reusable components with clear APIs
- Extensible architecture for adding new transformations

## Component Relationships
- **Controller** orchestrates the entire pipeline
- **Parser** feeds data to **Transformer**
- **Transformer** processes text using dual FSM architecture
- **Exporter** writes processed results
- **Config** provides system-wide constants to all components