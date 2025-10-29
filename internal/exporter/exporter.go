package exporter

import (
	"fmt"
	"os"
	filepath "path/filepath"
)

// WriteChunk writes content to a file, creating it if it doesn't exist
func WriteChunk(filePath, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	// Write content to file (overwrites if exists)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	
	return nil
}

// AppendChunk appends content to a file, creating it if it doesn't exist
func AppendChunk(filePath, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	// Open file for appending (create if doesn't exist)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s for appending: %w", filePath, err)
	}
	defer file.Close()
	
	// Write content
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to append to file %s: %w", filePath, err)
	}
	
	return nil
}