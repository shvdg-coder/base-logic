package database

import (
	"testing"
)

// TestDbContainer tests whether a database container can be created and a connection established.
func TestDbContainer(t *testing.T) {
	containerService := NewContainerService()
	config := NewPostgresContainerConfig()
	container, err := containerService.CreateContainer(config)
	if err != nil {
		t.Fatal(err)
	}

	defer container.Teardown()

	err = container.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
