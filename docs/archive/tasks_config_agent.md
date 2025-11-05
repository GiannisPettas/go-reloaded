# Config Agent Tasks — Go Reloaded

Tasks for the Config Agent focused on system configuration and constants.

---

## Task 1: Project Setup and Configuration
**Functionality**: Define system constants and basic project structure  
**TDD Steps**:
1. **Red**: Write tests for configuration validation
   - Test CHUNK_BYTES constant is positive and reasonable (e.g., 1024-8192 bytes)
   - Test OVERLAP_WORDS constant is positive and less than expected words per chunk
   - Test configuration struct validation
2. **Green**: Implement config package with CHUNK_BYTES, OVERLAP_WORDS constants
   - Create `internal/config/config.go`
   - Define `const CHUNK_BYTES = 4096`
   - Define `const OVERLAP_WORDS = 20`
   - Create `Config` struct with validation methods
3. **Refactor**: Add configuration validation and error handling
4. **Validate**: Configuration loads correctly and validation works

**Expected Output**:
- `internal/config/config.go` with system constants
- Configuration validation functions
- Test cases for all configuration scenarios

**Dependencies**: None
**Next Tasks**: All other agents depend on these constants