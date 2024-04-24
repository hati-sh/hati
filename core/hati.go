package core

type Hati struct {
	config    *Config
	serverTcp ServerTcp
}

func NewHati(config *Config) Hati {
	return Hati{
		config: config,
	}
}

func (h *Hati) Start() error {
	var err error

	h.serverTcp, err = NewServerTcp(h.config.ServerTcp.Host, h.config.ServerTcp.Port, h.config.ServerTcp.TlsEnabled)
	if err != nil {
		return err
	}

	if err := h.serverTcp.Start(); err != nil {
		return err
	}

	return nil
}

func (h *Hati) Stop() {}
