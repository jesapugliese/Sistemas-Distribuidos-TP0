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

// CreateClientSocket Initializes client socket. In case of
// failure, return error.
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
// If any error occurs during serialization or sending, it is returned to the caller.
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

func (cp ClientProtocol) RecvResult() (string, error) {
	// TODO
	// 1. Recibir la respuesta usando tcpProtocol
	// 2. Retornar la respuesta como string
	return "", nil
}
