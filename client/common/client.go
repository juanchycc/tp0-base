package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(apuesta Apuesta) {
	// autoincremental msgID to identify every message sent

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGTERM)

	c.createClientSocket()

	err := leerApuestas(c.config.ID, c.conn, sigchnl)
	if err != nil {
		log.Errorf(
			"action: send_data | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		c.finish()
		return
	}

	c.finish()

	// Wait a time between sending one message and the next one
	timeout := time.After(c.config.LoopPeriod)
	select {
	case <-timeout:
	case <-sigchnl:
		log.Infof("action: signal_detected | result: success | client_id: %v", c.config.ID)
		return
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)

}

func (c *Client) finish() {
	c.conn.Close()
}
