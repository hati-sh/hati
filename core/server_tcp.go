package core

import (
	"crypto/tls"
	"fmt"
	"io"
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
			conn: conn,
		}

		go client.handleRequest()
	}
}

func (client *Client) handleRequest() {
	const MAX_BUFFER_SIZE = 1024 * 8
	const TMP_BUFFER_SIZE = 1024 * 1

	buf := make([]byte, 0, MAX_BUFFER_SIZE)
	tmp := make([]byte, TMP_BUFFER_SIZE)

	for {
		n, err := client.conn.Read(tmp)

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)

				client.conn.Close()
				return
			}
			fmt.Println(err)
			break
		}
		buf = append(buf, tmp[:n]...)

		receivedMessage, errParse := ParseBytesToMessage(buf)
		if errParse != nil {
			client.conn.Write([]byte(errParse.Error()))
			client.conn.Close()
			return
		}

		client.conn.Write([]byte("+OK\n"))
		fmt.Println(string(receivedMessage.payload))
		// receivedMessage, err := ParseBytesToMessage(buf)
	}
}
