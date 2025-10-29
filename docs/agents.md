# AI Agent Roles — Go Reloaded

This document defines the specialized AI agent roles for developing the Go Reloaded text processing system using Test Driven Development (TDD).

---

## Agent Roles

### 1. **Test Agent**
- **Responsibility**: Write comprehensive unit tests before implementation
- **Focus**: Golden test validation, edge cases, error scenarios
- **Output**: Test files with clear assertions and test data

### 2. **Config Agent** 
- **Responsibility**: Define system constants and configuration
- **Focus**: CHUNK_BYTES, OVERLAP_WORDS, file paths
- **Output**: Configuration structs and validation

### 3. **Parser Agent**
- **Responsibility**: File reading with chunk management
- **Focus**: UTF-8 rune boundary alignment, chunk overlap
- **Output**: Streaming file parser with context preservation

### 4. **Transformer Agent**
- **Responsibility**: FSM implementation for text transformations
- **Focus**: Rule application, state management, context awareness
- **Output**: Finite State Machine with all transformation rules

### 5. **Exporter Agent**
- **Responsibility**: Progressive output writing
- **Focus**: Memory-efficient file writing, result formatting
- **Output**: Streaming file exporter

### 6. **Controller Agent**
- **Responsibility**: Orchestrate the complete workflow
- **Focus**: Component integration, error handling, main execution
- **Output**: Main application controller

### 7. **Integration Agent**
- **Responsibility**: End-to-end testing and validation
- **Focus**: Golden test execution, performance validation
- **Output**: Integration tests and benchmarks

---

## Development Workflow

Each agent follows TDD principles:
1. **Red**: Write failing tests first
2. **Green**: Implement minimal code to pass tests
3. **Refactor**: Improve code while keeping tests green

Agents work incrementally, building upon previous agent outputs to ensure seamless integration.