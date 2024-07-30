package csv

import (
	"database/sql"
	"fmt"
	"github.com/shvdg-dev/base-logic/pkg"
	"github.com/shvdg-dev/base-logic/pkg/testable/database"
	"log"
	"testing"
)

const compareMessageTemplate = "on row %d: [%s %s %s] from DB, [%s %s %s] from CSV"

// TestInsertingCSV verifies whether a .csv file can be inserted into the database.
func TestInsertingCSV(t *testing.T) {
	dbContainerService := database.NewContainerService()
	config := database.NewPostgresContainerConfig()
	dbContainer, err := dbContainerService.CreateContainer(config)
	if err != nil {
		t.Fatal(err)
	}

	defer dbContainer.Teardown()

	_, err = dbContainer.Query(createContactsTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	err = dbContainer.InsertCSVFile(myContactsCSVPath, contactsTable, contactsFields)
	if err != nil {
		log.Fatal(err)
	}

	csvRows, err := pkg.GetCSVRecords(myContactsCSVPath)
	if err != nil {
		log.Fatal(err)
	}

	tableRows, err := dbContainer.Query(getContactsQuery)
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

// compareRows compares the rows from the .csv file with those found in the table.
func compareRows(csvRows [][]string, tableRows *sql.Rows) error {
	index := 0
	for tableRows.Next() {
		var id, name, phone string
		if err := tableRows.Scan(&id, &name, &phone); err != nil {
			return err
		}

		match := id == csvRows[index][0] && name == csvRows[index][1] && phone == csvRows[index][2]
		compareMessage := fmt.Sprintf(compareMessageTemplate, index, id, name, phone, csvRows[index][0], csvRows[index][1], csvRows[index][2])

		if !match {
			return fmt.Errorf("FAIL: %s", compareMessage)
		} else {
			log.Printf("SUCC: %s", compareMessage)
		}

		index++
	}
	return nil
}
