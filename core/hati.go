package core

import (
	"context"
	"errors"

	"github.com/hati-sh/hati/storage"
)

type Hati struct {
	config        *Config
	storage       storage.Storage
	serverTcp     ServerTcp
	stopCtx       context.Context
	stopCtxCancel context.CancelFunc
}

func NewHati(ctx context.Context, config *Config) Hati {
	stopCtx, stopCtxCancel := context.WithCancel(ctx)

	return Hati{
		config:        config,
		storage:       storage.New(),
		stopCtx:       stopCtx,
		stopCtxCancel: stopCtxCancel,
	}
}

func (h *Hati) Start() error {
	var err error

	h.serverTcp, err = NewServerTcp(h.stopCtx, h.commandHandler, h.config.ServerTcp.Host, h.config.ServerTcp.Port, h.config.ServerTcp.TlsEnabled)
	if err != nil {
		return err
	}

	if err := h.serverTcp.Start(); err != nil {
		return err
	}

	return nil
}

func (h *Hati) commandHandler(payload []byte) ([]byte, error) {

	if payload != nil {
		response := []byte("+OK\n")

		return response, nil
	}

	return nil, errors.New("+ERR\n")
}

func (h *Hati) Stop() {
	h.serverTcp.Stop()
}
