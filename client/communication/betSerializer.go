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

// CalculateBetPacketSize calculates the size of the packet that would be generated
// by serializing the given bet.
func (as BetSerializer) CalculateBetPacketSize(bet utils.Bet) int {
	return 1 + 1 + 1 + len(bet.FirstName) + 1 + len(bet.LastName) + 4 + 2 + 2
}

// Serialize converts a Bet struct into a byte slice that can be sent
// over the network.
// Serialization format:
//   - MsgLen: 1 byte (integer)
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

	// Serialize MsgLen
	msgLen := as.CalculateBetPacketSize(bet)
	if msgLen > 255 {
		return nil, fmt.Errorf("Bet data too large to serialize")
	}
	serialized = append(serialized, byte(msgLen))

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
//   - MsgLen: 1 byte (integer)
//   - Success: 1 byte (value 1 for success)
//   - Document: 4 bytes (integer)
//   - Number: 2 bytes (integer)
func (as BetSerializer) DeserializeBetStoreResponse(betStoreResponseBytes []byte) (utils.BetStoreResponse, error) {
	if len(betStoreResponseBytes) != 7 {
		return utils.BetStoreResponse{}, fmt.Errorf("Invalid data length for bet store response")
	}

	msgLen := betStoreResponseBytes[0]
	if msgLen != 6 {
		return utils.BetStoreResponse{}, fmt.Errorf("Invalid MsgLen in bet store response: expected 6, got %d", msgLen)
	}
	success := betStoreResponseBytes[1]
	document := binary.BigEndian.Uint32(betStoreResponseBytes[2:6])
	number := binary.BigEndian.Uint16(betStoreResponseBytes[6:8])

	return utils.BetStoreResponse{
		Success:  success == 1,
		Document: document,
		Number:   number,
	}, nil
}
