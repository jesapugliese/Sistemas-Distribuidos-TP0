from common.utils import Bet


class BetSerializer:
    def __init__(self):
        pass

    def serialize_store_success(self, document, number) -> bytes:
        '''
        Serializes a success message for storing a bet, including the document and number of the bet.
        The format of the serialized message is:
        - document: 4 bytes (integer)
        - number: 2 bytes (integer)
        '''

        document_bytes = document.to_bytes(4, 'big')
        number_bytes = number.to_bytes(2, 'big')

        return document_bytes + number_bytes

    def deserialize_bet(self, bet_bytes) -> Bet:
        '''
        Deserializes the given bytes into a Bet object. The expected format of the bytes is:
        - agency_id: 1 byte (integer)
        - first_name: variable length string (preceded by its 1-byte length)
        - last_name: variable length string (preceded by its 1-byte length)
        - document: 4 bytes (integer)
        - birthdate:
          > year: 2 bytes (integer)
          > month: 1 byte (integer)
          > day: 1 byte (integer)
        - number: 2 bytes (integer)
        '''

        offset = 0

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
