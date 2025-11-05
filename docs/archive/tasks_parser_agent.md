# Parser Agent Tasks — Go Reloaded

Tasks for the Parser Agent focused on file reading and chunk management.

---

## Task 3: Parser - Basic File Reading
**Functionality**: Read file in fixed-size chunks  
**TDD Steps**:
1. **Red**: Write tests for reading CHUNK_BYTES from file
   - Test reading exact CHUNK_BYTES from file larger than chunk size
   - Test reading from file smaller than CHUNK_BYTES
   - Test reading from empty file
   - Test file not found error handling
2. **Green**: Implement basic file reader with chunk size limits
   - Create `internal/parser/parser.go`
   - Implement `ReadChunk(filepath string, offset int64) ([]byte, error)`
   - Use config.CHUNK_BYTES for chunk size
3. **Refactor**: Improve error handling and edge cases
4. **Validate**: Correct byte reading from test files

**Dependencies**: Config Agent (Task 1)

---

## Task 4: Parser - UTF-8 Rune Boundary Alignment
**Functionality**: Ensure chunks end at complete UTF-8 runes  
**TDD Steps**:
1. **Red**: Write tests with multi-byte Unicode characters at chunk boundaries
   - Test with UTF-8 characters (é, 中, 🚀) at chunk boundaries
   - Test chunk adjustment to complete runes
   - Test various Unicode scenarios
2. **Green**: Implement rune boundary detection and adjustment
   - Add `AdjustToRuneBoundary(data []byte) []byte`
   - Modify ReadChunk to ensure rune boundary alignment
   - Handle incomplete UTF-8 sequences at chunk end
3. **Refactor**: Optimize rune boundary detection
4. **Validate**: No Unicode characters are corrupted

**Dependencies**: Task 3

---

## Task 16: Parser - Chunk Overlap Implementation
**Functionality**: Implement OVERLAP_WORDS context preservation  
**TDD Steps**:
1. **Red**: Write tests for word overlap between chunks
   - Test extracting last N words from processed chunk
   - Test prepending overlap words to next chunk
   - Test word boundary detection
2. **Green**: Implement word extraction and storage from chunk end
   - Add `ExtractOverlapWords(text string) (overlap, remaining string)`
   - Add `PrependOverlapWords(overlap, newChunk string) string`
   - Use config.OVERLAP_WORDS for word count
3. **Refactor**: Optimize word extraction and merging
4. **Validate**: Context preservation across chunk boundaries

**Dependencies**: Task 4, Transformer Agent basic functionality