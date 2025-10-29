package main

import (
	"fmt"
	"go-reloaded/internal/controller"
	"os"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_file> <output_file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s input.txt output.txt\n", os.Args[0])
		os.Exit(1)
	}
	
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	
	// Process the file
	err := controller.ProcessFile(inputFile, outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing file: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Successfully processed %s -> %s\n", inputFile, outputFile)
}