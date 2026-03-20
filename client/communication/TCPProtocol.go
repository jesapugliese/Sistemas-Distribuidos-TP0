package communication

import "net"

type TCPProtocol struct {
	conn net.Conn
}

func NewTCPProtocol(conn net.Conn) TCPProtocol {
	return TCPProtocol{conn: conn}
}

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

func (tcpProt TCPProtocol) ReadAll(buffer []byte) (int, error) {
	// TODO
	return 0, nil
}
