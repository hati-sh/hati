package core

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type ServerTcpClient struct {
	conn                       net.Conn
	payloads                   chan []byte
	stopChan                   chan bool
	stopProcessingPayloadsChan chan bool
	stopWg                     sync.WaitGroup
	payloadHandlerCallback     PayloadHandlerCallback
}

func NewServerTcpClient(conn net.Conn, payloadHandlerCallback PayloadHandlerCallback) *ServerTcpClient {
	return &ServerTcpClient{
		conn:                       conn,
		payloads:                   make(chan []byte, TCP_PAYLOAD_HANDLER_CHAN_SIZE),
		stopChan:                   make(chan bool),
		stopProcessingPayloadsChan: make(chan bool),
		payloadHandlerCallback:     payloadHandlerCallback,
	}
}

func (client *ServerTcpClient) handleRequest() {
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

func (client *ServerTcpClient) processPayloads() {
Exit:
	for {
		select {
		case payload := <-client.payloads:
			{
				response, err := client.payloadHandlerCallback(payload)
				if err != nil {
					if _, err := client.conn.Write([]byte(err.Error())); err != nil {
						client.conn.Close()
						break
					}
					continue
				}

				_, err = client.conn.Write(response)
				if err != nil {
					fmt.Println(err)
					client.conn.Close()
					break
				}
			}
		case <-client.stopProcessingPayloadsChan:
			fmt.Println("stopChan: stopProcessPayloadsChan")
			break Exit
		}
	}
}
