package client

import (
	"crypto/tls"

	log "github.com/sirupsen/logrus"
)

func Dial(address string) (*tls.Conn, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	log.Debug("connecting to", address)
	conn, err := tls.Dial("tcp", address, conf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
