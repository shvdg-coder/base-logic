package pkg

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"time"
)

// DbManagerOption is used to instantiate a DbManager with the provided settings/configurations/actions.
type DbManagerOption func(*DbManager)

// DbOperations represents operations related to database actions
type DbOperations interface {
	Connect()
	Disconnect()
	StartMonitoring()
	StopMonitoring()
	InsertCSVFile(filePath, table string, fields []string) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// DbManager represents a manger of the database connection.
type DbManager struct {
	DriverName, URL     string
	IsMonitoringEnabled bool
	*sql.DB
}

// NewDbManager creates a new instance of DbManager.
func NewDbManager(driverName, URL string, options ...DbManagerOption) DbOperations {
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

// StartMonitoring monitors the database connection and attempts to reconnect whenever the database is not connected.
func (d *DbManager) StartMonitoring() {
	d.IsMonitoringEnabled = true
	for {
		if !d.IsMonitoringEnabled {
			break
		}
		err := d.Ping()
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

// InsertCSVFile is the main function that coordinates opening the file and inserting the records to the database
func (d *DbManager) InsertCSVFile(filePath, table string, fields []string) error {
	records, err := GetCSVRecords(filePath)
	if err != nil {
		return err
	}
	return d.InsertCSVRecords(table, fields, records)
}

// InsertCSVRecords inserts the contents of a .csv file into the database.
func (d *DbManager) InsertCSVRecords(table string, fields []string, records [][]string) error {
	transaction, err := d.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %s", err.Error())
	}

	statement, err := transaction.Prepare(pq.CopyIn(table, fields...))
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %s", err.Error())
	}

	for _, record := range records {
		data := make([]interface{}, len(record))
		for i, v := range record {
			data[i] = v
		}
		if _, err = statement.Exec(data...); err != nil {
			return fmt.Errorf("failed to execute statement: %s", err.Error())
		}
	}

	if err = statement.Close(); err != nil {
		return fmt.Errorf("failed to close statement: %s", err.Error())
	}

	if err = transaction.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %s", err.Error())
	}

	return nil
}
