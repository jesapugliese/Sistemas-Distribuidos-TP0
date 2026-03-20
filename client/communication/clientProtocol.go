package communication

import (
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type ClientProtocol struct {
	apuestaSerializer ApuestaSerializer
	tcpProtocol       TCPProtocol
}

func NewClientProtocol(conn net.Conn) ClientProtocol {
	return ClientProtocol{
		apuestaSerializer: NewApuestaSerializer(),
		tcpProtocol:       NewTCPProtocol(conn),
	}
}

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
