package config

import "fmt"

// System constants for chunk processing
const (
	CHUNK_BYTES   = 4096 // 4KB chunks for memory efficiency - can go from 1kb to 8kb
	OVERLAP_WORDS = 20   // Number of words to preserve between chunks - can go from 10 to 20
	// Also determines token buffer size (4x OVERLAP_WORDS = 80 tokens)
)

// ValidateConstants checks if all constants are within valid ranges
func ValidateConstants() error {
	if CHUNK_BYTES <= 0 {
		return fmt.Errorf("CHUNK_BYTES must be positive, got %d", CHUNK_BYTES)
	}
	if CHUNK_BYTES < 1024 {
		return fmt.Errorf("CHUNK_BYTES must be at least 1024 bytes, got %d", CHUNK_BYTES)
	}
	if CHUNK_BYTES > 8192 {
		return fmt.Errorf("CHUNK_BYTES must be at most 8192 bytes, got %d", CHUNK_BYTES)
	}
	if OVERLAP_WORDS <= 0 {
		return fmt.Errorf("OVERLAP_WORDS must be positive, got %d", OVERLAP_WORDS)
	}
	if OVERLAP_WORDS < 10 {
		return fmt.Errorf("OVERLAP_WORDS too small (min 10), got %d", OVERLAP_WORDS)
	}
	if OVERLAP_WORDS > 20 {
		return fmt.Errorf("OVERLAP_WORDS too large (max 20), got %d", OVERLAP_WORDS)
	}
	return nil
}
