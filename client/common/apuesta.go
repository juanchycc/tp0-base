package common

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
)

const MAX_BUFFER = 1024
const TYPE_MSG_POSITION = 0
const DOCUMENT_POSITION = 1
const NUMBER_POSITION = 2
const SUCCESS_BET_TYPE = "SUCCESS_BET"
const SINGLE_BET_TYPE = "SINGLE_BET"

type Apuesta struct {
	Name     string
	LastName string
	Document string
	Birthday string
	Number   string
}

type ApuestaMsg struct {
	Apuesta Apuesta
	Agency  string
}

// Crea un nuevo mensaje de apuesta a partir de los campos de la estructura ApuestaMsg recibida
// en caso de encontrar algún dato obligatorio vacio retorna un string vacio.
func (a *ApuestaMsg) CreateMsgString() string {

	if a.Agency == "" || a.Apuesta.Number == "" {
		return ""
	}

	apuesta := a.Apuesta

	return apuesta.Name + ";" + apuesta.LastName + ";" + apuesta.Document + ";" + apuesta.Birthday + ";" + apuesta.Number + "\n"
}

func enviarApuesta(conn net.Conn, ID string, apuesta string) error {
	err := sendSingleBet(apuesta, conn, ID)
	if err != nil {
		return err
	}
	return nil
}

func sendSingleBet(apuesta string, conn net.Conn, ID string) error {

	if len(apuesta) > MAX_BUFFER {
		log.Infof("action: sendBets| result: failed | client_id: %v | The number of bets exceeds %v", ID, MAX_BUFFER)

	} else {
		err := sendPacket(conn, SINGLE_BET_TYPE, ID, apuesta)
		if err != nil {
			return err
		}
	}

	_, err := getMsg(conn, SUCCESS_BET_TYPE)
	return err

}

func sendPacket(conn net.Conn, msgType string, ID string, msg string) error {
	msgLen := utf8.RuneCountInString(msg)
	header := msgType + ";" + strconv.Itoa(msgLen) + ";" + ID + "\n"

	packet := []byte(header + msg)
	packetLen := len(packet)

	totalWriteLen := 0

	//Si a la primera no escribe todo, manda lo que falta
	for totalWriteLen < packetLen {
		writeLen, err := conn.Write(packet[totalWriteLen:])
		if err != nil {
			return err
		}
		totalWriteLen += writeLen
		packet = packet[totalWriteLen:packetLen]
	}
	return nil
}

func getMsg(conn net.Conn, msgType string) ([]string, error) {
	leer := true
	var msg []string
	for leer {
		buffer := make([]byte, MAX_BUFFER)
		cant, err := bufio.NewReaderSize(conn, MAX_BUFFER).Read(buffer)
		if err != nil {
			return nil, err
		}
		recMsg := string(buffer[:cant])
		lines := strings.Split(recMsg, "\n")

		//Obtener Header:
		header := strings.Split(lines[0], ";")

		//Verifica el tipo de paquete
		if header[0] != msgType {
			return nil, nil
		}

		if header[1] == "0" {
			leer = false
		} else {
			//Verifica si llego todo:
			readLen := cant - len(lines[0]) - 1
			if header[1] == strconv.Itoa(readLen) {
				leer = false
			}
		}
		msg = append(msg, lines[1:]...)
	}
	return msg, nil
}
