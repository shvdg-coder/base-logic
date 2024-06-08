package pkg

import (
	"database/sql"
	"log"
	"time"
)

// DatabaseManager represents a manger of the database connection.
type DatabaseManager struct {
	DB *sql.DB
}

// NewDatabaseManager creates a new instance of DatabaseManager.
func NewDatabaseManager(driverName string, URL string) *DatabaseManager {
	manager := &DatabaseManager{}
	manager.Connect(driverName, URL)
	go manager.ConnectionMonitor(driverName, URL)
	return manager
}

// Connect establishes a connection to the database using the specified driver and URL.
func (d *DatabaseManager) Connect(driverName, URL string) {
	var err error
	d.DB, err = sql.Open(driverName, URL)
	if err != nil {
		log.Printf("Failed to connect to database")
	}
}

// ConnectionMonitor runs in a continuous loop to monitor the database connection.
func (d *DatabaseManager) ConnectionMonitor(driverName, URL string) {
	for {
		time.Sleep(15 * time.Second)
		err := d.DB.Ping()
		if err != nil {
			log.Printf("Lost connection to the database: %v", err)
			log.Printf("Attempting to reconnect...")
			d.Connect(driverName, URL)
		}
	}
}
