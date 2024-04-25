package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/hati-sh/hati/common"
)

const TCP_PAYLOAD_HANDLER_CHAN_SIZE = 2000000

// const TCP_READ_BUFFER_RATE_LIMIT_DELAY = 1 * time.Microsecond

type PayloadHandlerCallback func(payload []byte) ([]byte, error)

type ServerTcp struct {
	host                   string
	port                   string
	tlsEnabled             bool
	tlsCertificate         tls.Certificate
	listener               net.Listener
	payloadHandlerCallback PayloadHandlerCallback
	stopChan               chan bool
}

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

func NewServerTcp(ctx context.Context, payloadHandlerCallback PayloadHandlerCallback, host string, port string, tlsEnabled bool) (ServerTcp, error) {
	var cert tls.Certificate
	var err error

	if tlsEnabled {
		cert, err = common.GenX509KeyPair()
		if err != nil {
			return ServerTcp{}, err
		}
	}

	return ServerTcp{
		host:                   host,
		port:                   port,
		tlsEnabled:             tlsEnabled,
		tlsCertificate:         cert,
		payloadHandlerCallback: payloadHandlerCallback,
		stopChan:               make(chan bool),
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

func (s ServerTcp) Stop() {
	// s.stopChan <- true
}

func (s ServerTcp) startListener(listener net.Listener) {
	defer s.listener.Close()

	for {
		// select {
		// case <-s.ctx.Done():
		// 	{

		// 	}
		// default:
		// 	{
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		client := NewServerTcpClient(conn, s.payloadHandlerCallback)

		go client.handleRequest()
		// }
		// }
	}
}
