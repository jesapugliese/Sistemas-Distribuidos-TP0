import socket


class TCPProtocol:
    def __init__(self):
        pass

    def send_all(self, conn: socket.socket, data: bytes):
        '''
        Sends all bytes to the TCP connection, first sending the length of the data
        (2 bytes) and then sending the data itself.
        '''

        if len(data) > 2**16 - 1:
            raise ValueError("Message too large")

        length = len(data)
        length_bytes = length.to_bytes(2, byteorder='big', signed=False)
        conn.sendall(length_bytes + data)

    def recv_all(self, conn: socket.socket):
        '''
        Receives all bytes from the TCP connection, first reading the length of the data 
        (2 bytes) and then reading the data itself.
        '''
        
        length_bytes = self._recv_exact(conn, 2)
        length = int.from_bytes(length_bytes, byteorder='big', signed=False)

        return self._recv_exact(conn, length)
    
    def _recv_exact(self, conn: socket.socket, n: int) -> bytes:
        data = b''
        while len(data) < n:
            chunk = conn.recv(n - len(data))
            if not chunk:
                raise ConnectionError("connection closed")
            data += chunk
        return data