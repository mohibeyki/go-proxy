package client

import (
	"crypto/tls"

	log "github.com/sirupsen/logrus"
)

func Dial(address string) (*tls.Conn, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	log.Infoln("connecting to", address)
	conn, err := tls.Dial("tcp", address, conf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// func SendData(conn *tls.Conn, data []byte) error {
// 	_, err := conn.Write(data)
// 	if err != nil {
// 		return err
// 	}

// 	buff := make([]byte, 1024)
// 	n, err := conn.Read(buff)
// 	if err != nil {
// 		return err
// 	}

// 	println(string(buff[:n]))
// 	return nil
// }
