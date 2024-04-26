package core

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"net"
	"sync"

	"github.com/hati-sh/hati/common/logger"
)

type TcpServerClient struct {
	ctx                        context.Context
	ctxStop                    context.Context
	ctxStopCancel              context.CancelFunc
	conn                       net.Conn
	stopWg                     *sync.WaitGroup
	connReadOutChan            chan []byte
	payloads                   chan []byte
	stopProcessingPayloadsChan chan bool
	payloadHandler             PayloadHandler
	clientStoppedChan          chan net.Conn
}

func NewTcpServerClient(ctx context.Context, stopWg *sync.WaitGroup, clientStoppedChan chan net.Conn, conn net.Conn, payloadHandler PayloadHandler) *TcpServerClient {
	ctxStop, ctxStopCancel := context.WithCancel(ctx)

	return &TcpServerClient{
		ctx:                        ctx,
		ctxStop:                    ctxStop,
		ctxStopCancel:              ctxStopCancel,
		conn:                       conn,
		stopWg:                     stopWg,
		payloads:                   make(chan []byte, TCP_PAYLOAD_HANDLER_CHAN_SIZE),
		stopProcessingPayloadsChan: make(chan bool),
		connReadOutChan:            make(chan []byte, TCP_PAYLOAD_HANDLER_CHAN_SIZE),
		payloadHandler:             payloadHandler,
		clientStoppedChan:          clientStoppedChan,
	}
}

func (c *TcpServerClient) start() {
	c.stopWg.Add(1)
	defer c.conn.Close()
	defer c.stopWg.Done()

	go c.processPayloads()
	go c.scanForIncomingBytes()

OuterLoop:
	for {
		select {
		case payload := <-c.connReadOutChan:
			{
				c.payloads <- payload
			}
		case <-c.ctx.Done():
			{
				break OuterLoop
			}
		case <-c.ctxStop.Done():
			{
				break OuterLoop
			}
		}
	}

	if c.stopProcessingPayloadsChan != nil {
		close(c.stopProcessingPayloadsChan)
	}

	c.clientStoppedChan <- c.conn

	fmt.Println("stop TcpServerClient start")
}

func (c *TcpServerClient) scanForIncomingBytes() {
	c.stopWg.Add(1)

	defer c.conn.Close()
	defer close(c.connReadOutChan)
	defer c.ctxStopCancel()
	defer c.stopWg.Done()

	scanner := bufio.NewScanner(c.conn)
	buf := make([]byte, 0, 1<<20)
	scanner.Buffer(buf, math.MaxInt)

	for {
		ok := scanner.Scan()

		if !ok {
			if err := scanner.Err(); err != nil {
				logger.Error(err.Error())
				break
			}
			break
		}

		c.connReadOutChan <- scanner.Bytes()
	}

	logger.Debug("stop scanForIncomingBytes")
}

func (c *TcpServerClient) processPayloads() {
	c.stopWg.Add(1)
	defer c.conn.Close()

	defer c.ctxStopCancel()
	defer c.stopWg.Done()

OuterLoop:
	for {
		select {
		case payload := <-c.payloads:
			{
				response, err := c.payloadHandler(payload)
				if err != nil {
					if _, err := c.conn.Write([]byte(err.Error())); err != nil {
						logger.Error(err)

						break OuterLoop
					}
					continue
				}

				_, err = c.conn.Write(response)
				if err != nil {
					logger.Error(err)

					break OuterLoop
				}
			}
		case <-c.stopProcessingPayloadsChan:
			logger.Debug("stopChan: stopProcessPayloadsChan")

			break OuterLoop
		}
	}

	logger.Debug("stop processPayloads")
}
