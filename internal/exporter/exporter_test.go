package exporter

import (
	"go-reloaded/internal/testutils"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Purpose: Tests constants during development/CI

func TestWriteChunkNewFile(t *testing.T) {
	content := "Hello, World!\nThis is test content."

	// Create temporary output path
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "test-output.txt")
	defer os.Remove(outputPath)

	err := WriteChunk(outputPath, content)
	if err != nil {
		t.Fatalf("WriteChunk failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outputPath)
	}

	// Verify content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(data) != content {
		t.Errorf("Content mismatch. Expected: %q, Got: %q", content, string(data))
	}
}

func TestWriteChunkEmptyContent(t *testing.T) {
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "test-empty.txt")
	defer os.Remove(outputPath)

	err := WriteChunk(outputPath, "")
	if err != nil {
		t.Fatalf("WriteChunk with empty content failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read empty output file: %v", err)
	}

	if len(data) != 0 {
		t.Errorf("Empty file should have zero length, got %d", len(data))
	}
}

func TestAppendChunkExistingFile(t *testing.T) {
	initialContent := "Initial content\n"
	appendContent := "Appended content\n"

	// Create initial file
	filepath, err := testutils.CreateTestFile(initialContent)
	if err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}
	defer testutils.CleanupTestFile(filepath)

	err = AppendChunk(filepath, appendContent)
	if err != nil {
		t.Fatalf("AppendChunk failed: %v", err)
	}

	// Verify combined content
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read appended file: %v", err)
	}

	expected := initialContent + appendContent
	if string(data) != expected {
		t.Errorf("Content mismatch. Expected: %q, Got: %q", expected, string(data))
	}
}

func TestAppendChunkNewFile(t *testing.T) {
	content := "New file content"

	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "test-append-new.txt")
	defer os.Remove(outputPath)

	err := AppendChunk(outputPath, content)
	if err != nil {
		t.Fatalf("AppendChunk to new file failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read new appended file: %v", err)
	}

	if string(data) != content {
		t.Errorf("Content mismatch. Expected: %q, Got: %q", content, string(data))
	}
}

func TestWriteChunkUnicodeContent(t *testing.T) {
	content := "Hello 世界! 🚀 Café"

	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "test-unicode.txt")
	defer os.Remove(outputPath)

	err := WriteChunk(outputPath, content)
	if err != nil {
		t.Fatalf("WriteChunk with Unicode failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read Unicode file: %v", err)
	}

	if string(data) != content {
		t.Errorf("Unicode content mismatch. Expected: %q, Got: %q", content, string(data))
	}
}

func TestWriteChunkInvalidPath(t *testing.T) {
	var invalidPath string

	// Use OS-specific invalid paths
	if runtime.GOOS == "windows" {
		// Windows invalid characters: < > : " | ? *
		invalidPath = "<>:\"|?*invalid.txt"
	} else {
		// Unix/Linux: null byte is invalid
		invalidPath = "/tmp/invalid\x00file.txt"
	}

	err := WriteChunk(invalidPath, "content")
	if err == nil {
		t.Errorf("WriteChunk should return error for invalid path: %s", invalidPath)
	}
}

func TestAppendChunkMultiple(t *testing.T) {
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "test-multiple-append.txt")
	defer os.Remove(outputPath)

	chunks := []string{"Chunk 1\n", "Chunk 2\n", "Chunk 3\n"}

	for _, chunk := range chunks {
		err := AppendChunk(outputPath, chunk)
		if err != nil {
			t.Fatalf("AppendChunk failed for chunk %q: %v", chunk, err)
		}
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read multiple append file: %v", err)
	}

	expected := "Chunk 1\nChunk 2\nChunk 3\n"
	if string(data) != expected {
		t.Errorf("Multiple append content mismatch. Expected: %q, Got: %q", expected, string(data))
	}
}
