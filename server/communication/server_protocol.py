from communication.bet_serializer import BetSerializer
from communication.tcp_protocol import TCPProtocol
from common.utils import Bet


class ServerProtocol:
    def __init__(self):
        self._client_socket = None
        self._bet_serializer = BetSerializer()
        self._tcp_protocol = TCPProtocol()

    def recv_bet(self) -> Bet:
        """
        Receives a bet from the client socket and deserializes it into a 
        Bet object.
        """

        bet_bytes = self._tcp_protocol.recv_all(self._client_socket)
        return self._bet_serializer.deserialize_bet(bet_bytes)

    def send_store_success(self, bet_store_response):
        """
        Sends a success message to the client socket indicating that the bet was
        stored successfully. The message includes the document and number of the bet.
        """

        msg_store_success_bytes = self._bet_serializer.serialize_store_success(bet_store_response)
        self._tcp_protocol.send_all(self._client_socket, msg_store_success_bytes)

    def update_client_socket(self, client_socket):
        self._client_socket = client_socket

    def close(self):
        if self._client_socket:
            self._client_socket.close()
