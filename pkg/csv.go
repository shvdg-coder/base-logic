package pkg

import (
	"database/sql"
	"encoding/csv"
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

// CompareRows compares the rows from a .csv file with those from a query and returns an error when they are not equal.
func CompareRows(csvRows [][]string, tableRows *sql.Rows) error {
	index := 0
	for tableRows.Next() {
		var id, name, phone string
		if err := tableRows.Scan(&id, &name, &phone); err != nil {
			return err
		}

		match := id == csvRows[index][0] && name == csvRows[index][1] && phone == csvRows[index][2]

		if !match {
			compareMessage := fmt.Sprintf(
				"row %d: [%s %s %s] from table, [%s %s %s] from CSV", index,
				id, name, phone,
				csvRows[index][0], csvRows[index][1], csvRows[index][2])

			return fmt.Errorf("the contents of CSV and table are not equal: %s", compareMessage)
		}

		index++
	}
	return nil
}
