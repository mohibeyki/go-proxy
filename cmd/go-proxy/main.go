package main

import (
	"github.com/mohibeyki/go-proxy/cmd/go-proxy/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("go-proxy is starting")
	server.StartServer("127.0.0.1:8080")
}
