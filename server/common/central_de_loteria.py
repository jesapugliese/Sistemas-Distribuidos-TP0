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

        success = True

        bets = server_protocol.recv_store_bets_batch_msg(batch_bets_amount)

        try:
            store_bets(bets)
        except Exception:
            success = False
        
        server_protocol.send_store_bets_response(success)

        return None
