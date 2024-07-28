package testable

import (
	"testing"
)

// TestDbContainer tests whether a database container can be created and a connection established.
func TestDbContainer(t *testing.T) {
	db, err := NewDbContainer()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Teardown()

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
