package communication

import (
	"encoding/binary"
	"net"
)

type TCPProtocol struct {
	conn net.Conn
}

func NewTCPProtocol(conn net.Conn) TCPProtocol {
	return TCPProtocol{conn: conn}
}

// SendAll sends the given data over the TCP connection.
func (tcpProt TCPProtocol) SendAll(data []byte) error {
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

// RecvAll receives data from the TCP connection. It first reads
// 2 bytes to determine the length of the incoming message, and
// then reads the message data based on that length.
// It returns the complete message, including the length bytes.
func (tcpProt TCPProtocol) RecvAll() ([]byte, error) {
	lengthBytes, err := tcpProt.RecvExact(2)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(lengthBytes)

	data, err := tcpProt.RecvExact(int(length))
	if err != nil {
		return nil, err
	}

	return append(lengthBytes, data...), nil
}

// RecvExact reads exactly n bytes from the TCP connection.
func (tcpProt TCPProtocol) RecvExact(n int) ([]byte, error) {
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

func (tcpProt TCPProtocol) Close() {
	tcpProt.conn.Close()
}
