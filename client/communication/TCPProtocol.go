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

// SendAll sends the given data over the TCP connection. It first
// sends the length of the data (2 bytes) followed by the actual data.
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

// RecvAll receives data from the TCP connection. It first reads the length
// of the incoming data (2 bytes) and then reads the actual data based on
// that length. It returns the received data as a byte slice.
func (tcpProt TCPProtocol) RecvAll() ([]byte, error) {
	lengthBytes, err := tcpProt.recvExact(2)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(lengthBytes)

	return tcpProt.recvExact(int(length))
}

// recvExact reads exactly n bytes from the TCP connection.
func (tcpProt TCPProtocol) recvExact(n int) ([]byte, error) {
	data := make([]byte, n)
	totalRead := 0
	for totalRead < n {
		bytesRead, err := tcpProt.conn.Read(data[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += bytesRead
	}
	return data, nil
}
