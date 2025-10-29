package parser

import (
	"go-reloaded/internal/config"
	"go-reloaded/internal/testutils"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestReadChunkExactSize(t *testing.T) {
	// Create test file with exactly CHUNK_BYTES content
	content := strings.Repeat("a", config.CHUNK_BYTES)
	filepath, err := testutils.CreateTestFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer testutils.CleanupTestFile(filepath)
	
	data, err := ReadChunk(filepath, 0)
	if err != nil {
		t.Fatalf("ReadChunk failed: %v", err)
	}
	
	if len(data) != config.CHUNK_BYTES {
		t.Errorf("Expected %d bytes, got %d", config.CHUNK_BYTES, len(data))
	}
	
	if string(data) != content {
		t.Errorf("Content mismatch")
	}
}

func TestReadChunkLargerFile(t *testing.T) {
	// Create test file larger than CHUNK_BYTES
	content := strings.Repeat("b", config.CHUNK_BYTES*2)
	filepath, err := testutils.CreateTestFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer testutils.CleanupTestFile(filepath)
	
	data, err := ReadChunk(filepath, 0)
	if err != nil {
		t.Fatalf("ReadChunk failed: %v", err)
	}
	
	if len(data) != config.CHUNK_BYTES {
		t.Errorf("Expected %d bytes, got %d", config.CHUNK_BYTES, len(data))
	}
	
	expectedContent := strings.Repeat("b", config.CHUNK_BYTES)
	if string(data) != expectedContent {
		t.Errorf("Content mismatch for first chunk")
	}
}

func TestReadChunkSmallerFile(t *testing.T) {
	// Create test file smaller than CHUNK_BYTES
	content := "Small file content"
	filepath, err := testutils.CreateTestFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer testutils.CleanupTestFile(filepath)
	
	data, err := ReadChunk(filepath, 0)
	if err != nil {
		t.Fatalf("ReadChunk failed: %v", err)
	}
	
	if len(data) != len(content) {
		t.Errorf("Expected %d bytes, got %d", len(content), len(data))
	}
	
	if string(data) != content {
		t.Errorf("Content mismatch")
	}
}

func TestReadChunkEmptyFile(t *testing.T) {
	filepath, err := testutils.CreateTestFile("")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer testutils.CleanupTestFile(filepath)
	
	data, err := ReadChunk(filepath, 0)
	if err != nil {
		t.Fatalf("ReadChunk failed: %v", err)
	}
	
	if len(data) != 0 {
		t.Errorf("Expected 0 bytes for empty file, got %d", len(data))
	}
}

func TestReadChunkFileNotFound(t *testing.T) {
	_, err := ReadChunk("nonexistent.txt", 0)
	if err == nil {
		t.Errorf("ReadChunk should return error for nonexistent file")
	}
}

func TestReadChunkWithOffset(t *testing.T) {
	content := strings.Repeat("c", config.CHUNK_BYTES*2)
	filepath, err := testutils.CreateTestFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer testutils.CleanupTestFile(filepath)
	
	// Read second chunk
	data, err := ReadChunk(filepath, int64(config.CHUNK_BYTES))
	if err != nil {
		t.Fatalf("ReadChunk with offset failed: %v", err)
	}
	
	if len(data) != config.CHUNK_BYTES {
		t.Errorf("Expected %d bytes, got %d", config.CHUNK_BYTES, len(data))
	}
	
	expectedContent := strings.Repeat("c", config.CHUNK_BYTES)
	if string(data) != expectedContent {
		t.Errorf("Content mismatch for second chunk")
	}
}

func TestAdjustToRuneBoundary(t *testing.T) {
	// Test with complete UTF-8 characters
	completeUTF8 := []byte("Hello world")
	adjusted := AdjustToRuneBoundary(completeUTF8)
	if string(adjusted) != "Hello world" {
		t.Errorf("Complete UTF-8 should remain unchanged")
	}
}

func TestAdjustToRuneBoundaryIncomplete(t *testing.T) {
	// Test with incomplete UTF-8 character at end
	// "é" is 2 bytes: 0xC3 0xA9
	incompleteUTF8 := []byte("Hello é")
	incompleteUTF8 = append(incompleteUTF8[:len(incompleteUTF8)-1], 0xC3) // Remove last byte of é
	
	adjusted := AdjustToRuneBoundary(incompleteUTF8)
	expected := "Hello "
	if string(adjusted) != expected {
		t.Errorf("Expected %q, got %q", expected, string(adjusted))
	}
}

func TestAdjustToRuneBoundaryMultiByte(t *testing.T) {
	// Test with various multi-byte characters
	content := "Hello 世界! 🚀"
	bytes := []byte(content)
	
	// Truncate in middle of multi-byte character
	truncated := bytes[:len(bytes)-2] // Remove part of 🚀
	
	adjusted := AdjustToRuneBoundary(truncated)
	// Should end at complete character before 🚀
	expected := "Hello 世界! "
	if string(adjusted) != expected {
		t.Errorf("Expected %q, got %q", expected, string(adjusted))
	}
}

func TestReadChunkWithRuneBoundary(t *testing.T) {
	// Create content that will have UTF-8 characters at chunk boundary
	base := strings.Repeat("a", config.CHUNK_BYTES-10)
	unicode := "世界🚀" // Multi-byte characters
	content := base + unicode
	
	filepath, err := testutils.CreateTestFile(content)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer testutils.CleanupTestFile(filepath)
	
	data, err := ReadChunk(filepath, 0)
	if err != nil {
		t.Fatalf("ReadChunk failed: %v", err)
	}
	
	// Verify no UTF-8 corruption
	if !isValidUTF8(data) {
		t.Errorf("Chunk contains invalid UTF-8")
	}
	
	// Verify chunk ends at rune boundary
	if len(data) > 0 && data[len(data)-1] >= 0x80 {
		// If last byte is part of multi-byte sequence, verify it's complete
		lastRune := string(data[len(data)-1:])
		if len([]rune(lastRune)) == 0 {
			t.Errorf("Chunk ends with incomplete UTF-8 sequence")
		}
	}
}

func TestExtractOverlapWords(t *testing.T) {
	text := "word1 word2 word3 word4 word5 word6 word7 word8"
	
	overlap, remaining := ExtractOverlapWords(text)
	
	// Should extract last OVERLAP_WORDS words
	expectedOverlap := "word1 word2 word3 word4 word5 word6 word7 word8"
	expectedRemaining := ""
	
	if len(strings.Fields(text)) <= config.OVERLAP_WORDS {
		// If total words <= OVERLAP_WORDS, all words go to overlap
		if overlap != expectedOverlap || remaining != expectedRemaining {
			t.Errorf("Expected overlap=%q, remaining=%q, got overlap=%q, remaining=%q", expectedOverlap, expectedRemaining, overlap, remaining)
		}
	} else {
		// Should have exactly OVERLAP_WORDS in overlap
		overlapWords := strings.Fields(overlap)
		if len(overlapWords) != config.OVERLAP_WORDS {
			t.Errorf("Expected %d overlap words, got %d", config.OVERLAP_WORDS, len(overlapWords))
		}
	}
}

func TestPrependOverlapWords(t *testing.T) {
	overlap := "word1 word2"
	newChunk := "word3 word4 word5"
	
	result := PrependOverlapWords(overlap, newChunk)
	
	expected := "word1 word2 word3 word4 word5"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestPrependOverlapWordsEmpty(t *testing.T) {
	overlap := ""
	newChunk := "word1 word2"
	
	result := PrependOverlapWords(overlap, newChunk)
	
	if result != newChunk {
		t.Errorf("Expected %q, got %q", newChunk, result)
	}
}

func isValidUTF8(data []byte) bool {
	return utf8.Valid(data)
}