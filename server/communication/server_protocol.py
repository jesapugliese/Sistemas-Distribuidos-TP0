from communication.bet_serializer import BetSerializer
from communication.tcp_protocol import TCPProtocol


class ServerProtocol:
    def __init__(self):
        self._bet_serializer = BetSerializer()
        self._tcp_protocol = TCPProtocol()

    def recv_agency_id_msg(self, client_socket) -> int:
        """
        Receives the identification message from the client socket and returns the agency ID as an integer.
        """

        identification_msg_bytes = self._tcp_protocol.recv_exact(client_socket, 
                                                                self._bet_serializer.get_identification_msg_length())
        return self._bet_serializer.deserialize_identification_msg(identification_msg_bytes)

    def recv_batch_bets_amount_msg(self, client_socket) -> int:
        """
        Receives the batch bets amount message from the client socket and returns the batch 
        bets amount as an integer.
        """

        batch_bets_amount_bytes = self._tcp_protocol.recv_exact(client_socket, 
                                                                self._bet_serializer.get_batch_bets_amount_msg_length())
        return self._bet_serializer.deserialize_batch_bets_amount(batch_bets_amount_bytes)

    def recv_store_bets_batch_msg(self, client_socket, batch_bets_amount):
        """
        Receives a batch of bets store message from the client socket and deserializes the bets into a list of 
        Bet objects.
        It returns the list of bets.
        """
        bets = []

        for _ in range(batch_bets_amount):
            bet_bytes = self._tcp_protocol.recv_with_length(client_socket)
            bet = self._bet_serializer.deserialize_bet(bet_bytes)
            bets.append(bet)

        return bets

    def send_store_bets_response(self, client_socket, success):
        """
        Sends a response message to the client socket indicating whether the storage
        of the bets was successful or not.
        """

        msg_store_response_bytes = self._bet_serializer.serialize_store_bets_response(success)
        self._tcp_protocol.send_all(client_socket, msg_store_response_bytes)

    def send_winners_notification(self, client_socket, winners):
        """
        Sends a notification message to the client socket with the winners of the lottery.
        The winners are sent as a list of documents (integers).
        """

        winners_bytes = self._bet_serializer.serialize_winners_notification(winners)
        self._tcp_protocol.send_all(client_socket, winners_bytes)
