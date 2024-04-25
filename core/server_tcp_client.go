package core

import (
	"bufio"
	"net"
	"sync"

	"github.com/hati-sh/hati/common/logger"
)

type ServerTcpClient struct {
	conn                       net.Conn
	payloads                   chan []byte
	stopChan                   chan bool
	stopProcessingPayloadsChan chan bool
	stopWg                     sync.WaitGroup
	payloadHandlerCallback     PayloadHandlerCallback
	removeConnectionCallback   func(conn net.Conn)
	closedConnChan             chan<- bool
	connReadOutChan            chan []byte
}

func NewServerTcpClient(conn net.Conn, payloadHandlerCallback PayloadHandlerCallback, closedConnChan chan<- bool) *ServerTcpClient {
	return &ServerTcpClient{
		conn:                       conn,
		payloads:                   make(chan []byte, TCP_PAYLOAD_HANDLER_CHAN_SIZE),
		stopChan:                   make(chan bool),
		stopProcessingPayloadsChan: make(chan bool),
		payloadHandlerCallback:     payloadHandlerCallback,
		closedConnChan:             closedConnChan,
		connReadOutChan:            make(chan []byte, TCP_PAYLOAD_HANDLER_CHAN_SIZE),
	}
}

// func (client *ServerTcpClient) Stop() {
// 	client.stopWg.Wait()
// }

func (client *ServerTcpClient) scanForIncomingBytes() {
	defer client.conn.Close()
	defer close(client.connReadOutChan)
	defer client.stopWg.Done()

	scanner := bufio.NewScanner(client.conn)
	for {
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				logger.Error(err.Error())
				break
			}
			break
		}
		client.connReadOutChan <- scanner.Bytes()
	}
}

func (client *ServerTcpClient) handleRequest() {
	defer client.conn.Close()

	client.stopWg.Add(1)
	go client.processPayloads()

	// scanner := bufio.NewScanner(client.conn)
	go client.scanForIncomingBytes()
	client.stopWg.Add(1)
OuterLoop:
	for {
		select {
		case payload := <-client.connReadOutChan:
			{
				client.payloads <- payload
			}
		case <-client.stopChan:
			{
				break OuterLoop
			}
		}
	}

	if client.stopProcessingPayloadsChan != nil {
		client.stopProcessingPayloadsChan <- true
	}

	client.stopWg.Done()

	client.closedConnChan <- true

	logger.Debug("stop: handleRequest")
}

func (client *ServerTcpClient) processPayloads() {
OuterLoop:
	for {
		select {
		case payload := <-client.payloads:
			{
				response, err := client.payloadHandlerCallback(payload)
				if err != nil {
					if _, err := client.conn.Write([]byte(err.Error())); err != nil {
						logger.Error(err)

						break OuterLoop
					}
					continue
				}

				_, err = client.conn.Write(response)
				if err != nil {
					logger.Error(err)

					break OuterLoop
				}
			}
		case <-client.stopProcessingPayloadsChan:
			logger.Debug("stopChan: stopProcessPayloadsChan")

			break OuterLoop
		}
	}

	_ = client.conn.Close()
	close(client.stopProcessingPayloadsChan)

	client.stopWg.Done()
}
