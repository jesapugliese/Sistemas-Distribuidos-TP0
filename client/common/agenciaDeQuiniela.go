package common

import (
	"fmt"
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

// CreateBet creates a bet using the given parameters.
func (a *AgenciaDeQuiniela) CreateBet(firstName, lastName string,
	document string, birthdate string, number string) (utils.Bet, error) {
	if firstName == "" {
		return utils.Bet{}, fmt.Errorf("First name cannot be empty")
	}

	if lastName == "" {
		return utils.Bet{}, fmt.Errorf("Last name cannot be empty")
	}

	documentToInt, err := strconv.ParseUint(document, 10, 32)
	if err != nil {
		return utils.Bet{}, fmt.Errorf("Error at parsing document: %v", err)
	}
	if documentToInt > 99999999 {
		return utils.Bet{}, fmt.Errorf("Invalid document")
	}

	birthdateParsed, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return utils.Bet{}, fmt.Errorf("Invalid birthdate: %v", err)
	}
	birthdateDate := utils.Date{
		Year:  uint16(birthdateParsed.Year()),
		Month: uint8(birthdateParsed.Month()),
		Day:   uint8(birthdateParsed.Day()),
	}

	numberToInt, err := strconv.ParseUint(number, 10, 16)
	if err != nil {
		return utils.Bet{}, fmt.Errorf("Error at parsing bet number: %v", err)
	}
	if numberToInt > 9999 {
		return utils.Bet{}, fmt.Errorf("Invalid bet number")
	}

	return utils.Bet{
		AgencyID:  a.id,
		FirstName: firstName,
		LastName:  lastName,
		Document:  uint32(documentToInt),
		Birthdate: birthdateDate,
		Number:    uint16(numberToInt),
	}, nil
}

// Identificate sends the identification message (agency ID) to the server using the
// provided client protocol.
func (a *AgenciaDeQuiniela) Identificate(clientProtocol communication.ClientProtocol) error {
	return clientProtocol.SendIdentificationMsg(a.id)
}

// StoreBetsBatch sends the given batch of bets to the server using the provided client protocol.
func (a *AgenciaDeQuiniela) StoreBetsBatch(clientProtocol communication.ClientProtocol, batch []utils.Bet) error {
	return clientProtocol.SendStoreBetsBatchMsg(batch)
}

// GetStoreBetsResult receives the response from the server after sending a bet.
func (a *AgenciaDeQuiniela) GetStoreBetsResult(clientProtocol communication.ClientProtocol) (bool, error) {
	return clientProtocol.RecvStoreBetsResponse()
}

// GetWinnersNotification receives the winners notification from the server after the lottery draw.
func (a *AgenciaDeQuiniela) GetWinnersNotification(clientProtocol communication.ClientProtocol) ([]uint32, error) {
	return clientProtocol.RecvWinnersNotification()
}
