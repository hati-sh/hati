package broker

import "context"

type Broker struct {
	ctx    context.Context
	router map[string]Broker
	queue  map[string]Queue
}

func New(ctx context.Context) Broker {
	return Broker{
		ctx:    ctx,
		router: make(map[string]Broker),
		queue:  make(map[string]Queue),
	}
}
