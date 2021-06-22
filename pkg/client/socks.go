package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IPv4 struct {
	A, B, C, D byte
}

type IPv6 struct {
	A, B, C, D, E, F, G, H uint16
}

func (ip IPv4) toAddr() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip.A, ip.B, ip.C, ip.D)
}

func (ip IPv6) toAddr() string {
	return fmt.Sprintf("[%x:%x:%x:%x:%x:%x:%x:%x]", ip.A, ip.B, ip.C, ip.D, ip.E, ip.F, ip.G, ip.H)
}

func StartSocksServer(addr string) {
	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	defer server.Close()
	log.Info("accepting connections on", addr)
	for {
		client, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	defer client.Close()

	buffer := make([]byte, 1024)
	_, err := client.Read(buffer[:])
	if err != nil {
		log.Error(err)
		return
	}
	// sending back 0x05, 0x00 -> no authentication is supported
	client.Write([]byte{0x05, 0x00})

	n, err := client.Read(buffer)
	if err != nil {
		log.Error(err)
		return
	}

	var addr string
	var port uint16
	err = binary.Read(bytes.NewReader(buffer[n-2:n]), binary.BigEndian, &port)
	if err != nil {
		log.Error(err)
		return
	}

	// extracting address from buffer
	switch buffer[3] {
	// IPv4:port
	case 0x01:
		ip := IPv4{}
		if err := binary.Read(bytes.NewReader(buffer[4:n-2]), binary.BigEndian, &ip); err != nil {
			log.Error("Request parse error")
			return
		}
		addr = ip.toAddr()
	// n, domain:port
	case 0x03:
		addr = string(buffer[5 : n-2])
	// IPv6:port
	case 0x04:
		ip := IPv6{}
		if err := binary.Read(bytes.NewReader(buffer[4:n-2]), binary.BigEndian, &ip); err != nil {
			log.Error("Request parse error")
			return
		}
		addr = ip.toAddr()
	}
	addr = fmt.Sprintf("%s:%d", addr, port)

	upstream, err := Dial(viper.GetString("upstream.address"))
	if err != nil {
		log.Error(err)
		return
	}
	defer upstream.Close()

	// sending address as the first thing so that the upstream can establish connection
	upstream.Write([]byte(addr + "\n"))
	// telling client that we have established connection to the remote
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	go io.Copy(upstream, client)
	io.Copy(client, upstream)
}
