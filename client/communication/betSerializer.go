package communication

import (
	"encoding/binary"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type BetSerializer struct {
}

const (
	OpcodeStoreBetMsg  = 1
	OpcodeBatchSizeMsg = 2
)

func NewBetSerializer() BetSerializer {
	return BetSerializer{}
}

// CalculateStoreBetMsgPacketSize calculates the size of the packet that would be generated
// by serializing the given store bet message.
func (as BetSerializer) CalculateStoreBetMsgPacketSize(bet utils.Bet) int {
	return 2 + 1 + 1 + 1 + len(bet.FirstName) + 1 + len(bet.LastName) + 4 + 2 + 2
}

// SerializeStoreBetMsg generates serializes the message to store a bet.
// Serialization format:
//   - MsgLen: 2 bytes (integer)
//   - Opcode: 1 byte (integer)
//   - AgencyID: 1 byte (integer)
//   - FirstName: variable length string (preceded by its 1-byte length)
//   - LastName: variable length string (preceded by its 1-byte length)
//   - Document: 4 bytes (integer)
//   - Birthdate:
//     > Year: 2 bytes (integer)
//     > Month: 1 byte (integer)
//     > Day: 1 byte (integer)
//   - Number: 2 bytes (integer)
func (as BetSerializer) SerializeStoreBetMsg(bet utils.Bet) ([]byte, error) {
	var serialized []byte

	// Serialize MsgLen
	var tmp2 [2]byte
	binary.BigEndian.PutUint16(tmp2[:], uint16(as.CalculateStoreBetMsgPacketSize(bet)))
	serialized = append(serialized, tmp2[:]...)

	// Serialize Opcode
	serialized = append(serialized, OpcodeStoreBetMsg)

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
	binary.BigEndian.PutUint16(tmp2[:], uint16(bet.Number))
	serialized = append(serialized, tmp2[:]...)

	return serialized, nil
}

// DeserializeBetStoreResponse deserializes the response received from the server after
// sending a request for storage of a bet. It extracts the success status, document,
// and number from the byte slice.
// Data format:
//   - MsgLen: 2 bytes (integer)
//   - Success: 1 byte (value 1 for success)
//   - Document: 4 bytes (integer)
//   - Number: 2 bytes (integer)
func (as BetSerializer) DeserializeBetStoreResponse(betStoreResponseBytes []byte) (utils.BetStoreResponse, error) {
	if len(betStoreResponseBytes) != 8 {
		return utils.BetStoreResponse{}, fmt.Errorf("Invalid data length for bet store response")
	}

	msgLen := binary.BigEndian.Uint16(betStoreResponseBytes[0:2])
	if msgLen != 6 {
		return utils.BetStoreResponse{}, fmt.Errorf("Invalid MsgLen in bet store response: expected 6, got %d", msgLen)
	}
	success := betStoreResponseBytes[2]
	document := binary.BigEndian.Uint32(betStoreResponseBytes[3:7])
	number := binary.BigEndian.Uint16(betStoreResponseBytes[7:9])

	return utils.BetStoreResponse{
		Success:  success == 1,
		Document: document,
		Number:   number,
	}, nil
}

// SerializeBatchSize serialzes the message that indicates the batch size.
// Serialization format:
//   - MsgLen: 2 bytes (integer)
//   - Opcode: 1 byte (integer)
//   - BatchSize: 2 bytes (integer)
func (as BetSerializer) SerializeBatchSize(batchSize int) []byte {
	var serialized []byte

	// Serialize MsgLen
	var tmp2 [2]byte
	msgLen := 1 + 2
	binary.BigEndian.PutUint16(tmp2[:], uint16(msgLen))
	serialized = append(serialized, tmp2[:]...)

	// Serialize Opcode
	serialized = append(serialized, OpcodeBatchSizeMsg)

	// Serialize BatchSize
	binary.BigEndian.PutUint16(tmp2[:], uint16(batchSize))
	serialized = append(serialized, tmp2[:]...)

	return serialized
}
