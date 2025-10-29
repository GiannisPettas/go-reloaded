package testutils

import (
	"fmt"
	"os"
)

// CreateTestFile creates a temporary test file with the given content
func CreateTestFile(content string) (string, error) {
	tmpDir := os.TempDir()
	file, err := os.CreateTemp(tmpDir, "go-reloaded-test-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	
	_, err = file.WriteString(content)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return "", fmt.Errorf("failed to write content: %w", err)
	}
	
	err = file.Close()
	if err != nil {
		os.Remove(file.Name())
		return "", fmt.Errorf("failed to close file: %w", err)
	}
	
	return file.Name(), nil
}

// CleanupTestFile removes a test file
func CleanupTestFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove test file %s: %w", path, err)
	}
	return nil
}

// CompareFiles compares two files and returns error if they differ
func CompareFiles(expected, actual string) error {
	expectedData, err := os.ReadFile(expected)
	if err != nil {
		return fmt.Errorf("failed to read expected file %s: %w", expected, err)
	}
	
	actualData, err := os.ReadFile(actual)
	if err != nil {
		return fmt.Errorf("failed to read actual file %s: %w", actual, err)
	}
	
	if string(expectedData) != string(actualData) {
		return fmt.Errorf("files differ:\nExpected: %q\nActual: %q", string(expectedData), string(actualData))
	}
	
	return nil
}