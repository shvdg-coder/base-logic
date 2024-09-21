package database

import (
	"context"
	"github.com/shvdg-coder/base-logic/pkg"
	tstcon "github.com/testcontainers/testcontainers-go"
)

// ContainerOperations represents operations related to a database container.
type ContainerOperations interface {
	pkg.DbOps
	Teardown() error
}

// ContainerWrapper represents a wrapper around testcontainers.Container type.
type ContainerWrapper struct {
	tstcon.Container
}

// Container represents a database container, which can do both container and database actions.
type Container struct {
	ContainerWrapper
	pkg.DbOps
}

// NewContainer creates a new instance of Container.
func NewContainer(container tstcon.Container, database pkg.DbOps) *Container {
	return &Container{
		ContainerWrapper: ContainerWrapper{Container: container},
		DbOps:            database,
	}
}

// Teardown destroys the database container.
func (t *Container) Teardown() error {
	t.Disconnect()
	return t.Terminate(context.Background())
}
