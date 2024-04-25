package core

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hati-sh/hati/common"
)

const TcpHandlerPayloadChanSize = 5000000

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
	conn                   net.Conn
	payloads               chan []byte
	commandHandlerCallback CommandHandlerCallback
}

type ServerTcp struct {
	tlsCertificate tls.Certificate
	host           string
	port           string
	tlsEnabled     bool
	listener       net.Listener
}

func NewServerTcp(host string, port string, tlsEnabled bool) (ServerTcp, error) {
	var cert tls.Certificate
	var err error

	if tlsEnabled {
		cert, err = common.GenX509KeyPair()
		if err != nil {
			return ServerTcp{}, err
		}
	}

	return ServerTcp{
		tlsCertificate: cert,
		host:           host,
		port:           port,
		tlsEnabled:     tlsEnabled,
	}, nil
}

func (s ServerTcp) Start(commandHandlerCallback CommandHandlerCallback) error {
	var err error

	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}

	config := &tls.Config{}

	if s.tlsEnabled {
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0] = s.tlsCertificate

		s.listener = tls.NewListener(tcpKeepAliveListener{s.listener.(*net.TCPListener)}, config)

		go s.startListener(s.listener, commandHandlerCallback)
	} else {
		go s.startListener(s.listener, commandHandlerCallback)
	}

	return nil
}

func (s ServerTcp) startListener(listener net.Listener, commandHandlerCallback CommandHandlerCallback) {
	defer s.listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		client := &Client{
			conn:                   conn,
			payloads:               make(chan []byte, TcpHandlerPayloadChanSize),
			commandHandlerCallback: commandHandlerCallback,
		}

		go client.handleRequest()
	}
}

var globalCounter = 0
var lock sync.Mutex

func (client *Client) handleRequest() {
	go client.processPayloads()

	scanner := bufio.NewScanner(client.conn)
	for {
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				fmt.Println(err)
				break
			}

			break
		}
		client.payloads <- scanner.Bytes()
	}
}

func (client *Client) processPayloads() {

	for {
		payload := <-client.payloads

		lock.Lock()
		globalCounter++
		lock.Unlock()

		_, err := client.conn.Write([]byte("+OK\n"))
		if err != nil {
			fmt.Println(err)
			client.conn.Close()
		}

		fmt.Println(string(payload))

		client.commandHandlerCallback(payload)
	}

}
