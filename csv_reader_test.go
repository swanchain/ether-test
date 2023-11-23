package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

// ReadAndPrintCSV reads a CSV file, skips the header, and prints the records.
func ReadAndPrintCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isHeader := true

	for scanner.Scan() {
		if isHeader {
			// Skip the header
			isHeader = false
			continue
		}
		// fmt.Println(scanner.Text()) // Print each record
	}

	return scanner.Err()
}

func ReadCSVAndPrintCount(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	recordCount := 0
	isHeader := true

	for scanner.Scan() {
		if isHeader {
			// Skip the header
			isHeader = false
			continue
		}
		recordCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	fmt.Printf("Total records (excluding header): %d\n", recordCount)
	return recordCount, nil
}

func TestReadAndPrintCSV(t *testing.T) {

	// Call the ReadAndPrintCSV function
	recordCount, err := ReadCSVAndPrintCount("./../../../ethereum-address/000000000000.csv")
	if err != nil {
		t.Errorf("Failed to read CSV file and print count: %v", err)
	}

	// Assert the record count (change the expected count based on your actual test data)
	expectedCount := 4515349 // Set this to the actual number of records in your test data
	if recordCount != expectedCount {
		t.Errorf("Expected %d records, got %d", expectedCount, recordCount)
	}

}
