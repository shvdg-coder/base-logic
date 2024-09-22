package pkg

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

// GetCSVColumnValues retrieves the value of the provided column name of each record in the CSV.
func GetCSVColumnValues(filePath, columnName string) ([]string, error) {
	records, err := GetCSVRecords(filePath, true)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("file %s has no records", filePath)
	}

	columnIndex, err := findColumnIndex(records, columnName)
	if err != nil {
		return nil, err
	}

	return getColumnValues(records, columnIndex), nil
}

// findColumnIndex returns the index of the specified column name from the records.
func findColumnIndex(records [][]string, columnName string) (int, error) {
	headers := records[0]
	for i, name := range headers {
		if name == columnName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("column %s not found", columnName)
}

// getColumnValues retrieves the values of a column specified by its index.
func getColumnValues(records [][]string, index int) []string {
	valueRecords := records[1:] // Skipping headers
	columnValues := make([]string, len(valueRecords))
	for i := 0; i < len(valueRecords); i++ {
		columnValues[i] = valueRecords[i][index]
	}
	return columnValues
}

// GetCSVRecords opens a .csv file and returns the records.
func GetCSVRecords(filePath string, headers bool) ([][]string, error) {
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

	if headers {
		return records, nil
	} else {
		return records[1:], nil // Skip headers
	}
}

// CompareRows compares the rows from a CSV file with those from a SQL query and returns an error if they are not equal.
func CompareRows(csvRows [][]string, tableRows *sql.Rows) error {
	var scanResults [][]interface{}

	columns, err := tableRows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %v", err)
	}

	for tableRows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := tableRows.Scan(valuePointers...); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		scanResults = append(scanResults, values)
	}

	if len(csvRows) != len(scanResults) {
		return errors.New("number of rows in CSV and SQL result set do not match")
	}

	for index, scannedRow := range scanResults {
		if !areRowsEqual(csvRows[index], scannedRow) {
			return fmt.Errorf("rows are not equal at index %d", index)
		}
	}

	return nil
}

// areRowsEqual compares a single row from the table with a single row from the CSV.
func areRowsEqual(csvRow []string, tableRow []interface{}) bool {
	if len(csvRow) != len(tableRow) {
		return false
	}

	for i, tableValue := range tableRow {
		csvValue := csvRow[i]
		tableValueString := fmt.Sprintf("%v", tableValue)
		if csvValue != tableValueString {
			return false
		}
	}

	return true
}
