package integration_tests

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/shvdg-dev/base-logic/pkg/testable"
	"log"
	"os"
	"testing"
)

// TestInsertingCSV verifies whether the inserting of .csv files into the database works.
func TestInsertingCSV(t *testing.T) {
	db, err := testable.NewDbContainer()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Teardown()

	_, err = db.Query(createContactsTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	err = db.InsertCSV(myContactsCSVPath, contactsTable, contactsFields)
	if err != nil {
		log.Fatal(err)
	}

	csvRows, err := getRowsFromCSV(myContactsCSVPath)
	if err != nil {
		log.Fatal(err)
	}

	tableRows, err := db.Query(getContactsQuery)
	if err != nil {
		log.Fatal(err)
	}

	defer tableRows.Close()

	err = compareRows(csvRows, tableRows)
	if err != nil {
		log.Fatal(err)
	}

	if err = tableRows.Err(); err != nil {
		log.Fatal(err)
	}
}

// getRowsFromCSV retrieves the rows, minus the headers, of the provided .csv file.
func getRowsFromCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return rows[1:], nil
}

// compareRows compares the rows from the .csv file with those found in the table.
func compareRows(csvRows [][]string, tableRows *sql.Rows) error {
	index := 0
	for tableRows.Next() {
		var id, name, phone string
		if err := tableRows.Scan(&id, &name, &phone); err != nil {
			return err
		}

		if id != csvRows[index][0] || name != csvRows[index][1] || phone != csvRows[index][2] {
			return fmt.Errorf("mismatch on row %d: got [%s %s %s] from DB, [%s %s %s] from CSV", index, id, name, phone, csvRows[index][0], csvRows[index][1], csvRows[index][2])
		}

		index++
	}
	return nil
}
