package database

import (
	"testing"
)

// TestDbContainer tests whether a database container can be created and a connection established.
func TestDbContainer(t *testing.T) {
	dbContainerService := NewContainerService()
	config := NewPostgresContainerConfig()
	dbContainer, err := dbContainerService.CreateContainer(config)
	if err != nil {
		t.Fatal(err)
	}

	defer dbContainer.Teardown()

	err = dbContainer.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
