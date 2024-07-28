package pkg

import (
	"encoding/csv"
	"fmt"
	"os"
)

// GetCSVRecords opens a .csv file and returns the records
func GetCSVRecords(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read the CSV: %s", err.Error())
	}

	// Skip headers
	return records[1:], nil
}
