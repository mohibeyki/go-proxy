package server

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

func StartTunnel(address string) {
	certificates, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("starting tcp server on", address)
	config := &tls.Config{Certificates: []tls.Certificate{certificates}}
	ln, err := tls.Listen("tcp", address, config)
	if err != nil {
		log.Error(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(client net.Conn) {
	defer client.Close()
	reader := bufio.NewReader(client)

	addr, err := reader.ReadString('\n')
	addr = strings.TrimSpace(addr)
	log.Debugf("address: %s", addr)
	if err != nil {
		log.Error(err)
		return
	}

	target, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error(err)
		return
	}

	go io.Copy(target, client)
	io.Copy(client, target)
}
