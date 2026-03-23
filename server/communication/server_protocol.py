from communication.bet_serializer import BetSerializer
from communication.tcp_protocol import TCPProtocol


class ServerProtocol:
    def __init__(self):
        self._client_socket = None
        self._bet_serializer = BetSerializer()
        self._tcp_protocol = TCPProtocol()

    def recv_batch_bets_amount_msg(self) -> int:
        """
        Receives the batch bets amount message from the client socket and returns the batch 
        bets amount as an integer.
        """

        batch_bets_amount_bytes = self._tcp_protocol.recv_exact(self._client_socket, 
                                                                self._bet_serializer.get_batch_bets_amount_msg_length())
        return self._bet_serializer.deserialize_batch_bets_amount(batch_bets_amount_bytes)

    def recv_store_bets_batch_msg(self, batch_bets_amount):
        """
        Receives a batch of bets store message from the client socket and deserializes the bets into a list of 
        Bet objects.
        It returns the list of bets.
        """
        bets = []

        for _ in range(batch_bets_amount):
            bet_bytes = self._tcp_protocol.recv_with_length(self._client_socket)
            bet = self._bet_serializer.deserialize_bet(bet_bytes)
            bets.append(bet)

        return bets

    def send_store_bets_response(self, success):
        """
        Sends a response message to the client socket indicating whether the storage
        of the bets was successful or not.
        """

        msg_store_response_bytes = self._bet_serializer.serialize_store_bets_response(success)
        self._tcp_protocol.send_all(self._client_socket, msg_store_response_bytes)

    def update_client_socket(self, client_socket):
        self._client_socket = client_socket

    def close_client_socket(self):
        if self._client_socket:
            self._client_socket.close()
