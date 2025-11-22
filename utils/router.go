package utils

import (
	"context"
	"fmt"
)

type Consumer interface {
	Consume(ctx context.Context, data string) error
}

type Router struct {
	routes map[string]Consumer
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Consumer),
	}
}

func (r *Router) Register(queue string, c Consumer) {
	r.routes[queue] = c
}

func (r *Router) GetQueues() []string {
	queues := make([]string, 0, len(r.routes))
	for q := range r.routes {
		queues = append(queues, q)
	}
	return queues
}

func (r *Router) Route(ctx context.Context, queue string, data string) error {
	c, ok := r.routes[queue]
	if !ok {
		return fmt.Errorf("no consumer found for queue: %s", queue)
	}
	return c.Consume(ctx, data)
}
