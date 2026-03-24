package communication

import (
	"fmt"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type ClientProtocol struct {
	betSerializer   BetSerializer
	tcpProtocol     TCPProtocol
	batch           []utils.Bet
	batchPacketSize int
}

func NewClientProtocol(serverAddress string, clientID string) (ClientProtocol, error) {
	clientSocket, err := createClientSocket(serverAddress, clientID)
	if err != nil {
		return ClientProtocol{}, err
	}
	return ClientProtocol{
		betSerializer:   NewBetSerializer(),
		tcpProtocol:     NewTCPProtocol(clientSocket),
		batch:           nil,
		batchPacketSize: 0,
	}, nil
}

// CreateClientSocket Initializes client socket.
func createClientSocket(serverAddress string, clientID string) (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf(
			"action: connect | result: fail | client_id: %v | error: %v",
			clientID,
			err,
		)
	}
	return conn, nil
}

// SendIdentificationMsg sends the identification message (agency ID) to the server
// using the TCP protocol.
func (cp ClientProtocol) SendIdentificationMsg(agencyID uint8) error {
	identificationMsgBytes := cp.betSerializer.SerializeIdentificationMsg(agencyID)
	return cp.tcpProtocol.SendAll(identificationMsgBytes)
}

// SendStoreBetsBatchMsg sends the given batch of bets to the server using the TCP protocol.
func (cp ClientProtocol) SendStoreBetsBatchMsg(batch []utils.Bet) error {
	storeBetsBatchMsgBytes, err := cp.betSerializer.SerializeStoreBetsBatchMsg(batch)
	if err != nil {
		return err
	}
	return cp.tcpProtocol.SendAll(storeBetsBatchMsgBytes)
}

// SendBatchBetsAmountMsg sends a message to the server indicating the size of the
// batch of bets that will be sent next.
func (cp ClientProtocol) SendBatchBetsAmountMsg(batchBetsAmount int) error {
	batchBetsAmountMsgBytes := cp.betSerializer.SerializeBatchBetsAmount(batchBetsAmount)
	return cp.tcpProtocol.SendAll(batchBetsAmountMsgBytes)
}

// RecvStoreBetsResponse receives the response from the server after sending
// a batch of bets. It returns whether the batch was successfully stored.
func (cp ClientProtocol) RecvStoreBetsResponse() (bool, error) {
	storeBetsResponseBytes, err := cp.tcpProtocol.RecvExact(cp.betSerializer.CalculateStoreBetsResponsePacketSize())
	if err != nil {
		return false, err
	}
	return cp.betSerializer.DeserializeStoreBetsResponse(storeBetsResponseBytes)
}

// Close closes the TCP connection.
func (cp ClientProtocol) Close() {
	cp.tcpProtocol.Close()
}

// BatchReachMaxSize checks if adding the current bet to the batch would
// exceed either the maximum number of bets allowed in a batch or the
// maximum packet size allowed for a batch.
func (cp ClientProtocol) BatchReachMaxSize(currentBet utils.Bet, batchMaxAmount int, batchMaxPacketSize int) bool {
	currentBetPacketSize := cp.betSerializer.CalculateStoreBetMsgPacketSize(currentBet)
	if cp.batchPacketSize+currentBetPacketSize > batchMaxPacketSize ||
		len(cp.batch)+1 > batchMaxAmount {
		return true
	}
	return false
}

// GetBatch returns the current batch of bets and resets the batch and
// its packet size.
func (cp *ClientProtocol) GetBatch() []utils.Bet {
	batch := cp.batch
	cp.batch = nil
	cp.batchPacketSize = 0
	return batch
}

// AppendToBatch adds the given bet to the current batch and updates
// the batch packet size accordingly.
func (cp *ClientProtocol) AppendToBatch(bet utils.Bet) {
	cp.batch = append(cp.batch, bet)
	cp.batchPacketSize += cp.betSerializer.CalculateStoreBetMsgPacketSize(bet)
}
