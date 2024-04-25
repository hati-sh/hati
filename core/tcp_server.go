package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/hati-sh/hati/common"
	"github.com/hati-sh/hati/common/logger"
)

type PayloadHandler func(payload []byte) ([]byte, error)

type TcpServerConfig struct {
	Host       string
	Port       string
	TlsEnabled bool
}

type TcpServer struct {
	ctx              context.Context
	config           *TcpServerConfig
	tlsCertificate   tls.Certificate
	listener         net.Listener
	payloadHandler   PayloadHandler
	stopChan         chan bool
	stopWg           sync.WaitGroup
	listenerStopChan chan bool
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
		ctx:              ctx,
		config:           config,
		tlsCertificate:   cert,
		payloadHandler:   payloadHandler,
		stopChan:         make(chan bool),
		listenerStopChan: make(chan bool),
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
	fmt.Println("tcpServer Stop")
}

func (s *TcpServer) WaitForStop() {
	s.stopWg.Wait()
	fmt.Println("tcpServer WaitForStop")
}

func (s *TcpServer) startListener() {
	s.stopWg.Add(1)
	defer s.stopWg.Done()

OuterLoop:
	for {
		_, err := s.listener.Accept()

		if err != nil {
			select {
			case <-s.listenerStopChan:
				break OuterLoop
			default:
				logger.Error(err.Error())
			}
		} else {
			fmt.Println("accept conn?")
		}
	}

	fmt.Println("stop startListener")
}

func (s *TcpServer) daemon(stopChan <-chan bool) {
	s.stopWg.Add(1)

	defer close(s.stopChan)
	defer s.stopWg.Done()

OuterLoop:
	for {
		select {
		case <-stopChan:
			s.listener.Close()
			fmt.Println("daemon stop")
			s.listenerStopChan <- true

			break OuterLoop
		case <-s.ctx.Done():
			s.listener.Close()
			s.listenerStopChan <- true

			fmt.Println("ctx daemon stop")

			break OuterLoop
		}
	}
}
