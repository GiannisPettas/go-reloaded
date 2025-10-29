# Controller Agent Tasks — Go Reloaded

Tasks for the Controller Agent focused on workflow orchestration and main application.

---

## Task 18: Controller - Workflow Integration
**Functionality**: Orchestrate Parser → Transformer → Exporter workflow  
**TDD Steps**:
1. **Red**: Write tests for complete workflow execution
   - Test processing single chunk through complete pipeline
   - Test multi-chunk processing with overlap
   - Test error handling at each stage
2. **Green**: Implement controller that manages component interaction
   - Create `internal/controller/controller.go`
   - Implement `ProcessFile(inputPath, outputPath string) error`
   - Orchestrate Parser → Transformer → Exporter flow
   - Handle chunk overlap and context management
3. **Refactor**: Optimize workflow and error handling
4. **Validate**: End-to-end processing with simple test cases

**Dependencies**: Parser Agent (Tasks 3, 4, 16), Transformer Agent (all tasks), Exporter Agent (Task 5)

---

## Task 21: Main Application
**Functionality**: Command-line interface for the application  
**TDD Steps**:
1. **Red**: Write tests for CLI argument parsing and file handling
   - Test with valid input/output file arguments
   - Test with missing or invalid arguments
   - Test file not found scenarios
   - Test permission errors
2. **Green**: Implement main.go with proper error handling
   - Create `cmd/go-reloaded/main.go`
   - Parse command line arguments (input file, output file)
   - Call controller.ProcessFile()
   - Handle and display errors appropriately
3. **Refactor**: Improve CLI user experience and error messages
4. **Validate**: Application runs correctly from command line

**Expected Output**:
- `cmd/go-reloaded/main.go` with CLI interface
- Proper argument parsing and validation
- User-friendly error messages
- Integration with controller workflow

**Dependencies**: Task 18 (Controller workflow integration)
**Next Tasks**: Integration Agent final validation