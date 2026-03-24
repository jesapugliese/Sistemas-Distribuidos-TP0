from common.utils import store_bets, load_bets, has_won
import threading


class CentralDeLoteriaNacional:
    def __init__(self):
        self._winners = {}

        self._lock_bets_storage = threading.Lock()
    
    def store_bets(self, server_protocol, client_socket, batch_bets_amount):
        """
        Receive a batch of bets from the client through the server protocol, 
        store them using the store_bets function, and send a response back 
        to the client indicating whether the storage was successful or not.
        """
        success = True

        bets = server_protocol.recv_store_bets_batch_msg(client_socket, batch_bets_amount)

        with self._lock_bets_storage:
            try:
                store_bets(bets)
            except Exception:
                success = False

        server_protocol.send_store_bets_response(client_socket, success)

        return success

    def draw_winners(self):
        """
        Load all the bets using the load_bets function, determine the winners using 
        the has_won function, and store the winners in the _winners attribute.
        """

        bets = load_bets()

        for bet in bets:
            if has_won(bet):
                if bet.agency not in self._winners:
                    self._winners[bet.agency] = []
                self._winners[bet.agency].append(int(bet.document))
        
    def notify_winners_to_agencies(self, server_protocol, clients):
        """
        Notify the winners to the agencies through the server protocol.
        """

        for agency_id, client_socket in clients:
            agency_winners = []
            if agency_id in self._winners:
                agency_winners = self._winners[agency_id]
            server_protocol.send_winners_notification(client_socket, agency_winners)
