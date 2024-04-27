package core

import (
	"context"
	"sync"

	"github.com/hati-sh/hati/broker"
	"github.com/hati-sh/hati/storage"
)

type Hati struct {
	config         *Config
	stopCtx        context.Context
	stopCtxCancel  context.CancelFunc
	stopWg         sync.WaitGroup
	tcpServer      *TcpServer
	rpcServer      *RpcServer
	commandHandler CommandHandler
	storage        storage.Storage
	broker         broker.Broker
}

func NewHati(ctx context.Context, config *Config) *Hati {
	stopCtx, stopCtxCancel := context.WithCancel(ctx)

	hati := &Hati{
		config:        config,
		stopCtx:       stopCtx,
		stopCtxCancel: stopCtxCancel,
	}

	hati.storage = storage.New(hati.stopCtx)
	hati.broker = broker.New(hati.stopCtx)

	hati.commandHandler = CommandHandler{
		ctx:     hati.stopCtx,
		storage: &hati.storage,
		broker:  &hati.broker,
	}

	hati.tcpServer = NewTcpServer(hati.stopCtx, config.ServerTcp, hati.commandHandler.processPayload)
	hati.rpcServer = NewRpcServer(hati.stopCtx, &hati.storage, config.ServerRpc)

	return hati
}

func (h *Hati) Start() error {
	var err error

	if err = h.tcpServer.Start(); err != nil {
		return err
	}

	if h.config.ServerRpc.Enabled {
		if err = h.rpcServer.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (h *Hati) Stop() {
	h.stopCtxCancel()

	h.tcpServer.WaitForStop()

	h.stopWg.Wait()
}
