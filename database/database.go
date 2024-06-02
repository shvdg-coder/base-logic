package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

// Connection represents a connection to a database.
type Connection struct {
	DB *sql.DB
}

// NewConnection creates a new Connection and establishes a connection to the database using the specified URL.
func NewConnection(URL string) *Connection {
	database, err := sql.Open("postgres", URL)
	if err != nil {
		log.Fatal(err)
	}
	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &Connection{DB: database}
}
