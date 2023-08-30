package common

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

	return SINGLE_BET_TYPE + ";" + a.Agency + ";" + apuesta.Name + ";" + apuesta.LastName + ";" + apuesta.Document + ";" + apuesta.Birthday + ";" + apuesta.Number + "\n"
}
