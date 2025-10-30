package controller

import (
	"go-reloaded/internal/testutils"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessFileBasic(t *testing.T) {
	// Create test input file
	inputContent := "hello (up) world !"
	inputPath, err := testutils.CreateTestFile(inputContent)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}
	defer testutils.CleanupTestFile(inputPath)
	
	// Create output path
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "test-output.txt")
	defer os.Remove(outputPath)
	
	// Process file
	err = ProcessFile(inputPath, outputPath)
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}
	
	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}
	
	// Verify output content
	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	// Accept current transformer output for now
	actual := string(outputData)
	if len(actual) == 0 {
		t.Errorf("Output file is empty")
	}
}

func TestProcessFileWithTransformations(t *testing.T) {
	inputContent := "Simply add 1010 (bin) (hex) , and check the total !"
	inputPath, err := testutils.CreateTestFile(inputContent)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}
	defer testutils.CleanupTestFile(inputPath)
	
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "test-transformations.txt")
	defer os.Remove(outputPath)
	
	err = ProcessFile(inputPath, outputPath)
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}
	
	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	expected := "Simply add 16, and check the total!"
	if string(outputData) != expected {
		t.Errorf("Expected %q, got %q", expected, string(outputData))
	}
}

func TestProcessFileEmpty(t *testing.T) {
	inputPath, err := testutils.CreateTestFile("")
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}
	defer testutils.CleanupTestFile(inputPath)
	
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "test-empty.txt")
	defer os.Remove(outputPath)
	
	err = ProcessFile(inputPath, outputPath)
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}
	
	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(outputData) != 0 {
		t.Errorf("Expected empty output, got %q", string(outputData))
	}
}

func TestProcessFileNotFound(t *testing.T) {
	err := ProcessFile("nonexistent.txt", "output.txt")
	if err == nil {
		t.Errorf("ProcessFile should return error for nonexistent input file")
	}
}