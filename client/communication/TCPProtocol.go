package communication

import (
	"encoding/binary"
	"fmt"
	"net"
)

type TCPProtocol struct {
	conn net.Conn
}

func NewTCPProtocol(conn net.Conn) TCPProtocol {
	return TCPProtocol{conn: conn}
}

func (tcpProt TCPProtocol) SendAll(data []byte) error {
	// Send the length of the data first (2 bytes)
	var lengthBytes [2]byte

	length := uint16(len(data))
	binary.BigEndian.PutUint16(lengthBytes[:], length)

	n, err := tcpProt.conn.Write(lengthBytes[:])
	if err != nil {
		return err
	}
	if n != 2 {
		return fmt.Errorf("Failed to send data length")
	}

	// Send the total data
	totalSend := 0

	for totalSend < len(data) {
		n, err := tcpProt.conn.Write(data[totalSend:])
		if err != nil {
			return err
		}
		totalSend += n
	}

	return nil
}

func (tcpProt TCPProtocol) ReadAll(buffer []byte) (int, error) {
	// TODO
	return 0, nil
}
