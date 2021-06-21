package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

type sockIP struct {
	A, B, C, D byte
	PORT       uint16
}

func (ip sockIP) toAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", ip.A, ip.B, ip.C, ip.D, ip.PORT)
}

func StartServer(addr string) {
	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	defer server.Close()
	log.Infoln("accepting connections on", addr)
	for {
		client, err := server.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handleClient(client)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buffer [1024]byte

	_, err := conn.Read(buffer[:])
	if err != nil {
		log.Errorln(err)
		return
	}
	conn.Write([]byte{0x05, 0x00})

	n, err := conn.Read(buffer[:])
	if err != nil {
		log.Errorln(err)
		return
	}

	var addr string
	switch buffer[3] {
	case 0x01:
		sip := sockIP{}
		if err := binary.Read(bytes.NewReader(buffer[4:n]), binary.BigEndian, &sip); err != nil {
			log.Errorln("Request parse error")
			return
		}
		addr = sip.toAddr()
	case 0x03:
		host := string(buffer[5 : n-2])
		var port uint16
		err = binary.Read(bytes.NewReader(buffer[n-2:n]), binary.BigEndian, &port)
		if err != nil {
			log.Errorln(err)
			return
		}
		addr = fmt.Sprintf("%s:%d", host, port)
	}

	server, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorln(err)
		return
	}
	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	go io.Copy(server, conn)
	io.Copy(conn, server)
}
