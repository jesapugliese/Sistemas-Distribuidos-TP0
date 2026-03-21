from utils import BetStoreResponse, store_bets


class CentralDeLoteriaNacional:
    def __init__(self):
        pass
    
    def recv_msg_store_bets(self, server_protocol):
        """
        Receive a bet from a client and store it in the system. Then sends a success 
        message back to the client with the document and number of the bet.
        """
        
        bet = server_protocol.recv_bet()
        store_bets([bet])
        bet_store_response = BetStoreResponse(success=True, document=bet.document, number=bet.number)
        server_protocol.send_store_success(bet_store_response)

        return bet.document, bet.number
