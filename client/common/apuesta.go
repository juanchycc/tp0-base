package common

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const SINGLE_BET_TYPE = "SINGLE_BET"
const MULTIPLE_BET_TYPE = "MULTIPLE_BET"
const FINISH_BET_TYPE = "FINISH_BET"
const SUCCESS_BET_TYPE = "SUCCESS_BET"
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

func leerApuestas(ID string, conn net.Conn) error {

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
	//TODO: finish
	_, err = fmt.Fprintf(
		conn,
		"%v;%v;%v\n", FINISH_BET_TYPE, 0, ID,
	)

	return err
}

func sendBets(apuestas string, conn net.Conn, ID string) error {

	totalLen := len(apuestas)
	header := MULTIPLE_BET_TYPE + ";" + strconv.Itoa(totalLen) + ";" + ID + "\n"
	headerLen := len(header)

	if totalLen > MAX_BUFFER {

		totalWrite := 0
		for totalWrite < totalLen {
			end := totalWrite + MAX_BUFFER
			packet := MULTIPLE_BET_TYPE + ";" + strconv.Itoa(totalLen) + ";" + ID + "\n" + apuestas[totalWrite:end] + "\n"
			writeLen, err := conn.Write([]byte(packet))
			if err != nil {
				return err
			}
			totalWrite += writeLen - headerLen
		}
	} else {
		//log.Infof("Envio: %v %v", header+msg, totalLen+headerLen)
		packet := header + apuestas
		_, err := fmt.Fprint(conn, packet)
		if err != nil {
			return err
		}
		/*TODO: No se envio todo entonces reenvio lo que falta
		if (writeLen - headerLen) != totalLen {

			return sendBets(MULTIPLE_BET_TYPE+";"+strconv.Itoa(totalLen)+";"+ID+"\n"+msg[writeLen:totalLen-1], conn, ID)
		}*/
	}

	return waitSuccess(ID, conn)

}

func waitSuccess(ID string, conn net.Conn) error {

	msg, err := bufio.NewReaderSize(conn, MAX_BUFFER).ReadString('\n')

	if err != nil {
		return err
	}

	res := strings.Split(msg, ";")
	if res[TYPE_MSG_POSITION] != SUCCESS_BET_TYPE {
		return waitSuccess(ID, conn)
	}

	return nil
}
