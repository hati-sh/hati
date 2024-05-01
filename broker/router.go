package broker

import "context"

type Router struct {
	ctx    context.Context
	config RouterConfig
}

func NewRouter(ctx context.Context, config RouterConfig) *Router {
	return &Router{
		ctx:    ctx,
		config: config,
	}
}
