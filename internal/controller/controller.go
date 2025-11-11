package controller

import (
	"fmt"
	"go-reloaded/internal/config"
	"go-reloaded/internal/exporter"
	"go-reloaded/internal/parser"
	"go-reloaded/internal/transformer"
	"os"
	"strings"
)

// ProcessFile orchestrates the complete workflow: Parser → Transformer → Exporter
func ProcessFile(inputPath, outputPath string) error {
	// Check if input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}
	// Get file size to determine if we need chunked processing
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// For small files, process in one chunk
	if fileInfo.Size() <= int64(config.CHUNK_BYTES) {
		return processSingleChunk(inputPath, outputPath)
	}

	// For larger files, use chunked processing with overlap
	return processChunkedFile(inputPath, outputPath)
}

// processSingleChunk handles files that fit in a single chunk
func processSingleChunk(inputPath, outputPath string) error {
	// Read entire file
	data, err := parser.ReadChunk(inputPath, 0)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Convert to text
	text := string(data)

	// Apply transformations in single pass
	result := transformer.ProcessText(text)

	// Write to output
	err = exporter.WriteChunk(outputPath, result)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

// processChunkedFile handles large files with proper chunked processing
func processChunkedFile(inputPath, outputPath string) error {
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	var offset int64 = 0
	var overlapContext string
	isFirstChunk := true

	for {
		// Read chunk
		data, err := parser.ReadChunk(inputPath, offset)
		if err != nil {
			return fmt.Errorf("failed to read chunk at offset %d: %w", offset, err)
		}

		// If no data, we're done
		if len(data) == 0 {
			break
		}

		// Convert to text
		chunkText := string(data)

		// Merge with overlap context
		var textToProcess string
		if overlapContext != "" {
			textToProcess = overlapContext + chunkText
		} else {
			textToProcess = chunkText
		}

		// Apply single-pass FSM transformation to this chunk
		processedChunk := transformer.ProcessText(textToProcess)

		// If we had overlap context, remove it from the processed result to avoid duplication
		if overlapContext != "" {
			// Skip the overlap words from the processed result
			overlapWordCount := len(strings.Fields(overlapContext))
			processedWords := strings.Fields(processedChunk)
			if len(processedWords) > overlapWordCount {
				processedChunk = strings.Join(processedWords[overlapWordCount:], " ")
			} else {
				processedChunk = "" // All words were overlap
			}
		}

		// Extract overlap for next chunk and get remaining text
		newOverlap, remaining := parser.ExtractOverlapWords(processedChunk)

		// Write remaining text to output
		if remaining != "" {
			if isFirstChunk {
				err = exporter.WriteChunk(outputPath, remaining)
				isFirstChunk = false
			} else {
				err = exporter.AppendChunk(outputPath, remaining)
			}
			if err != nil {
				return fmt.Errorf("failed to write chunk: %w", err)
			}
		}

		// Update context and offset
		overlapContext = newOverlap
		offset += int64(len(data))

		// Safety check to prevent infinite loops
		if offset >= fileInfo.Size() {
			break
		}

		// If chunk was smaller than expected, we're at end of file
		if len(data) < config.CHUNK_BYTES {
			break
		}
	}

	// Write any remaining overlap context at the end
	if overlapContext != "" {
		if isFirstChunk {
			err = exporter.WriteChunk(outputPath, overlapContext)
		} else {
			err = exporter.AppendChunk(outputPath, overlapContext)
		}
		if err != nil {
			return fmt.Errorf("failed to write final overlap: %w", err)
		}
	}

	return nil
}
