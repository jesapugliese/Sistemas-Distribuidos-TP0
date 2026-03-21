package communication

import (
	"fmt"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type ClientProtocol struct {
	betSerializer BetSerializer
	tcpProtocol   TCPProtocol
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
