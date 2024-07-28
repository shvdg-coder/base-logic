package integration_tests

import (
	"github.com/shvdg-dev/base-logic/pkg/testable"
	"testing"
)

func TestDbContainer(t *testing.T) {
	db, err := testable.NewDbContainer()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Teardown()

	// Create a table
	_, err = db.Query("CREATE TABLE test_table (id INT PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Drop a table
	_, err = db.Query("DROP TABLE IF EXISTS test_table")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
}
