package testutils

import (
	"os"
	"testing"
)

func TestCreateTestFile(t *testing.T) {
	content := "Hello, World!\nThis is a test file."
	
	filepath, err := CreateTestFile(content)
	if err != nil {
		t.Fatalf("CreateTestFile failed: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Errorf("Test file was not created: %s", filepath)
	}
	
	// Verify content
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	
	if string(data) != content {
		t.Errorf("File content mismatch. Expected: %q, Got: %q", content, string(data))
	}
	
	// Cleanup
	CleanupTestFile(filepath)
}

func TestCleanupTestFile(t *testing.T) {
	content := "Temporary test content"
	filepath, err := CreateTestFile(content)
	if err != nil {
		t.Fatalf("CreateTestFile failed: %v", err)
	}
	
	// Verify file exists before cleanup
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Errorf("Test file should exist before cleanup: %s", filepath)
	}
	
	// Cleanup
	err = CleanupTestFile(filepath)
	if err != nil {
		t.Errorf("CleanupTestFile failed: %v", err)
	}
	
	// Verify file is removed
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		t.Errorf("Test file should be removed after cleanup: %s", filepath)
	}
}

func TestCompareFiles(t *testing.T) {
	content1 := "Same content"
	content2 := "Same content"
	content3 := "Different content"
	
	file1, _ := CreateTestFile(content1)
	file2, _ := CreateTestFile(content2)
	file3, _ := CreateTestFile(content3)
	
	defer CleanupTestFile(file1)
	defer CleanupTestFile(file2)
	defer CleanupTestFile(file3)
	
	// Test identical files
	if err := CompareFiles(file1, file2); err != nil {
		t.Errorf("Identical files should not return error: %v", err)
	}
	
	// Test different files
	if err := CompareFiles(file1, file3); err == nil {
		t.Errorf("Different files should return error")
	}
}

func TestCompareFilesNonExistent(t *testing.T) {
	err := CompareFiles("nonexistent1.txt", "nonexistent2.txt")
	if err == nil {
		t.Errorf("Comparing nonexistent files should return error")
	}
}

func TestCreateTestFileEmpty(t *testing.T) {
	filepath, err := CreateTestFile("")
	if err != nil {
		t.Fatalf("CreateTestFile with empty content failed: %v", err)
	}
	
	defer CleanupTestFile(filepath)
	
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read empty test file: %v", err)
	}
	
	if len(data) != 0 {
		t.Errorf("Empty file should have zero length, got %d", len(data))
	}
}

func TestCreateTestFileWithUnicode(t *testing.T) {
	content := "Hello 世界! 🚀 Café"
	
	filepath, err := CreateTestFile(content)
	if err != nil {
		t.Fatalf("CreateTestFile with Unicode failed: %v", err)
	}
	
	defer CleanupTestFile(filepath)
	
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read Unicode test file: %v", err)
	}
	
	if string(data) != content {
		t.Errorf("Unicode content mismatch. Expected: %q, Got: %q", content, string(data))
	}
}