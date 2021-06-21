package main

import (
	"github.com/mohibeyki/go-proxy/pkg/config"
	"github.com/mohibeyki/go-proxy/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	log.Infoln("go-proxy is starting")
	config.Init()
	server.StartTunnel(viper.GetString("upstream.address"))
}
