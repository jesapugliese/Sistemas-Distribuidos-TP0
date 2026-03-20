package main

import (
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
)

func main() {
	client := common.NewClient()
	client.Start()
}
