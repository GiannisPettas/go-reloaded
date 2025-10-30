package parser

import (
	"fmt"
	"go-reloaded/internal/config"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// ReadChunk reads a chunk of data from file starting at the given offset
func ReadChunk(filepath string, offset int64) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()
	
	// Seek to offset
	if offset > 0 {
		_, err = file.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("failed to seek to offset %d: %w", offset, err)
		}
	}
	
	// Read up to CHUNK_BYTES
	buffer := make([]byte, config.CHUNK_BYTES)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read from file: %w", err)
	}
	
	// Return only the bytes that were actually read
	chunk := buffer[:n]
	
	// Adjust to rune boundary to avoid UTF-8 corruption
	adjusted := AdjustToRuneBoundary(chunk)
	return adjusted, nil
}

// AdjustToRuneBoundary ensures the byte slice ends at a complete UTF-8 rune
func AdjustToRuneBoundary(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	
	// Check if data is already valid UTF-8
	if utf8.Valid(data) {
		return data
	}
	
	// Find the last valid rune boundary by working backwards
	for i := len(data) - 1; i >= 0; i-- {
		if utf8.Valid(data[:i+1]) {
			return data[:i+1]
		}
	}
	
	// If no valid UTF-8 found, return empty
	return []byte{}
}

// ExtractOverlapWords extracts the last OVERLAP_WORDS from processed text
// Returns (overlap, remaining) where overlap contains the last words
func ExtractOverlapWords(text string) (overlap, remaining string) {
	words := strings.Fields(text)
	
	if len(words) <= config.OVERLAP_WORDS {
		// If we have fewer words than overlap size, return all as overlap
		return text, ""
	}
	
	// Split into remaining and overlap
	remainingWords := words[:len(words)-config.OVERLAP_WORDS]
	overlapWords := words[len(words)-config.OVERLAP_WORDS:]
	
	remaining = strings.Join(remainingWords, " ")
	overlap = strings.Join(overlapWords, " ")
	
	return overlap, remaining
}

// PrependOverlapWords prepends overlap words to new chunk text
func PrependOverlapWords(overlap, newChunk string) string {
	if overlap == "" {
		return newChunk
	}
	if newChunk == "" {
		return overlap
	}
	return overlap + " " + newChunk
}