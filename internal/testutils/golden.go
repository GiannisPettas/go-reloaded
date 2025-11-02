package testutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type GoldenTest struct {
	Name     string
	Input    string
	Expected string
}

// ParseGoldenTests reads and parses golden_tests.md file
func ParseGoldenTests(filePath string) ([]GoldenTest, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open golden tests file: %w", err)
	}
	defer file.Close()

	var tests []GoldenTest
	scanner := bufio.NewScanner(file)
	
	var currentTest GoldenTest
	var inInput, inExpected bool
	var inputBuilder, expectedBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		
		// Parse test name
		if strings.HasPrefix(line, "## T") && strings.Contains(line, "—") {
			// Save previous test if exists
			if currentTest.Name != "" {
				currentTest.Input = strings.TrimSpace(inputBuilder.String())
				currentTest.Expected = strings.TrimSpace(expectedBuilder.String())
				tests = append(tests, currentTest)
			}
			
			// Start new test
			parts := strings.Split(line, "—")
			if len(parts) >= 2 {
				currentTest = GoldenTest{Name: strings.TrimSpace(parts[0][3:])}
				inputBuilder.Reset()
				expectedBuilder.Reset()
				inInput = false
				inExpected = false
			}
		}
		
		// Parse input section
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "**Input:**" {
			inInput = true
			inExpected = false
			continue
		}
		
		// Parse expected output section
		if trimmedLine == "**Expected Output:**" {
			inInput = false
			inExpected = true
			continue
		}
		
		// Stop parsing when hitting next section
		if strings.HasPrefix(trimmedLine, "**") && trimmedLine != "**Input:**" && trimmedLine != "**Expected Output:**" {
			inInput = false
			inExpected = false
		}
		
		// Collect input/expected content
		if inInput && line != "" {
			if inputBuilder.Len() > 0 {
				inputBuilder.WriteByte('\n')
			}
			inputBuilder.WriteString(line)
		}
		
		if inExpected && line != "" {
			if expectedBuilder.Len() > 0 {
				expectedBuilder.WriteByte('\n')
			}
			expectedBuilder.WriteString(line)
		}
	}
	
	// Save last test
	if currentTest.Name != "" {
		currentTest.Input = strings.TrimSpace(inputBuilder.String())
		currentTest.Expected = strings.TrimSpace(expectedBuilder.String())
		tests = append(tests, currentTest)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	
	return tests, nil
}