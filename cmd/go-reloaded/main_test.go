package main

import (
	"go-reloaded/internal/testutils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainWithValidArgs(t *testing.T) {
	// Create test input file
	inputContent := "hello (up) world"
	inputPath, err := testutils.CreateTestFile(inputContent)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}
	defer testutils.CleanupTestFile(inputPath)
	
	// Create output path
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "main-test-output.txt")
	defer os.Remove(outputPath)
	
	// Run main with arguments
	cmd := exec.Command("go", "run", "main.go", inputPath, outputPath)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Fatalf("Main execution failed: %v, output: %s", err, string(output))
	}
	
	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}
}

func TestMainWithInvalidArgs(t *testing.T) {
	// Test with no arguments
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	
	if err == nil {
		t.Errorf("Expected error for no arguments, but got none")
	}
	
	if !strings.Contains(string(output), "Usage:") {
		t.Errorf("Expected usage message, got: %s", string(output))
	}
}

func TestMainWithNonexistentFile(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "nonexistent.txt", "output.txt")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	
	if err == nil {
		t.Errorf("Expected error for nonexistent file, but got none")
	}
	
	if !strings.Contains(string(output), "does not exist") {
		t.Errorf("Expected file not found error, got: %s", string(output))
	}
}