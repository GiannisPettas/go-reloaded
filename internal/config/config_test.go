package config

import "testing"

// Purpose: Tests constants during development/CI

func TestChunkBytesConstant(t *testing.T) {
	if CHUNK_BYTES <= 0 {
		t.Errorf("CHUNK_BYTES must be positive, got %d", CHUNK_BYTES)
	}
	if CHUNK_BYTES < 1024 || CHUNK_BYTES > 8192 {
		t.Errorf("CHUNK_BYTES should be reasonable size (1024-8192), got %d", CHUNK_BYTES)
	}
}

func TestOverlapWordsConstant(t *testing.T) {
	if OVERLAP_WORDS <= 0 {
		t.Errorf("OVERLAP_WORDS must be positive, got %d", OVERLAP_WORDS)
	}
	if OVERLAP_WORDS > 50 {
		t.Errorf("OVERLAP_WORDS should be reasonable (<=50), got %d", OVERLAP_WORDS)
	}
}

func TestValidateConstants(t *testing.T) {
	if err := ValidateConstants(); err != nil {
		t.Errorf("ValidateConstants should not return error with current constants: %v", err)
	}
}
