package database

import (
	"database/sql"
	"log"
	"time"
)

// Manager represents a manger of the database connection.
type Manager struct {
	DB *sql.DB
}

// NewManager creates a new instance of Manager.
func NewManager(driverName string, URL string) *Manager {
	manager := &Manager{}
	manager.Connect(driverName, URL)
	go manager.ConnectionMonitor(driverName, URL)
	return manager
}

// Connect establishes a connection to the database using the specified driver and URL.
func (m *Manager) Connect(driverName, URL string) {
	var err error
	m.DB, err = sql.Open(driverName, URL)
	if err != nil {
		log.Printf("Failed to connect to database")
	}
}

// ConnectionMonitor runs in a continuous loop to monitor the database connection.
func (m *Manager) ConnectionMonitor(driverName, URL string) {
	for {
		time.Sleep(time.Minute)
		err := m.DB.Ping()
		if err != nil {
			log.Printf("Lost connection to the database: %v", err)
			log.Printf("Attempting to reconnect...")
			m.Connect(driverName, URL)
		}
	}
}
