package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/shvdg-dev/base-logic/pkg"
	tstcon "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ContainerOperations represents the operations regarding container management.
type ContainerOperations interface {
	CreateContainer(config *ContainerConfig) (*Container, error)
}

// ContainerService is responsible for managing containers.
type ContainerService struct {
	ContainerOperations
}

// NewContainerService instantiates a new ContainerService.
func NewContainerService() ContainerOperations {
	return &ContainerService{}
}

// CreateContainer creates a new instance of Container.
func (c *ContainerService) CreateContainer(config *ContainerConfig) (*Container, error) {
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

	dbm := pkg.NewDbManager(config.Driver, url, pkg.WithConnection())

	return &Container{
		Container: container,
		DbManager: dbm,
	}, nil
}

// newContainerRequest instantiates a new request for a container.
func (c *ContainerService) newContainerRequest(config *ContainerConfig) tstcon.ContainerRequest {
	exposedPort := config.Port + "/" + config.Protocol
	return tstcon.ContainerRequest{
		Image:        config.Image,
		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(exposedPort)),
		Env:          config.Env,
	}
}

// newContainer instantiates a new container using the provided context and request.
func (c *ContainerService) newContainer(ctx context.Context, request tstcon.ContainerRequest) (tstcon.Container, error) {
	return tstcon.GenericContainer(ctx, tstcon.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
}

// createURL constructs the database connection URL.
func (c *ContainerService) createURL(host, port string, config *ContainerConfig) (string, error) {
	if config.Driver == "postgres" {
		return c.createPostgresURL(host, port, config)
	}
	return "", errors.New(fmt.Sprintf("unsupported driver for database URL creation: %s", config.Driver))
}

// createPostgresURL constructs Postgres database connection URL.
func (c *ContainerService) createPostgresURL(host, port string, config *ContainerConfig) (string, error) {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		config.User,
		config.Password,
		config.DbName), nil
}
