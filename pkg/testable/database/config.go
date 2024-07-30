package database

type ContainerConfig struct {
	Driver   string
	Image    string
	Port     string
	Protocol string

	User     string
	Password string
	DbName   string

	Env map[string]string
}

// NewPostgresContainerConfig creates a default configuration for a Postgres database container.
func NewPostgresContainerConfig() *ContainerConfig {
	config := &ContainerConfig{
		Driver:   "postgres",
		Image:    "postgres:13",
		Port:     "5432",
		Protocol: "tcp",
		User:     "docker",
		Password: "docker",
		DbName:   "test",
	}

	config.Env = map[string]string{
		"POSTGRES_PASSWORD": config.Password,
		"POSTGRES_USER":     config.User,
		"POSTGRES_DB":       config.DbName,
	}

	return config
}
