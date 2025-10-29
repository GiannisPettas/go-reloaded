package config

import "fmt"

// System constants for chunk processing
const (
	CHUNK_BYTES   = 4096 // 4KB chunks for memory efficiency
	OVERLAP_WORDS = 20   // Number of words to preserve between chunks
)

// Config holds system configuration
type Config struct {
	ChunkBytes   int
	OverlapWords int
}

// Validate checks if configuration values are valid
func (c Config) Validate() error {
	if c.ChunkBytes <= 0 {
		return fmt.Errorf("chunk bytes must be positive, got %d", c.ChunkBytes)
	}
	if c.OverlapWords <= 0 {
		return fmt.Errorf("overlap words must be positive, got %d", c.OverlapWords)
	}
	return nil
}