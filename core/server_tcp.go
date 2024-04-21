package core

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hati-sh/hati/common"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

type Client struct {
	conn net.Conn
}

type ServerTcp struct {
	tlsCertificate tls.Certificate
	host           string
	port           string
}

func NewServerTcp(host string, port string) (ServerTcp, error) {
	cert, err := common.GenX509KeyPair()
	if err != nil {
		return ServerTcp{}, err
	}

	return ServerTcp{
		tlsCertificate: cert,
		host:           host,
		port:           port,
	}, nil
}

func (s ServerTcp) Start() error {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}

	config := &tls.Config{
		Certificates: make([]tls.Certificate, 1),
	}
	config.Certificates[0] = s.tlsCertificate

	tlsListener := tls.NewListener(tcpKeepAliveListener{listener.(*net.TCPListener)}, config)

	defer listener.Close()
	defer tlsListener.Close()

	fmt.Printf("TCP listening at: %s:%s\n", s.host, s.port)
	fmt.Printf("TLS: ON\n")

	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequest()
	}

	return nil
}

func (client *Client) handleRequest() {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}
		fmt.Printf("Message incoming: %s", string(message))
		client.conn.Write([]byte("Message received.\n"))
	}
}
