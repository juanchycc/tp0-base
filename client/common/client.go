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
const SUCCESS_BET_TYPE = "SUCCESS"

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

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGTERM)

	apuestaMsg := ApuestaMsg{
		Apuesta: apuesta,
		Agency:  c.config.ID,
	}
	msgString := apuestaMsg.CreateMsgString()
	totalLen := len(msgString)

	// si msg vacio, no tengo nada que mandar
	if msgString == "" {
		log.Errorf(
			"action: create_msg | result: fail | client_id: %v",
			c.config.ID,
		)
		return
	}
	// por ahora len 1, una sola apuesta
	len := 1

	c.createClientSocket()

	for i := 0; i < len; i++ {
		select {
		case sig := <-sigchnl:

			log.Infof("action: signal_detected -> %v | result: success | client_id: %v", sig, c.config.ID)
			c.finish()
			return

		default:

			msgError := false
			if totalLen > MAX_BUFFER {

				totalMsg := totalLen / MAX_BUFFER
				initialPos := 0

				for i := 0; i < totalMsg; i++ {
					s := msgString[initialPos : MAX_BUFFER+initialPos]
					msgError = c.sengMsg(s)
					if msgError {
						break
					}
					initialPos += MAX_BUFFER
				}

			} else {
				msgError = c.sengMsg(msgString)
			}

			if msgError {
				c.finish()
				return
			}

			msg, err := bufio.NewReaderSize(c.conn, MAX_BUFFER).ReadString('\n')

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

	// Wait a time between sending one message and the next one
	timeout := time.After(c.config.LoopPeriod)
	select {
	case <-timeout:
	case sig := <-sigchnl:
		log.Infof("action: signal_detected -> %v | result: success | client_id: %v", sig, c.config.ID)
		return
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)

}

func (c *Client) finish() {
	c.conn.Close()
}

func (c *Client) sengMsg(msg string) bool {

	len := len(msg)

	writeBytes, err := fmt.Fprint(
		c.conn,
		msg,
	)

	if err != nil {
		log.Errorf(
			"action: send_data | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		c.finish()
		return true
	}
	// parte del msg no se envio, llamado recursvio para que mande el resto
	if len > writeBytes {
		return c.sengMsg(msg[writeBytes:len])
	}
	return false
}

func MsgHandler(msg string) {

	msgLen := len(msg)
	if msgLen == 0 {
		return
	}

	msg = msg[0 : msgLen-1]
	res := strings.Split(msg, ";")
	if res[TYPE_MSG_POSITION] == SUCCESS_BET_TYPE {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			res[DOCUMENT_POSITION],
			res[NUMBER_POSITION],
		)
	}
}
