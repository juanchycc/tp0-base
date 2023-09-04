package common

import (
	"bufio"
	"encoding/csv"
	"net"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
)

const SINGLE_BET_TYPE = "SINGLE_BET"
const MULTIPLE_BET_TYPE = "MULTIPLE_BET"
const FINISH_BET_TYPE = "FINISH_BET"
const SUCCESS_BET_TYPE = "SUCCESS_BET"
const WINNERS_BET_TYPE = "WINNERS_BET"
const TOPE_APUESTAS = 100
const MAX_BUFFER = 8192
const TYPE_MSG_POSITION = 0
const DOCUMENT_POSITION = 1
const NUMBER_POSITION = 2

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

	apuesta := a.Apuesta

	return apuesta.Name + ";" + apuesta.LastName + ";" + apuesta.Document + ";" + apuesta.Birthday + ";" + apuesta.Number + "\n"
}

func leerApuestas(conn net.Conn, ID string) error {

	file, err := os.Open("./app/data/agency-" + ID + ".csv")
	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	apuestas := ""
	cantFilas := 1

	for {

		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				//llega al final y hay apuestas pendientes, se mandan
				if len(apuestas) != 0 {
					sendBets(apuestas, conn, ID)
				}
				break
			}
			return err
		}

		nuevaApuesta := Apuesta{
			Name:     record[0],
			LastName: record[1],
			Document: record[2],
			Birthday: record[3],
			Number:   record[4],
		}

		nuevaApuestaMsg := ApuestaMsg{
			Apuesta: nuevaApuesta,
			Agency:  ID,
		}

		apuestas = apuestas + nuevaApuestaMsg.CreateMsgString()

		if cantFilas == TOPE_APUESTAS {
			err = sendBets(apuestas, conn, ID)
			if err != nil {
				return err
			}

			apuestas = ""
			cantFilas = 1

		} else {
			cantFilas++
		}

	}

	return sendPacket(conn, FINISH_BET_TYPE, ID, "")

}

func sendBets(apuestas string, conn net.Conn, ID string) error {

	if len(apuestas) > MAX_BUFFER {
		log.Infof("action: sendBets| result: failed | client_id: %v | The number of bets exceeds 8KB", ID)

	} else {
		err := sendPacket(conn, MULTIPLE_BET_TYPE, ID, apuestas)
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

func getWinners(conn net.Conn, ID string) error {

	var msg []string
	var err error

	for len(msg) == 0 {
		msg, err = getMsg(conn, WINNERS_BET_TYPE)
		if err != nil {
			return err
		}
	}

	for _, m := range msg {
		data := strings.Split(m, ";")
		if data[0] == ID {
			log.Infof("action: consulta_ganadores | result: success | client_id: %v | cant_ganadores: %v", ID, data[1])
		}
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
