package pkg

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type DbManagerOption func(*DbManager)

// DbManager represents a manger of the database connection.
type DbManager struct {
	DriverName, URL string
	*sql.DB
	IsMonitoringEnabled bool
}

// NewDbManager creates a new instance of DbManager.
func NewDbManager(driverName, URL string, options ...DbManagerOption) *DbManager {
	validateDriver(driverName)
	dbm := &DbManager{
		DriverName: driverName,
		URL:        URL,
	}
	for _, option := range options {
		option(dbm)
	}
	return dbm
}

// validateDriver checks if the given DriverName is "postgres"
func validateDriver(driverName string) {
	if driverName != "postgres" {
		log.Fatalf("Invalid driver; only 'postgres' is supported. Received: %s", driverName)
	}
}

// WithConnection attempts to connect with the database.
func WithConnection() DbManagerOption {
	return func(dbm *DbManager) {
		dbm.Connect()
	}
}

// WithMonitoring enables the connection monitoring for the database.
func WithMonitoring() DbManagerOption {
	return func(dbm *DbManager) {
		dbm.StartMonitoring()
	}
}

// Connect establishes a connection to the database using the specified driver and URL.
func (d *DbManager) Connect() {
	var err error
	d.DB, err = sql.Open(d.DriverName, d.URL)
	if err != nil {
		log.Printf("Failed to connect to database: %s", err.Error())
	}
}

// StartMonitoring monitors the database connection and attempts to reconnect whenever the database is not connected.
func (d *DbManager) StartMonitoring() {
	d.IsMonitoringEnabled = true
	for {
		if !d.IsMonitoringEnabled {
			break
		}
		err := d.DB.Ping()
		if err != nil {
			log.Printf("Lost connection to the database: %v", err)
			log.Printf("Attempting to reconnect...")
			d.Connect()
		}
		time.Sleep(15 * time.Second)
	}
}

// StopMonitoring disables the connection monitoring.
func (d *DbManager) StopMonitoring() {
	d.IsMonitoringEnabled = false
}

// Disconnect disconnects from the database.
func (d *DbManager) Disconnect() {
	d.IsMonitoringEnabled = false
	if d.DB == nil {
		return
	}
	err := d.DB.Close()
	if err != nil {
		log.Printf("Failed to diconnect from database: %s", err.Error())
	}
}

// CloseRows attempts to close the rows, a failure is logged, but no error is returned, as it is safe to ignore.
func (d *DbManager) CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("Failed to close rows: %s", err)
	}
}
