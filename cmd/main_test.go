package main

import (
	"os"
	"strings"
	"testing"
)

// TestProcessCustomerData tests the main application logic
func TestProcessCustomerData(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		outputPath  string
		expectError bool
	}{
		{
			name:        "Process with terminal output",
			inputPath:   "../customerimporter/test_data.csv",
			outputPath:  "",
			expectError: false,
		},
		{
			name:        "Process with file output",
			inputPath:   "../customerimporter/test_data.csv",
			outputPath:  "test_output.csv",
			expectError: false,
		},
		{
			name:        "Process with invalid input file",
			inputPath:   "nonexistent.csv",
			outputPath:  "",
			expectError: true,
		},
		{
			name:        "Process with large dataset",
			inputPath:   "../customerimporter/benchmark10k.csv",
			outputPath:  "large_test_output.csv",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.outputPath != "" {
				os.Remove(tt.outputPath)
			}

			err := ProcessCustomerData(tt.inputPath, tt.outputPath)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If we expected success and had an output file, verify it was created
			if !tt.expectError && tt.outputPath != "" {
				if _, err := os.Stat(tt.outputPath); os.IsNotExist(err) {
					t.Errorf("Expected output file %s to be created", tt.outputPath)
				} else {
					// Clean up the test file
					os.Remove(tt.outputPath)
				}
			}
		})
	}
}

// TestProcessCustomerDataTerminalOutput tests terminal output specifically
func TestProcessCustomerDataTerminalOutput(t *testing.T) {
	err := ProcessCustomerData("../customerimporter/test_data.csv", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestProcessCustomerDataFileOutput tests file output specifically
func TestProcessCustomerDataFileOutput(t *testing.T) {
	outputFile := "test_main_output.csv"
	defer os.Remove(outputFile) // Clean up

	err := ProcessCustomerData("../customerimporter/test_data.csv", outputFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the file was created and has content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
	}

	// Check that the file contains the expected header
	if !strings.Contains(string(content), "domain,number_of_customers") {
		t.Error("Output file should contain CSV header")
	}

	// Check that the file has some data
	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 { // Header + at least one data line
		t.Error("Output file should contain data")
	}
}

// TestProcessCustomerDataErrorHandling tests error scenarios
func TestProcessCustomerDataErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		inputPath  string
		outputPath string
	}{
		{
			name:       "Non-existent input file",
			inputPath:  "definitely_does_not_exist.csv",
			outputPath: "",
		},
		{
			name:       "Invalid output directory",
			inputPath:  "../customerimporter/test_data.csv",
			outputPath: "/invalid/path/that/does/not/exist/output.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessCustomerData(tt.inputPath, tt.outputPath)
			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

// BenchmarkProcessCustomerData benchmarks the main function
func BenchmarkProcessCustomerData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := ProcessCustomerData("../customerimporter/test_data.csv", "")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkProcessCustomerDataLarge benchmarks with large dataset
func BenchmarkProcessCustomerDataLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := ProcessCustomerData("../customerimporter/benchmark10k.csv", "")
		if err != nil {
			b.Fatal(err)
		}
	}
}
