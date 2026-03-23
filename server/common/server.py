import logging
import os
import signal
import socket

from configparser import ConfigParser
from common.central_de_loteria import CentralDeLoteriaNacional
from communication.server_protocol import ServerProtocol


class Server:
    def __init__(self):
        config_params = self._initialize_config()
        logging_level = config_params["logging_level"]
        port = config_params["port"]
        listen_backlog = config_params["listen_backlog"]

        self._initialize_log(logging_level)

        logging.debug(f"action: config | result: success | port: {port} | "
                      f"listen_backlog: {listen_backlog} | "
                      f"logging_level: {logging_level}")

        self._running = True
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(0.5)

        self._server_protocol = ServerProtocol()
        self._central_de_loteria = CentralDeLoteriaNacional()

        signal.signal(signal.SIGTERM, self._sigterm_handler)

    def start(self):
        """
        Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again.
        """

        while self._running:
            client_socket = self._accept_new_connection()
            if not client_socket:
                continue
            self._server_protocol.update_client_socket(client_socket)
            self._handle_client_connection()

    def _handle_client_connection(self):
        """
        Receives messages batches of bets from the client and processes them until 
        a batch size message with a non-positive batch size is received, 
        indicating the end of the communication. 
        """

        try:
            while True:
                batch_bets_amount = self._server_protocol.recv_batch_bets_amount_msg()
                logging.info(f"action: recibir_cantidad_apuestas_en_batch | result: success | cantidad_apuestas_en_batch: {batch_bets_amount}")
                if batch_bets_amount <= 0:
                    break
                err = self._central_de_loteria.store_bets(self._server_protocol, batch_bets_amount)
                if err:
                    logging.error(f"action: apuesta_recibida | result: fail | cantidad: {batch_bets_amount}")
                else:
                    logging.info(f"action: apuesta_recibida | result: success | cantidad: {batch_bets_amount}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            logging.info('action: close_client_connection | result: in_progress')
            self._server_protocol.close_client_socket()
            logging.info('action: close_client_connection | result: success')

    def _accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        try:
            logging.info('action: accept_connections | result: in_progress')
            client_socket, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        except socket.timeout:
            return None
        except OSError:
            return None
        
        return client_socket

    def _initialize_config(self):
        """ Parse env variables or config file to find program config params

        Function that search and parse program configuration parameters in the
        program environment variables first and the in a config file. 
        If at least one of the config parameters is not found a KeyError exception 
        is thrown. If a parameter could not be parsed, a ValueError is thrown. 
        If parsing succeeded, the function returns a ConfigParser object 
        with config parameters
        """

        config = ConfigParser(os.environ)
        # If config.ini does not exists original config object is not modified
        config.read("config.ini")

        config_params = {}
        try:
            config_params["port"] = int(os.getenv('SERVER_PORT', config["DEFAULT"]["SERVER_PORT"]))
            config_params["listen_backlog"] = int(os.getenv('SERVER_LISTEN_BACKLOG', config["DEFAULT"]["SERVER_LISTEN_BACKLOG"]))
            config_params["logging_level"] = os.getenv('LOGGING_LEVEL', config["DEFAULT"]["LOGGING_LEVEL"])
        except KeyError as e:
            raise KeyError("Key was not found. Error: {} .Aborting server".format(e))
        except ValueError as e:
            raise ValueError("Key could not be parsed. Error: {}. Aborting server".format(e))

        return config_params

    def _initialize_log(self, logging_level):
        """
        Python custom logging initialization

        Current timestamp is added to be able to identify in docker
        compose logs the date when the log has arrived
        """
        logging.basicConfig(
            format='%(asctime)s %(levelname)-8s %(message)s',
            level=logging_level,
            datefmt='%Y-%m-%d %H:%M:%S',
        )

    def _sigterm_handler(self, signum, frame):
        """
        SIGTERM signal handler

        Function that handles the SIGTERM signal to gracefully 
        shutdown the server and the current client. 
        """

        logging.info('action: signal_handler | result: in_progress | signal: SIGTERM')
        self._server_socket.close()
        self._server_protocol.close_client_socket()
        self._running = False
        logging.info('action: signal_handler | result: success | signal: SIGTERM')
