package ssh

import (
	_ "github.com/lib/pq"
	"github.com/shvdg-coder/base-logic/pkg"
	"testing"
)

// TestStartTunnel verifies it is possible to connect to a server via a SSH tunnel.
func TestStartTunnel(t *testing.T) {
	// Set up the ssh tunnel configuration
	sshConfig := &pkg.SSHConfig{
		User:        pkg.GetEnvValueAsString(userKey),
		Password:    pkg.GetEnvValueAsString(passwordKey),
		Server:      pkg.GetEnvValueAsString(serverKey),
		Destination: pkg.GetEnvValueAsString(destinationKey),
		LocalPort:   pkg.GetEnvValueAsString(localPortKey),
	}

	// Try to connect to the database
	dbService := pkg.NewDbService("postgres",
		pkg.GetEnvValueAsString(databaseURL),
		pkg.WithSSH(sshConfig),
		pkg.WithConnection())
	defer dbService.Disconnect()

	dbService.SSHTunnel.Start()

	// Test if able to ping the database
	err := dbService.Ping()
	if err != nil {
		t.Fatalf("Could not ping database: %s", err.Error())
	}

	dbService.SSHTunnel.Close()
}
