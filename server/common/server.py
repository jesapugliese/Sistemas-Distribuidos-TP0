from configparser import ConfigParser
import logging
import os
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

        self._port = port
        self._listen_backlog = listen_backlog
        self._central_de_loteria = CentralDeLoteriaNacional()

    def start(self):
        # server_protocol = ServerProtocol(self._port, self._listen_backlog)
        # self._central_de_loteria.recibir_apuesta(server_protocol)
        # self._central_de_loteria.enviar_resultado(server_protocol)
        pass

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
