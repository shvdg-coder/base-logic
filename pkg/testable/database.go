package testable

import (
	"context"
	"fmt"
	"github.com/shvdg-dev/base-logic/pkg"
	tstcon "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// DbContainer is used to spin up a database in a container for integration testing.
type DbContainer struct {
	tstcon.Container
	*pkg.DbManager
}

// NewDbContainer creates a new instance of DbContainer.
func NewDbContainer() (*DbContainer, error) {
	ctx := context.Background()

	container, err := NewPostgresContainer(ctx)
	if err != nil {
		return nil, err
	}

	host, port, err := GetContainerInfo(ctx, container)
	if err != nil {
		return nil, err
	}

	URL := CreateURL(host, port)
	dbm := pkg.NewDbManager("postgres", URL, pkg.WithConnection())

	return &DbContainer{
		Container: container,
		DbManager: dbm,
	}, nil
}

// NewPostgresContainer sets up a new Postgres container in Docker and returns the container instance.
func NewPostgresContainer(ctx context.Context) (tstcon.Container, error) {
	req := tstcon.ContainerRequest{
		Image:        PostgresImage,
		ExposedPorts: []string{PostgresExposedPorts},
		WaitingFor:   wait.ForListeningPort(PostgresExposedPorts),
		Env: map[string]string{
			"POSTGRES_PASSWORD": PostgresPassword,
			"POSTGRES_USER":     PostgresUser,
			"POSTGRES_DB":       PostgresDB,
		},
	}
	return tstcon.GenericContainer(ctx, tstcon.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// GetContainerInfo fetches and returns the host IP and mapped port from a running Docker container.
func GetContainerInfo(ctx context.Context, container tstcon.Container) (string, string, error) {
	ip, err := container.Host(ctx)
	if err != nil {
		return "", "", err
	}

	port, err := container.MappedPort(ctx, PostgresPort)
	if err != nil {
		return "", "", err
	}

	return ip, port.Port(), nil
}

// CreateURL constructs and returns the database connection URL using the constants and provided host and port.
func CreateURL(host, port string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		PostgresUser,
		PostgresPassword,
		PostgresDB)
}

// Teardown destroys the database container.
func (t *DbContainer) Teardown() error {
	t.Disconnect()
	return t.Terminate(context.Background())
}
