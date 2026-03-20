package common

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/communication"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("log")

type ClientConfig struct {
	ID            string
	ServerAddress string
}

type Client struct {
	config            ClientConfig
	agenciaDeQuiniela AgenciaDeQuiniela
}

// InitConfig Function that uses viper library to parse configuration parameters.
func InitConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix("cli")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.BindEnv("id")
	v.BindEnv("server", "address")
	v.BindEnv("loop", "period")
	v.BindEnv("loop", "amount")
	v.BindEnv("log", "level")

	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from config file. Using env variables instead")
	}

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_PERIOD env var as time.Duration.")
	}

	return v, nil
}

// InitLogger Receives the log level to be set in go-logging as a string. This method
// parses the string and set the level to the logger. If the level string is not
// valid an error is returned
func InitLogger(logLevel string) error {
	baseBackend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05} %{level:.5s}     %{message}`,
	)
	backendFormatter := logging.NewBackendFormatter(baseBackend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	logLevelCode, err := logging.LogLevel(logLevel)
	if err != nil {
		return err
	}
	backendLeveled.SetLevel(logLevelCode, "")

	logging.SetBackend(backendLeveled)
	return nil
}

// PrintConfig Print all the configuration parameters of the program.
// For debugging purposes only
func PrintConfig(v *viper.Viper) {
	log.Infof("action: config | result: success | client_id: %s | server_address: %s | loop_amount: %v | loop_period: %v | log_level: %s",
		v.GetString("id"),
		v.GetString("server.address"),
		v.GetInt("loop.amount"),
		v.GetDuration("loop.period"),
		v.GetString("log.level"),
	)
}

// NewClient Initializes a new client. This method is responsible for
// initializing the configuration and the logger.
func NewClient() *Client {
	v, err := InitConfig()
	if err != nil {
		log.Criticalf("%s", err)
	}

	if err := InitLogger(v.GetString("log.level")); err != nil {
		log.Criticalf("%s", err)
	}

	PrintConfig(v)

	config := ClientConfig{
		ServerAddress: v.GetString("server.address"),
		ID:            v.GetString("id"),
	}

	nombreAgencia := "Agencia de Quiniela " + config.ID
	agenciaDeQuiniela := NewAgenciaDeQuiniela(nombreAgencia, config.ID)

	return &Client{config, *agenciaDeQuiniela}
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and returned to the caller.
func (c *Client) createClientSocket() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return nil, err
	}
	return conn, nil
}

func (c *Client) Start() {
	clientSocket, err := c.createClientSocket()
	if err != nil {
		log.Criticalf("%s", err)
	}

	clientProtocol := communication.NewClientProtocol(clientSocket)

	err = c.agenciaDeQuiniela.RegistrarApuesta(clientProtocol)
	if err != nil {
		log.Criticalf("%s", err)
	}
	_, err = c.agenciaDeQuiniela.RecibirResultado(clientProtocol)
	if err != nil {
		log.Criticalf("%s", err)
	}

	// log.Infof("%s", result)
}
