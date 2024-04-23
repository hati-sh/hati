package core

type Hati struct{}

func NewHati() Hati {
	return Hati{}
}

func (h *Hati) Start() error {
	serverTcp, err := NewServerTcp("0.0.0.0", "4242")
	if err != nil {
		return err
	}

	if err := serverTcp.Start(); err != nil {
		return err
	}

	// fmt.Println(MESSAGE_HEADER)

	return nil
}
