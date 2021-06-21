package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	log.Infoln("reading config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %w", err)
	}
}
