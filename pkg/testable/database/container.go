package database

import (
	"context"
	"github.com/shvdg-dev/base-logic/pkg"
	tstcon "github.com/testcontainers/testcontainers-go"
)

// Container is used to spin up a database in a container for integration testing.
type Container struct {
	tstcon.Container
	*pkg.DbManager
}

// Teardown destroys the database container.
func (t *Container) Teardown() error {
	t.Disconnect()
	return t.Terminate(context.Background())
}
