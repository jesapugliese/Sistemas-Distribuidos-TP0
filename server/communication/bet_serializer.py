from common.utils import Bet

class BetSerializer:
    def __init__(self):
        self.batch_bets_amount_msg_length = 2

    def get_batch_bets_amount_msg_length(self):
        return self.batch_bets_amount_msg_length
    
    def get_identification_msg_length(self):
        return 1

    def serialize_store_bets_response(self, success) -> bytes:
        """
        Serializes a response message for storing a bet, indicating whether the operation was successful.
        The format of the serialized message is:
        - success: 1 byte (integer, 1 = success)
        """

        success_byte = b'\x01' if success else b'\x00'
        return success_byte

    def serialize_winners_notification(self, winners) -> bytes:
        """
        Serializes a notification message with the winners of the lottery.
        The winners are sent as a list of documents (integers).
        The format of the serialized message is:
        - winners_count: 2 bytes (integer)
        - winners_documents: for each winner, 4 bytes (integer) representing the document number
        """

        winners_count = len(winners)
        winners_count_bytes = winners_count.to_bytes(2, 'big')
        winners_documents_bytes = b''.join([winner.to_bytes(4, 'big') for winner in winners])
        return winners_count_bytes + winners_documents_bytes

    def deserialize_identification_msg(self, identification_msg_bytes) -> int:
        """
        Deserializes the given bytes into an integer representing the agency ID. 
        The expected format of the bytes is:
        - agency_id: 1 byte (integer)
        """

        return identification_msg_bytes[0]

    def deserialize_bet(self, bet_bytes) -> Bet:
        """
        Deserializes the given bytes into a Bet object. The expected format of the bytes is:
        - msg_len: 2 bytes (integer)
        - agency_id: 1 byte (integer)
        - first_name: variable length string (preceded by its 1-byte length)
        - last_name: variable length string (preceded by its 1-byte length)
        - document: 4 bytes (integer)
        - birthdate:
          > year: 2 bytes (integer)
          > month: 1 byte (integer)
          > day: 1 byte (integer)
        - number: 2 bytes (integer)
        """

        offset = 2 # Ignore msg_len

        # Deserialize agency ID
        agency_id = bet_bytes[offset]
        offset += 1

        # Deserialize first name
        first_name_length = bet_bytes[offset]
        offset += 1
        first_name = bet_bytes[offset:offset + first_name_length].decode('utf-8')
        offset += first_name_length

        # Deserialize last name
        last_name_length = bet_bytes[offset]
        offset += 1
        last_name = bet_bytes[offset:offset + last_name_length].decode('utf-8')
        offset += last_name_length

        # Deserialize document number
        document = int.from_bytes(bet_bytes[offset:offset + 4], 'big')
        offset += 4

        # Deserialize birthdate
        birthdate_year = int.from_bytes(bet_bytes[offset:offset + 2], 'big')
        offset += 2
        birthdate_month = bet_bytes[offset]
        offset += 1
        birthdate_day = bet_bytes[offset]
        offset += 1

        # Deserialize bet number
        number = int.from_bytes(bet_bytes[offset:offset + 2], 'big')
        offset += 2

        return Bet(
            agency=agency_id, 
            first_name=first_name, 
            last_name=last_name, 
            document=document,
            birthdate=f"{birthdate_year:04d}-{birthdate_month:02d}-{birthdate_day:02d}", 
            number=number
        )

    def deserialize_batch_bets_amount(self, batch_bets_amount_bytes) -> int:
        """
        Deserializes the given bytes into an integer representing the batch bets amount. 
        The expected format of the bytes is:
        - batch_bets_amount: 2 bytes (integer)
        """

        return int.from_bytes(batch_bets_amount_bytes, 'big')