package pkg

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	strings "strings"
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
	columns, err := tableRows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %v", err)
	}

	values, valuePointers := make([]interface{}, len(columns)), make([]interface{}, len(columns))
	for i := range values {
		valuePointers[i] = &values[i]
	}

	index := 0
	for tableRows.Next() {
		if err := tableRows.Scan(valuePointers...); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		if index >= len(csvRows) {
			return fmt.Errorf("the contents of CSV and table are not equal: table has more rows than CSV")
		}

		if err := compareRows(csvRows[index], values, columns, index); err != nil {
			return err
		}

		index++
	}

	if index < len(csvRows) {
		return fmt.Errorf("the contents of CSV and table are not equal: table has fewer rows than CSV")
	}

	return nil
}

// compareRows compares a single row from the table with a single row from the CSV, and throws an error if they do not match.
func compareRows(csvRow []string, tableRow []interface{}, columns []string, rowIndex int) error {
	rowCompare, csvCompare := make([]string, len(columns)), make([]string, len(columns))
	match := true

	for i, colVal := range tableRow {
		val := fmt.Sprintf("%v", colVal)
		rowCompare[i] = val

		if i < len(csvRow) {
			csvCompare[i] = csvRow[i]
			if csvRow[i] != val {
				match = false
			}
		} else {
			match = false
		}
	}

	if !match {
		return createMismatchError(rowIndex, rowCompare, csvCompare)
	}

	return nil
}

// createMismatchError formats a meaningful error message for mismatched rows.
func createMismatchError(rowIndex int, tableRow, csvRow []string) error {
	compareMessage := fmt.Sprintf(
		"row %d: [%s] from table, [%s] from CSV", rowIndex,
		strings.Join(tableRow, ", "),
		strings.Join(csvRow, ", "),
	)
	return fmt.Errorf("the contents of CSV and table are not equal: %s", compareMessage)
}
