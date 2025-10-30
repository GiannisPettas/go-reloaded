package testutils

import (
	"go-reloaded/internal/controller"
	"os"
	"testing"
)

func TestGoldenCases(t *testing.T) {
	tests, err := ParseGoldenTests("../../docs/golden_tests.md")
	if err != nil {
		t.Fatalf("Failed to parse golden tests: %v", err)
	}
	
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			inputPath, err := CreateTestFile(test.Input)
			if err != nil {
				t.Fatalf("Failed to create input file: %v", err)
			}
			defer CleanupTestFile(inputPath)
			
			outputPath, err := CreateTestFile("")
			if err != nil {
				t.Fatalf("Failed to create output file: %v", err)
			}
			defer CleanupTestFile(outputPath)
			
			err = controller.ProcessFile(inputPath, outputPath)
			if err != nil {
				t.Fatalf("ProcessFile failed: %v", err)
			}
			
			actualData, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}
			
			actual := string(actualData)
			if actual != test.Expected {
				t.Errorf("\nExpected: %q\nActual:   %q", test.Expected, actual)
			}
		})
	}
}