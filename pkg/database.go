package pkg

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

// DbSvcOption is used to instantiate a DbSvc with the provided settings/configurations/actions.
type DbSvcOption func(*DbSvc)

// DbOps represents operations related to database actions
type DbOps interface {
	Connect()
	Disconnect()
	StartMonitoring()
	StopMonitoring()
	InsertCSVFile(filePath, table string, fields []string) error
	BulkInsert(table string, fields []string, data [][]interface{}) error
	DB() *sql.DB
}

// DbSvc represents a manger of the database connection.
type DbSvc struct {
	DriverName, URL     string
	IsMonitoringEnabled bool
	SSHTunnel           *SSHTunnel
	db                  *sql.DB
}

// NewDbSvc creates a new instance of DbSvc.
func NewDbSvc(driverName, URL string, options ...DbSvcOption) *DbSvc {
	validateDriver(driverName)
	dbm := &DbSvc{
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

// DB returns the underlying *sql.DB instance used for the database connection.
func (d *DbSvc) DB() *sql.DB {
	return d.db
}

// WithSSHTunnel establishes an SSH tunnel for connecting to the database.
func WithSSHTunnel(config *SSHConfig) DbSvcOption {
	return func(dbs *DbSvc) {
		sshTunnel, err := NewSSHTunnel(config)
		if err != nil {
			log.Printf("Unable to establish SSH tunnel: %s", err.Error())
		}
		dbs.SSHTunnel = sshTunnel
	}
}

// WithConnection attempts to connect with the database.
func WithConnection() DbSvcOption {
	return func(dbs *DbSvc) {
		dbs.Connect()
	}
}

// WithMonitoring enables the connection monitoring for the database.
func WithMonitoring() DbSvcOption {
	return func(dbs *DbSvc) {
		dbs.StartMonitoring()
	}
}

// Connect establishes a connection to the database using the specified driver and URL.
func (d *DbSvc) Connect() {
	dbURL := d.URL
	if d.SSHTunnel != nil {
		d.SSHTunnel.Start()
		dbURL = strings.Replace(dbURL, "<PORT>", strconv.Itoa(d.SSHTunnel.Local.Port), 1)
	}

	var err error
	d.db, err = sql.Open(d.DriverName, dbURL)
	if err != nil {
		log.Printf("Failed to connect to database: %s", err.Error())
	}

	err = d.DB().Ping()
	if err != nil {
		log.Printf("Failed to reach database: %s", err.Error())
	}
}

// Disconnect disconnects from the database.
func (d *DbSvc) Disconnect() {
	d.IsMonitoringEnabled = false
	if d.db == nil {
		return
	}

	err := d.DB().Close()
	if err != nil {
		log.Printf("Failed to diconnect from database: %s", err.Error())
	}

	if d.SSHTunnel != nil {
		d.SSHTunnel.Close()
	}
}

// StartMonitoring monitors the database connection and attempts to reconnect whenever the database is not connected.
func (d *DbSvc) StartMonitoring() {
	d.IsMonitoringEnabled = true
	for {
		if !d.IsMonitoringEnabled {
			break
		}
		err := d.DB().Ping()
		if err != nil {
			log.Printf("Lost connection to the database: %v", err)
			log.Printf("Attempting to reconnect...")
			d.Connect()
		}
		time.Sleep(15 * time.Second)
	}
}

// StopMonitoring disables the connection monitoring.
func (d *DbSvc) StopMonitoring() {
	d.IsMonitoringEnabled = false
}

// InsertCSVFile is the main function that coordinates opening the file and inserting the records to the database
func (d *DbSvc) InsertCSVFile(filePath, table string, fields []string) error {
	records, err := GetCSVRecords(filePath, false)
	if err != nil {
		return err
	}
	return d.insertCSVRecords(table, fields, records)
}

// insertCSVRecords inserts the contents of a .csv file into the database.
func (d *DbSvc) insertCSVRecords(table string, fields []string, records [][]string) error {
	transaction, err := d.DB().Begin()
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

// BulkInsert helps inserting data in bulk.
func (d *DbSvc) BulkInsert(table string, fields []string, data [][]interface{}) error {
	txn, err := d.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(table, fields...))
	if err != nil {
		return fmt.Errorf("failed preparing statement: %w", err)
	}

	for _, row := range data {
		_, err := stmt.Exec(row...)
		if err != nil {
			return fmt.Errorf("failed executing statement: %w", err)
		}
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("failed closing statement: %w", err)
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	return nil
}
