package config

import "testing"

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

func TestConfigValidation(t *testing.T) {
	config := Config{
		ChunkBytes:   CHUNK_BYTES,
		OverlapWords: OVERLAP_WORDS,
	}
	
	if err := config.Validate(); err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}
}

func TestConfigValidationInvalid(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{"zero chunk bytes", Config{ChunkBytes: 0, OverlapWords: 20}},
		{"negative chunk bytes", Config{ChunkBytes: -1, OverlapWords: 20}},
		{"zero overlap words", Config{ChunkBytes: 4096, OverlapWords: 0}},
		{"negative overlap words", Config{ChunkBytes: 4096, OverlapWords: -1}},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); err == nil {
				t.Errorf("Invalid config should return error")
			}
		})
	}
}