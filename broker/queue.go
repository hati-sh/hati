package broker

import "context"

type Queue struct {
	ctx    context.Context
	name   string
	config QueueConfig
}

func NewQueue(ctx context.Context, config QueueConfig) *Queue {
	return &Queue{
		ctx:    ctx,
		config: config,
	}
}
