package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/communication"
	"github.com/op/go-logging"
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
	v.BindEnv("log", "level")

	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from config file. Using env variables instead")
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
	log.Infof("action: config | result: success | client_id: %s | server_address: %s | log_level: %s",
		v.GetString("id"),
		v.GetString("server.address"),
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
	id, err := strconv.ParseUint(config.ID, 10, 8)
	if err != nil {
		log.Criticalf("Invalid client ID: %s", config.ID)
	}

	agencyName := "Agencia de Quiniela " + config.ID
	agenciaDeQuiniela := NewAgenciaDeQuiniela(agencyName, uint8(id))

	return &Client{config, *agenciaDeQuiniela}
}

// sendStoreBetMsg extracts the bet information from the given line and
// sends it to agenciaDeQuiniela to be stored. Then it receives the response
// from the server and logs the result of the bet storage operation.
func (c *Client) sendStoreBetMsg(clientProtocol communication.ClientProtocol, line []string) error {
	if len(line) != 5 {
		return fmt.Errorf("Invalid number of fields in line: %v", line)
	}
	firstName, lastName, document, birthdate, number := line[0], line[1], line[2], line[3], line[4]
	log.Infof("action: crear_apuesta | result: success | dni: %v | number: %v", document, number)
	err := c.agenciaDeQuiniela.StoreBet(clientProtocol, firstName, lastName, document, birthdate, number)
	if err != nil {
		return err
	}
	log.Infof("action: registrar_apuesta | result: success")

	betStoreResponse, err := c.agenciaDeQuiniela.RecvBetStoreResponse(clientProtocol)
	if err != nil {
		return err
	}
	log.Infof("action: apuesta_enviada | result: %s | dni: %v | number: %v",
		func() string {
			if betStoreResponse.Success {
				return "success"
			}
			return "fail"
		}(),
		betStoreResponse.Document,
		betStoreResponse.Number,
	)
	return nil
}

// sendStoreBetMsgs reads the bets from the ./data/data.csv file and sends them to the
// server using the given ClientProtocol.
func (c *Client) sendStoreBetMsgs(clientProtocol communication.ClientProtocol) error {
	dataFile, err := os.Open("./data/data.csv")
	if err != nil {
		return err
	}

	reader := csv.NewReader(dataFile)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = c.sendStoreBetMsg(clientProtocol, line)
		if err != nil {
			return err
		}
	}

	defer dataFile.Close()
	return nil
}

// Start is the main method of the client. It is responsible for starting the client and
// handling the SIGTERM signal to gracefully shutdown the client. It creates a new
// ClientProtocol to communicate with the server, sends the bet and receives the
// response from the server. Finally, it logs the result of the bet storage operation.
func (c *Client) Start() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM)
	clientProtocol, err := communication.NewClientProtocol(c.config.ServerAddress, c.config.ID)
	if err != nil {
		log.Criticalf("%s", err)
		clientProtocol.Close()
		return
	}

	select {
	case <-signalChannel:
		log.Infof("action: sigterm_received | result: success | client_id: %v", c.config.ID)
		clientProtocol.Close()
		return
	default:
	}

	err = c.sendStoreBetMsgs(clientProtocol)
	if err != nil {
		log.Criticalf("%s", err)
	}

	clientProtocol.Close()
}
