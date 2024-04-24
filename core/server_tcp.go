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
	conn     net.Conn
	payloads chan []byte
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

func (s ServerTcp) Start() error {
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

		go s.startListener(s.listener)
	} else {
		go s.startListener(s.listener)
	}

	return nil
}

func (s ServerTcp) startListener(listener net.Listener) {
	defer s.listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		client := &Client{
			conn:     conn,
			payloads: make(chan []byte, 5000000),
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
		// scanner.
		// fmt.Println(scanner.Text())
	}
	// bufReader := bufio.NewReader(client.conn)
	// // bufio.Scanner()
	// for {
	// 	receivedBytes, err := bufReader.ReadBytes('\n')

	// 	if len(receivedBytes) < 1 {
	// 		client.conn.Close()
	// 		return
	// 	}

	// 	if err != nil {
	// 		if err != io.EOF {
	// 			fmt.Println("read error:", err)

	// 			client.conn.Close()
	// 			return
	// 		}
	// 		fmt.Println(err)
	// 		break
	// 	}

	// 	// payloadBuffer := make([]byte, contentLength)
	// 	// n, _ := io.ReadFull(client.conn, payloadBuffer)

	// 	// fmt.Println(string(contentLengthBytes))
	// 	// fmt.Println("=======")
	// 	// fmt.Println(n)
	// 	// fmt.Println(payloadBuffer)
	// 	// fmt.Println(contentLength)

	// 	// fmt.Println(string(receivedContentLengthBytes[0 : len(receivedContentLengthBytes)-1]))

	// 	// payloadBuffer := make([]byte,)

	// 	// cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))
	// 	// args := bytes.TrimSpace(bytes.TrimPrefix(message, cmd))

	// 	// buffer := make([]byte, 0, 10)
	// 	// n, _ := io.ReadFull(client.conn, buffer)

	// 	// fmt.Println(n)
	// 	// fmt.Println(buffer)
	// 	// receivedMore , _ := bufio.NewReader(client.conn).Read
	// 	// buf = append(buf, tmp[:n]...)

	// 	// receivedMessage, errParse := ParseBytesToMessage(receivedBytes)
	// 	// if errParse != nil {
	// 	// 	client.conn.Write([]byte(errParse.Error()))
	// 	// 	client.conn.Close()
	// 	// 	return
	// 	// }
	// 	client.payloads <- receivedBytes

	// 	// receivedMessage, err := ParseBytesToMessage(buf)
	// }
}

func (client *Client) processPayloads() {

	for {
		<-client.payloads

		_, err := client.conn.Write([]byte("+OK\n"))
		if err != nil {
			fmt.Println(err)
			client.conn.Close()
		}

		lock.Lock()
		globalCounter++
		lock.Unlock()

		fmt.Println(globalCounter)
	}

}
