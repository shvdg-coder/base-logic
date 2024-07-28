package testable

const (
	PostgresImage        = "postgres:13"
	PostgresPort         = "5432"
	Protocol             = "tcp"
	PostgresExposedPorts = PostgresPort + "/" + Protocol
	PostgresUser         = "docker"
	PostgresPassword     = "docker"
	PostgresDB           = "test"
)
