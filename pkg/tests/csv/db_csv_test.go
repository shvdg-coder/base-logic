package csv

import (
	"github.com/shvdg-coder/base-logic/pkg"
	"github.com/shvdg-coder/base-logic/pkg/testable/database"
	"testing"
)

// TestInsertingCSV verifies whether a .csv file can be inserted into the database.
func TestInsertingCSV(t *testing.T) {
	dbContainer := setup(t)
	defer dbContainer.Teardown()

	// Execute
	err := dbContainer.InsertCSVFile(contactsCSVPath, contactsTableName, columnNames)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the rows in the CSV file
	csvRows, err := pkg.GetCSVRecords(contactsCSVPath, false)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve the rows in the table
	tableRows, err := dbContainer.DB().Query(getContactsQuery)
	if err != nil {
		t.Fatal(err)
	}

	defer tableRows.Close()

	if err = tableRows.Err(); err != nil {
		t.Fatal(err)
	}

	// Test
	err = pkg.CompareRows(csvRows, tableRows)
	if err != nil {
		t.Fatal(err)
	}

}

// setup prepares the tests by performing the minimally required steps.
func setup(t *testing.T) database.ContainerOps {
	// Instantiate a database container
	dbContainerService := database.NewContainerSvc()
	config := database.NewPostgresContainerConfig()
	dbContainer, err := dbContainerService.CreateContainer(config)
	if err != nil {
		t.Fatal(err)
	}

	// Create table
	_, err = dbContainer.DB().Query(createContactsTableQuery)
	if err != nil {
		t.Fatal(err)
	}

	return dbContainer
}
