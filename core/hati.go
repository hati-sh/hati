package core

type Hati struct {
	config *Config
}

func NewHati(config *Config) Hati {
	return Hati{
		config: config,
	}
}

func (h *Hati) Start() error {
	serverTcp, err := NewServerTcp(h.config.ServerTcp.Host, h.config.ServerTcp.Port)
	if err != nil {
		return err
	}

	if err := serverTcp.Start(); err != nil {
		return err
	}

	// fmt.Println(MESSAGE_HEADER)

	return nil
}
