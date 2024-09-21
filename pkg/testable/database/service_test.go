package database

import (
	"testing"
)

// TestDbContainer tests whether a database container can be created and a connection established.
func TestDbContainer(t *testing.T) {
	dbContainerService := NewContainerSvc()
	config := NewPostgresContainerConfig()
	dbContainer, err := dbContainerService.CreateContainer(config)
	if err != nil {
		t.Fatal(err)
	}

	defer dbContainer.Teardown()

	err = dbContainer.DB().Ping()
	if err != nil {
		t.Fatal(err)
	}
}
