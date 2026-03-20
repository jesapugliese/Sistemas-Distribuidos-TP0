package communication

import (
	"fmt"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type ClientProtocol struct {
	apuestaSerializer ApuestaSerializer
	tcpProtocol       TCPProtocol
}

func NewClientProtocol(serverAddress string, clientID string) (ClientProtocol, error) {
	clientSocket, err := createClientSocket(serverAddress, clientID)
	if err != nil {
		return ClientProtocol{}, err
	}
	return ClientProtocol{
		apuestaSerializer: NewApuestaSerializer(),
		tcpProtocol:       NewTCPProtocol(clientSocket),
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

// Send serializes the given Apuesta and sends it to the server using TCPProtocol.
// If any error occurs during serialization or sending, it is returned to the caller.
func (cp ClientProtocol) Send(apuesta utils.Apuesta) error {
	apuestaSerialized, err := cp.apuestaSerializer.Serialize(apuesta)
	if err != nil {
		return err
	}
	err = cp.tcpProtocol.SendAll(apuestaSerialized)
	if err != nil {
		return err
	}
	return nil
}

func (cp ClientProtocol) Receive() (string, error) {
	// TODO
	// 1. Recibir la respuesta usando tcpProtocol
	// 2. Retornar la respuesta como string
	return "", nil
}
