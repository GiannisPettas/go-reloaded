# Agent Workflow Instructions — Go Reloaded

This document provides specific instructions for each AI agent role in the Go Reloaded development process.

---

## General TDD Workflow for All Agents

1. **Red Phase**: Write failing tests first
2. **Green Phase**: Write minimal code to make tests pass
3. **Refactor Phase**: Improve code quality while keeping tests green
4. **Validate**: Ensure all tests pass before moving to next task

---

## Agent-Specific Instructions

### 1. Test Agent
**Role**: Foundation testing and test utilities
**Working Directory**: `internal/` and root level
**Key Responsibilities**:
- Create test helper functions and utilities
- Set up test data and file fixtures
- Write comprehensive test cases with edge scenarios
- Ensure test coverage and validation frameworks

**Instructions**:
- Always write tests before any implementation
- Use table-driven tests for multiple scenarios
- Create reusable test utilities for file operations
- Focus on edge cases and error conditions

### 2. Config Agent
**Role**: System configuration and constants
**Working Directory**: `internal/config/`
**Key Responsibilities**:
- Define CHUNK_BYTES and OVERLAP_WORDS constants
- Create configuration validation
- Set up system-wide settings

**Instructions**:
- Start with test for configuration loading
- Keep constants easily modifiable
- Add validation for configuration values
- Document all configuration options

### 3. Parser Agent
**Role**: File reading and chunk management
**Working Directory**: `internal/parser/`
**Key Responsibilities**:
- Implement chunked file reading
- Handle UTF-8 rune boundary alignment
- Manage chunk overlap with word preservation

**Instructions**:
- Test with various file sizes and Unicode content
- Ensure no data corruption at chunk boundaries
- Implement streaming approach for memory efficiency
- Handle file reading errors gracefully

### 4. Transformer Agent
**Role**: Text transformation and FSM implementation
**Working Directory**: `internal/transformer/`
**Key Responsibilities**:
- Implement dual FSM architecture (character parser + token processor)
- Create all transformation rules (hex, bin, case, punctuation, articles)
- Handle single-pass processing with fixed-size buffers
- Manage context across chunks with pending command state

**Instructions**:
- Implement one transformation rule at a time
- Test each rule with golden test cases
- Handle invalid commands gracefully
- Maintain state for cross-chunk operations

### 5. Exporter Agent
**Role**: Output file writing
**Working Directory**: `internal/exporter/`
**Key Responsibilities**:
- Write processed text to output files
- Handle progressive writing for memory efficiency
- Manage file creation and error handling

**Instructions**:
- Test with various output scenarios
- Implement streaming write operations
- Handle file system errors appropriately
- Ensure output formatting is correct

### 6. Controller Agent
**Role**: Workflow orchestration and main application
**Working Directory**: `internal/controller/` and `cmd/go-reloaded/`
**Key Responsibilities**:
- Integrate all components (Parser → Transformer → Exporter)
- Handle command-line interface
- Manage overall application flow
- Implement error handling and logging

**Instructions**:
- Test complete workflow integration
- Handle component failures gracefully
- Implement proper CLI argument parsing
- Ensure clean resource management

### 7. Integration Agent
**Role**: End-to-end testing and validation
**Working Directory**: Root level and test directories
**Key Responsibilities**:
- Execute all golden tests (T1-T27)
- Performance and memory validation
- System integration testing
- Final quality assurance

**Instructions**:
- Run all golden tests and ensure they pass
- Test with large files for memory efficiency
- Validate system performance requirements
- Create comprehensive integration test suite

---

## Agent Execution Order

1. **Test Agent** → Set up testing framework
2. **Config Agent** → Define system constants
3. **Parser Agent** → Implement file reading
4. **Exporter Agent** → Implement file writing
5. **Transformer Agent** → Implement transformations (longest phase)
6. **Controller Agent** → Integrate components
7. **Integration Agent** → Final validation

---

## Inter-Agent Dependencies

- **Config Agent** outputs are used by all other agents
- **Parser Agent** outputs feed into **Transformer Agent**
- **Transformer Agent** outputs feed into **Exporter Agent**
- **Controller Agent** orchestrates Parser → Transformer → Exporter
- **Integration Agent** validates work from all previous agents

---

## Success Criteria for Each Agent

- All tests pass (100% success rate)
- Code follows Go best practices
- Memory usage remains efficient
- Error handling is comprehensive
- Documentation is clear and complete