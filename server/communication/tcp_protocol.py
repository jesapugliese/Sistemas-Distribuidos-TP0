import socket


class TCPProtocol:
    def __init__(self):
        pass

    def send_all(self, conn: socket.socket, data: bytes):
        """
        Sends all bytes to the TCP connection.
        """
        conn.sendall(data)

    def recv_with_length(self, conn: socket.socket):
        """
        Receives all bytes from the TCP connection, first reading the length of the data 
        (2 bytes) and then reading the data itself.
        """
        length_bytes = self.recv_exact(conn, 2)
        length = int.from_bytes(length_bytes, byteorder='big', signed=False)

        data = self.recv_exact(conn, length)

        return length_bytes + data

    def recv_exact(self, conn: socket.socket, n: int) -> bytes:
        """
        Receives exactly n bytes from the TCP connection. If the connection is closed 
        before n bytes are received, raises a ConnectionError.
        """
        data = b''
        while len(data) < n:
            chunk = conn.recv(n - len(data))
            if not chunk:
                raise ConnectionError("connection closed")
            data += chunk
        return data