package core

import (
	"fmt"

	"github.com/hati-sh/hati/storage"
)

type Hati struct {
	config    *Config
	storage   storage.Storage
	serverTcp ServerTcp
}

func NewHati(config *Config) Hati {
	return Hati{
		config:  config,
		storage: storage.New(),
	}
}

func (h *Hati) Start() error {
	var err error

	h.serverTcp, err = NewServerTcp(h.config.ServerTcp.Host, h.config.ServerTcp.Port, h.config.ServerTcp.TlsEnabled)
	if err != nil {
		return err
	}

	if err := h.serverTcp.Start(h.processCommand); err != nil {
		return err
	}

	return nil
}

func (h *Hati) processCommand(payload []byte) {
	fmt.Println("processCommand")
}

func (h *Hati) Stop() {}
