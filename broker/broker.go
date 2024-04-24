package broker

type Broker struct {
	router map[string]Broker
	queue  map[string]Queue
}

func NewBroker() Broker {
	return Broker{
		router: make(map[string]Broker),
		queue:  make(map[string]Queue),
	}
}
