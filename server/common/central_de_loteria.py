from common.utils import store_bets


class CentralDeLoteriaNacional:
    def __init__(self):
        pass
    
    def store_bets(self, server_protocol, batch_size):
        """
        Receive a batch of bets from the client through the server protocol, 
        store them using the store_bets function, and send a response back 
        to the client indicating whether the storage was successful or not.
        """
        bets = []

        bytes_read = 0
        bets_amount = 0
        while bytes_read < batch_size:
            bet, msg_size = server_protocol.recv_store_bet_msg()
            bets.append(bet)
            bytes_read += msg_size
            bets_amount += 1

        try:
            store_bets(bets)
        except Exception as e:
            return bets_amount, e
        
        server_protocol.send_store_bets_response(False if e else True)

        return bets_amount, None
