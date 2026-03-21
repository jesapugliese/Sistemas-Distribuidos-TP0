from utils import store_bets


class CentralDeLoteriaNacional:
    def __init__(self):
        pass
    
    def recv_msg_store_bets(self, server_protocol):
        """
        Receive a bet from a client and store it in the system.
        """
        
        bet = server_protocol.recv_bet()
        store_bets([bet])
        server_protocol.send_store_success(bet.document, bet.number)

        return bet.document, bet.number
