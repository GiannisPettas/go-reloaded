package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestAllProject runs all tests in the project
func TestAllProject(t *testing.T) {
	// Get project root directory
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Run tests on all packages except testutils to avoid recursion
	cmd := exec.Command("go", "test", "-count=1", "-v", "./cmd/...", "./internal/config", "./internal/controller", "./internal/exporter", "./internal/parser", "./internal/transformer")
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("🧪 Running all tests in go-reloaded project...")
	fmt.Println(strings.Repeat("=", 50))

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			t.Fatalf("Main tests failed with exit code %d", exitError.ExitCode())
		}
		t.Fatalf("Main tests failed: %v", err)
	}

	// Run golden tests separately
	fmt.Println("\n🏆 Running Golden Test Suite...")
	cmdGolden := exec.Command("go", "test", "-count=1", "-v", "-run=TestGoldenCases")
	cmdGolden.Dir = filepath.Join(projectRoot, "internal", "testutils")
	cmdGolden.Stdout = os.Stdout
	cmdGolden.Stderr = os.Stderr

	err = cmdGolden.Run()
	if err != nil {
		t.Fatalf("Golden tests failed: %v", err)
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✅ All tests passed!")
}

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}