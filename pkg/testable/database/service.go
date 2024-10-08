package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/shvdg-coder/base-logic/pkg"
	tstcon "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ContainerManagement represents the actions regarding container management.
type ContainerManagement interface {
	CreateContainer(config *ContainerConfig) (*Container, error)
}

// ContainerSvc is responsible for managing containers.
type ContainerSvc struct {
	ContainerManagement
}

// NewContainerSvc instantiates a new ContainerSvc.
func NewContainerSvc() *ContainerSvc {
	return &ContainerSvc{}
}

// CreateContainer creates a new instance of Container.
func (c *ContainerSvc) CreateContainer(config *ContainerConfig) (*Container, error) {
	ctx := context.Background()

	request := c.newContainerRequest(config)
	container, err := c.newContainer(ctx, request)
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	extPort, err := container.MappedPort(ctx, nat.Port(config.Port))
	if err != nil {
		return nil, err
	}

	url, err := c.createURL(host, extPort.Port(), config)
	if err != nil {
		return nil, err
	}

	dbs := pkg.NewDbSvc(config.Driver, url, pkg.WithConnection())
	dbContainer := NewContainer(container, dbs)

	return dbContainer, nil
}

// newContainerRequest instantiates a new request for a container.
func (c *ContainerSvc) newContainerRequest(config *ContainerConfig) tstcon.ContainerRequest {
	exposedPort := config.Port + "/" + config.Protocol
	return tstcon.ContainerRequest{
		Image:        config.Image,
		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(exposedPort)),
		Env:          config.Env,
	}
}

// newContainer instantiates a new container using the provided context and request.
func (c *ContainerSvc) newContainer(ctx context.Context, request tstcon.ContainerRequest) (tstcon.Container, error) {
	return tstcon.GenericContainer(ctx, tstcon.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
}

// createURL constructs the database connection URL.
func (c *ContainerSvc) createURL(host, port string, config *ContainerConfig) (string, error) {
	if config.Driver == "postgres" {
		return c.createPostgresURL(host, port, config)
	}
	return "", errors.New(fmt.Sprintf("unsupported driver for database URL creation: %s", config.Driver))
}

// createPostgresURL constructs Postgres database connection URL.
func (c *ContainerSvc) createPostgresURL(host, port string, config *ContainerConfig) (string, error) {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		config.User,
		config.Password,
		config.DbName), nil
}
