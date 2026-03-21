package communication

import (
	"encoding/binary"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type BetSerializer struct {
}

func NewBetSerializer() BetSerializer {
	return BetSerializer{}
}

// Serialize converts a Bet struct into a byte slice that can be sent
// over the network.
// Serialization format:
//   - AgencyID: 1 byte (integer)
//   - FirstName: variable length string (preceded by its 1-byte length)
//   - LastName: variable length string (preceded by its 1-byte length)
//   - Document: 4 bytes (integer)
//   - Birthdate:
//     > Year: 2 bytes (integer)
//     > Month: 1 byte (integer)
//     > Day: 1 byte (integer)
//   - Number: 2 bytes (integer)
func (as BetSerializer) SerializeBet(bet utils.Bet) ([]byte, error) {
	var serialized []byte

	// Serialize AgencyID
	serialized = append(serialized, byte(bet.AgencyID))

	// Serialize FirstName
	firstNameBytes := []byte(bet.FirstName)
	if len(firstNameBytes) > 255 {
		return nil, fmt.Errorf("FirstName too long")
	}
	serialized = append(serialized, byte(len(firstNameBytes)))
	serialized = append(serialized, firstNameBytes...)

	// Serialize LastName
	lastNameBytes := []byte(bet.LastName)
	if len(lastNameBytes) > 255 {
		return nil, fmt.Errorf("LastName too long")
	}
	serialized = append(serialized, byte(len(lastNameBytes)))
	serialized = append(serialized, lastNameBytes...)

	// Serialize Document
	var tmp4 [4]byte
	binary.BigEndian.PutUint32(tmp4[:], bet.Document)
	serialized = append(serialized, tmp4[:]...)

	// Serialize Birthdate
	binary.BigEndian.PutUint16(tmp4[0:2], bet.Birthdate.Year)
	tmp4[2] = bet.Birthdate.Month
	tmp4[3] = bet.Birthdate.Day
	serialized = append(serialized, tmp4[:]...)

	// Serialize Number
	var tmp2 [2]byte
	binary.BigEndian.PutUint16(tmp2[:], uint16(bet.Number))
	serialized = append(serialized, tmp2[:]...)

	return serialized, nil
}

// DeserializeBetStoreResponse deserializes the response received from the server after
// sending a request for storage of a bet. It extracts the success status, document,
// and number from the byte slice.
// Data format:
// - success: 1 byte (value 1 for success)
// - document: 4 bytes (integer)
// - number: 2 bytes (integer)
func (as BetSerializer) DeserializeBetStoreResponse(betStoreResponseBytes []byte) (utils.BetStoreResponse, error) {
	if len(betStoreResponseBytes) != 7 {
		return utils.BetStoreResponse{}, fmt.Errorf("Invalid data length for bet store response")
	}

	success := betStoreResponseBytes[0]
	document := binary.BigEndian.Uint32(betStoreResponseBytes[1:5])
	number := binary.BigEndian.Uint16(betStoreResponseBytes[5:7])

	return utils.BetStoreResponse{
		Success:  success == 1,
		Document: document,
		Number:   number,
	}, nil
}
