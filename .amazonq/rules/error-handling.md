# Go-Reloaded Error Handling Rules

## Error Wrapping Standards
- Always use `fmt.Errorf` with `%w` verb for error chaining
- Include operation context in error messages
- Preserve original error information through the chain

## Error Message Format
```go
return fmt.Errorf("failed to read chunk at offset %d: %w", offset, err)
return fmt.Errorf("invalid command '%s' in transformer: %w", cmd, err)
return fmt.Errorf("unable to write output file '%s': %w", filepath, err)
```

## Early Return Pattern
- Check errors immediately after operations
- Return early, don't continue processing on error
- Avoid nested error handling

## Context Information
Include relevant details in error messages:
- File paths for file operations
- Byte offsets for parsing operations  
- Command names for transformation operations
- Token positions for FSM operations

## Error Propagation
- Controller layer: Orchestrate error handling across components
- Parser layer: Include file path and read position
- Transformer layer: Include command and token context
- Exporter layer: Include output file path and write position

## Testing Error Cases
- Test both success and failure scenarios
- Verify error message content and wrapped errors
- Use table-driven tests for multiple error conditions