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
		betSerializer: NewBetSerializer(),
		tcpProtocol:   NewTCPProtocol(clientSocket),
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

// Send serializes the given bet and sends it to the server using TCPProtocol.
func (cp ClientProtocol) SendBet(bet utils.Bet) error {
	betSerialized, err := cp.betSerializer.SerializeBet(bet)
	if err != nil {
		return err
	}
	err = cp.tcpProtocol.SendAll(betSerialized)
	if err != nil {
		return err
	}
	return nil
}

// RecvBetStoreResponse receives the response from the server after sending a bet.
func (cp ClientProtocol) RecvBetStoreResponse() (utils.BetStoreResponse, error) {
	betStoreResponseBytes, err := cp.tcpProtocol.RecvAll()
	if err != nil {
		return utils.BetStoreResponse{}, err
	}
	return cp.betSerializer.DeserializeBetStoreResponse(betStoreResponseBytes)
}

// Close closes the TCP connection.
func (cp ClientProtocol) Close() {
	cp.tcpProtocol.Close()
}

// BatchReachMaxSize checks if adding the current bet to the batch would
// exceed either the maximum number of bets allowed in a batch or the
// maximum packet size allowed for a batch.
func (cp ClientProtocol) BatchReachMaxSize(currentBet utils.Bet, batchMaxAmount int, batchMaxPacketSize int) bool {
	currentBetPacketSize := cp.betSerializer.CalculateBetPacketSize(currentBet)
	if cp.batchPacketSize+currentBetPacketSize > batchMaxPacketSize ||
		len(cp.batch) > batchMaxAmount {
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
	cp.batchPacketSize += cp.betSerializer.CalculateBetPacketSize(bet)
}
