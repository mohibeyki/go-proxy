package main

import (
	"github.com/mohibeyki/go-proxy/pkg/client"
	"github.com/mohibeyki/go-proxy/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	log.Infoln("starting client")
	config.Init()
	client.StartSocksServer(viper.GetString("socks.address"))
}
