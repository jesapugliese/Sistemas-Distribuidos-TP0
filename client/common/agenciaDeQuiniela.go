package common

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/communication"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type AgenciaDeQuiniela struct {
	nombre string
	id     string
}

// NewAgenciaDeQuiniela initializes a new AgenciaDeQuiniela struct with
// the given name and id.
func NewAgenciaDeQuiniela(nombre string, id string) *AgenciaDeQuiniela {
	return &AgenciaDeQuiniela{
		nombre: nombre,
		id:     id,
	}
}

// CrearApuesta reads the bet information from environment variables,
// validates it and creates an Apuesta struct. If any of the fields is
// invalid an error is returned.
func (a *AgenciaDeQuiniela) CrearApuesta() (utils.Apuesta, error) {
	nombre := os.Getenv("CLI_NOMBRE")

	apellido := os.Getenv("CLI_APELLIDO")

	documento, err := strconv.ParseUint(os.Getenv("CLI_DOCUMENTO"), 10, 32)
	if err != nil {
		return utils.Apuesta{}, err
	}
	if documento <= 0 || documento > 99999999 {
		return utils.Apuesta{}, fmt.Errorf("Documento inválido")
	}

	nacimientoString := os.Getenv("CLI_NACIMIENTO")
	nacimientoParsed, err := time.Parse("2006-01-02", nacimientoString)
	if err != nil {
		return utils.Apuesta{}, fmt.Errorf("Fecha de nacimiento inválida")
	}
	nacimiento := utils.Fecha{
		Anio: uint16(nacimientoParsed.Year()),
		Mes:  uint8(nacimientoParsed.Month()),
		Dia:  uint8(nacimientoParsed.Day()),
	}

	numero, err := strconv.ParseUint(os.Getenv("CLI_NUMERO"), 10, 16)
	if err != nil {
		return utils.Apuesta{}, err
	}
	if numero <= 0 || numero > 9999 {
		return utils.Apuesta{}, fmt.Errorf("Número inválido")
	}

	return utils.Apuesta{
		Nombre:     nombre,
		Apellido:   apellido,
		Documento:  uint32(documento),
		Nacimiento: nacimiento,
		Numero:     uint16(numero),
	}, nil
}

func (a *AgenciaDeQuiniela) RegistrarApuesta(clientProtocol communication.ClientProtocol) error {
	apuesta, err := a.CrearApuesta()
	if err != nil {
		return err
	}
	clientProtocol.Send(apuesta)
	return nil
}

func (a *AgenciaDeQuiniela) RecibirResultado(clientProtocol communication.ClientProtocol) (string, error) {
	return clientProtocol.Receive()
}
