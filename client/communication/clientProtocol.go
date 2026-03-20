package communication

import (
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type ClientProtocol struct {
	conn net.Conn
	// TODO: apuestaSerializer ApuestaSerializer
	// TODO: tcpProtocol TCPProtocol
}

func NewClientProtocol(conn net.Conn) ClientProtocol {
	return ClientProtocol{
		conn: conn,
	}
}

func (cp ClientProtocol) Send(apuesta utils.Apuesta) error {
	// TODO:
	// 1. Serializar la apuesta usando apuestaSerializer
	// 2. Enviar la apuesta serializada usando tcpProtocol
	return nil
}

func (cp ClientProtocol) Receive() (string, error) {
	// TODO
	// 1. Recibir la respuesta usando tcpProtocol
	// 2. Retornar la respuesta como string
	return "", nil
}
