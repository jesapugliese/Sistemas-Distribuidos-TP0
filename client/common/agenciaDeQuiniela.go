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
	name string
	id   uint8
}

func NewAgenciaDeQuiniela(name string, id uint8) *AgenciaDeQuiniela {
	return &AgenciaDeQuiniela{
		name: name,
		id:   id,
	}
}

// CreateBet reads the bet information from environment variables,
// validates it and creates an Bet struct.
func (a *AgenciaDeQuiniela) CreateBet() (utils.Bet, error) {
	firstName := os.Getenv("CLI_NOMBRE")

	lastName := os.Getenv("CLI_APELLIDO")

	document, err := strconv.ParseUint(os.Getenv("CLI_DOCUMENTO"), 10, 32)
	if err != nil {
		return utils.Bet{}, err
	}
	if document <= 0 || document > 99999999 {
		return utils.Bet{}, fmt.Errorf("Invalid document")
	}

	birthdateString := os.Getenv("CLI_NACIMIENTO")
	birthdateParsed, err := time.Parse("2006-01-02", birthdateString)
	if err != nil {
		return utils.Bet{}, fmt.Errorf("Invalid birthdate")
	}
	birthdate := utils.Date{
		Year:  uint16(birthdateParsed.Year()),
		Month: uint8(birthdateParsed.Month()),
		Day:   uint8(birthdateParsed.Day()),
	}

	number, err := strconv.ParseUint(os.Getenv("CLI_NUMERO"), 10, 16)
	if err != nil {
		return utils.Bet{}, err
	}
	if number <= 0 || number > 9999 {
		return utils.Bet{}, fmt.Errorf("Invalid bet number")
	}

	return utils.Bet{
		AgencyID:  a.id,
		FirstName: firstName,
		LastName:  lastName,
		Document:  uint32(document),
		Birthdate: birthdate,
		Number:    uint16(number),
	}, nil
}

// StoreBet creates a bet using the CreateBet method and sends it to the server
// using the given ClientProtocol.
func (a *AgenciaDeQuiniela) StoreBet(clientProtocol communication.ClientProtocol) error {
	bet, err := a.CreateBet()
	if err != nil {
		return err
	}
	clientProtocol.SendBet(bet)
	return nil
}

// RecvBetStoreResponse receives the response from the server after sending a bet.
func (a *AgenciaDeQuiniela) RecvBetStoreResponse(clientProtocol communication.ClientProtocol) (utils.BetStoreResponse, error) {
	return clientProtocol.RecvBetStoreResponse()
}
