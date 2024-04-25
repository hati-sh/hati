package core

import (
	"context"
	"fmt"
	"sync"
)

type PayloadHandler func(payload []byte) ([]byte, error)

type TcpServerConfig struct {
	Host       string
	Port       string
	TlsEnabled bool
}

type TcpServer struct {
	ctx            context.Context
	config         *TcpServerConfig
	payloadHandler PayloadHandler
	stopChan       chan bool
	stopWg         sync.WaitGroup
}

func NewTcpServer(ctx context.Context, config *TcpServerConfig, payloadHandler PayloadHandler) TcpServer {
	return TcpServer{
		ctx:            ctx,
		config:         config,
		payloadHandler: payloadHandler,
		stopChan:       make(chan bool),
	}
}

func (s *TcpServer) Start() error {
	s.stopWg.Add(1)
	go s.daemon(s.stopChan)

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

func (s *TcpServer) daemon(stopChan <-chan bool) {
	for {
		select {
		case <-stopChan:
			fmt.Println("daemon stop")

			s.stopWg.Done()
			return
		case <-s.ctx.Done():
			close(s.stopChan)
			fmt.Println("ctx daemon stop")

			s.stopWg.Done()
			return
		}
	}
}
