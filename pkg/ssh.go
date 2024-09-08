package pkg

import (
	"fmt"
	"github.com/elliotchance/sshtunnel"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

// SSHConfig holds data to create a SSH Tunnel.
type SSHConfig struct {
	User        string
	Password    string
	Server      string
	Destination string
	LocalPort   string
}

// SSHTunnel is helps with creating an SSH Tunnel.
type SSHTunnel struct {
	*sshtunnel.SSHTunnel
}

// NewSSHTunnel instantiates a SSHTunnel.
func NewSSHTunnel(config *SSHConfig) (*SSHTunnel, error) {
	server := fmt.Sprintf("%s@%s", config.User, config.Server)

	tunnel, err := sshtunnel.NewSSHTunnel(
		server,
		ssh.Password(config.Password),
		config.Destination,
		config.LocalPort,
	)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create SSH Tunnel: %s", err.Error()))
		return nil, err
	}

	return &SSHTunnel{SSHTunnel: tunnel}, nil
}

// Start starts the SSH Tunnel.
func (t *SSHTunnel) Start() {
	go func() {
		err := t.SSHTunnel.Start()
		if err != nil {
			log.Printf(fmt.Sprintf("Failed to start SSH Tunnel: %s" + err.Error()))
			return
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

// Close closes the SSH Tunnel.
func (t *SSHTunnel) Close() {
	t.SSHTunnel.Close()
}
