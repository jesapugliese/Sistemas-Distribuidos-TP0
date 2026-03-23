from common.utils import store_bets


class CentralDeLoteriaNacional:
    def __init__(self):
        pass
    
    def store_bets(self, server_protocol, batch_bets_amount):
        """
        Receive a batch of bets from the client through the server protocol, 
        store them using the store_bets function, and send a response back 
        to the client indicating whether the storage was successful or not.
        """
        bets = []

        bets_processed = 0
        while bets_processed < batch_bets_amount:
            bet = server_protocol.recv_store_bet_msg()
            bets.append(bet)
            bets_processed += 1
        success = True

        try:
            store_bets(bets)
        except Exception:
            success = False
        
        server_protocol.send_store_bets_response(success)

        return None
