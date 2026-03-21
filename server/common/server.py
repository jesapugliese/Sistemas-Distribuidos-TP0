import logging
import os
import signal
import socket

from configparser import ConfigParser
from central_de_loteria import CentralDeLoteriaNacional
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

        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        self._server_protocol = ServerProtocol()
        self._central_de_loteria = CentralDeLoteriaNacional()

        signal.signal(signal.SIGTERM, self._sigterm_handler)

    def start(self):
        """
        Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while True:
            client_socket = self._accept_new_connection()
            if not client_socket:
                continue
            self._server_protocol.update_client_socket(client_socket)
            self._handle_client_connection()

    def _handle_client_connection(self):
        """
        Receives message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """

        try:
            logging.info('action: recibir_mensaje | result: in_progress')
            document, number = self._central_de_loteria.recv_msg_store_bets(self._server_protocol)
            logging.info('action: apuesta_almacenada | result: success | '
                         f'dni: {document} | numero: {number}')
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            logging.info('action: close_client_connection | result: in_progress')
            self._server_protocol.close_client_connection()
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
            logging.warning('action: accept_connections | result: fail')
            return None
        except OSError:
            logging.error('action: accept_connections | result: fail')
            return None
        
        return client_socket

    def _initialize_config(self):
        config = ConfigParser()
        config.read(os.path.join(os.path.dirname(__file__), "config.ini"))
        return {
            "logging_level": getattr(logging, config["DEFAULT"]["LoggingLevel"]
                                     .upper()),
            "port": int(config["DEFAULT"]["Port"]),
            "listen_backlog": int(config["DEFAULT"]["ListenBacklog"])
        }

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
        self._server_protocol.shutdown()
        logging.info('action: signal_handler | result: success | signal: SIGTERM')

        exit(0)
