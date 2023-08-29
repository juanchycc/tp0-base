package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const MAX_BUFFER = 1024
const TYPE_MSG_POSITION = 0
const DOCUMENT_POSITION = 1
const NUMBER_POSITION = 2

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
	msgID := 1

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGTERM)

	apuestaMsg := ApuestaMsg{
		Apuesta: apuesta,
		Agency:  c.config.ID,
	}
	msgString := apuestaMsg.CreateMsgString()
	// si msg vacio, no tengo nada que mandar
	if msgString == "" || len(msgString) > MAX_BUFFER {
		log.Errorf(
			"action: create_msg | result: fail | client_id: %v",
			c.config.ID,
		)
		return
	}
	// por ahora len 1
	len := 1

	c.createClientSocket()

	for i := 0; i < len; i++ {
		select {
		case sig := <-sigchnl:

			log.Infof("action: signal_detected -> %v | result: success | client_id: %v", sig, c.config.ID)
			c.finish()
			return

		default:

			_, err := fmt.Fprint(
				c.conn,
				msgString,
			)
			if err != nil {
				log.Errorf(
					"action: send_data | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				c.finish()
				return
			}

			msg, err := bufio.NewReaderSize(c.conn, MAX_BUFFER).ReadString('\n')
			msgID++

			if err != nil {
				log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				c.finish()
				return
			}

			MsgHandler(msg)
		}

	}
	c.finish()
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) finish() {
	c.conn.Close()
}

func MsgHandler(msg string) {

	// TODO: mejorar
	res := strings.Split(msg, ";")
	if res[TYPE_MSG_POSITION] == "success" {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			res[DOCUMENT_POSITION],
			res[NUMBER_POSITION],
		)
	}

}
