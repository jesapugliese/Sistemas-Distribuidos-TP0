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

// CalculateStoreBetMsgPacketSize calculates the size of the packet that would be generated
// by serializing the given store bet message.
func (as BetSerializer) CalculateStoreBetMsgPacketSize(bet utils.Bet) int {
	return 2 + 1 + 1 + len(bet.FirstName) + 1 + len(bet.LastName) + 4 + 2 + 1 + 1 + 2
}

// CalculateStoreBetsResponsePacketSize calculates the size of the packet that contains
// the response from the server after sending a batch of bets.
func (as BetSerializer) CalculateStoreBetsResponsePacketSize() int {
	return 1
}

// SerializeStoreBetMsg generates serializes the message to store a bet.
// Serialization format:
//   - MsgLen: 2 bytes (integer)
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
	msgLen := as.CalculateStoreBetMsgPacketSize(bet) - 2
	binary.BigEndian.PutUint16(tmp2[:], uint16(msgLen))
	serialized = append(serialized, tmp2[:]...)

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

// SerializeStoreBetsBatchMsg generates the serialized message to store a batch of bets.
func (as BetSerializer) SerializeStoreBetsBatchMsg(batch []utils.Bet) ([]byte, error) {
	var serialized []byte

	for _, bet := range batch {
		betBytes, err := as.SerializeStoreBetMsg(bet)
		if err != nil {
			return nil, err
		}
		serialized = append(serialized, betBytes...)
	}

	return serialized, nil
}

// DeserializeStoreBetsResponse deserializes the response received from the server after
// sending a batch of bets.
// Deserialization format:
//   - Success: 1 byte (integer, 1 = success)
func (as BetSerializer) DeserializeStoreBetsResponse(storeBetsResponseBytes []byte) (bool, error) {
	if len(storeBetsResponseBytes) != 1 {
		return false, fmt.Errorf("Invalid data length for store bets response")
	}

	success := storeBetsResponseBytes[0]

	return success == 1, nil
}

// SerializeBatchBetsAmount serializes the message that indicates the batch bets amount.
// Serialization format:
//   - BatchBetsAmount: 2 bytes (integer)
func (as BetSerializer) SerializeBatchBetsAmount(batchBetsAmount int) []byte {
	var serialized []byte
	var tmp2 [2]byte

	binary.BigEndian.PutUint16(tmp2[:], uint16(batchBetsAmount))
	serialized = append(serialized, tmp2[:]...)

	return serialized
}
