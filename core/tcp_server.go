package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hati-sh/hati/common"
	"github.com/hati-sh/hati/common/logger"
)

const TCP_PAYLOAD_HANDLER_CHAN_SIZE = 2000000

type PayloadHandler func(payload []byte) ([]byte, error)

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

type TcpServerConfig struct {
	Host       string
	Port       string
	TlsEnabled bool
}

type TcpServer struct {
	ctx               context.Context
	config            *TcpServerConfig
	tlsCertificate    tls.Certificate
	listener          net.Listener
	payloadHandler    PayloadHandler
	stopChan          chan bool
	stopWg            sync.WaitGroup
	listenerStopChan  chan bool
	clients           map[net.Conn]*TcpServerClient
	clientsMutex      sync.Mutex
	clientStoppedChan chan net.Conn
}

func NewTcpServer(ctx context.Context, config *TcpServerConfig, payloadHandler PayloadHandler) TcpServer {
	var cert tls.Certificate
	var err error

	if config.TlsEnabled {
		cert, err = common.GenX509KeyPair()

		if err != nil {
			panic(err)
		}
	}

	return TcpServer{
		ctx:               ctx,
		config:            config,
		tlsCertificate:    cert,
		payloadHandler:    payloadHandler,
		stopChan:          make(chan bool),
		listenerStopChan:  make(chan bool),
		clients:           make(map[net.Conn]*TcpServerClient),
		clientStoppedChan: make(chan net.Conn, 1000),
	}
}

func (s *TcpServer) Start() error {
	go s.daemon(s.stopChan)

	var err error

	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))
	if err != nil {
		return err
	}

	if s.config.TlsEnabled {
		config := &tls.Config{}
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0] = s.tlsCertificate

		s.listener = tls.NewListener(tcpKeepAliveListener{s.listener.(*net.TCPListener)}, config)
	}

	go s.startListener()

	return nil
}

func (s *TcpServer) Stop() {
	s.stopChan <- true

	s.stopWg.Wait()
	logger.Debug("tcpServer Stop")
}

func (s *TcpServer) WaitForStop() {
	s.stopWg.Wait()
	logger.Debug("tcpServer WaitForStop")
}

func (s *TcpServer) startListener() {
	s.stopWg.Add(1)
	defer s.stopWg.Done()

OuterLoop:
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			select {
			case <-s.listenerStopChan:
				break OuterLoop
			default:
				logger.Error(err.Error())
			}
		} else {
			s.handleConnection(conn)
		}
	}

	logger.Debug("stop startListener")
}

func (s *TcpServer) handleConnection(conn net.Conn) {
	s.clientsMutex.Lock()
	s.clients[conn] = NewTcpServerClient(s.ctx, &s.stopWg, s.clientStoppedChan, conn, s.payloadHandler)
	s.clientsMutex.Unlock()

	go s.clients[conn].start()
}

func (s *TcpServer) daemon(stopChan <-chan bool) {
	s.stopWg.Add(1)

	defer close(s.stopChan)
	defer s.stopWg.Done()

OuterLoop:
	for {
		select {
		case conn := <-s.clientStoppedChan:
			conn.Close()
			s.clientsMutex.Lock()
			delete(s.clients, conn)
			s.clientsMutex.Unlock()
			logger.Debug("tcp client removed from map")
		case <-stopChan:
			s.listener.Close()

			logger.Debug("daemon stop")

			s.listenerStopChan <- true

			break OuterLoop
		case <-s.ctx.Done():
			s.listener.Close()
			s.listenerStopChan <- true

			logger.Debug("ctx daemon stop")

			break OuterLoop
		}
	}
}
