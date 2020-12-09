package main

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// handler is used to handler responses from redis pubsub
type handler struct {
	ctx       context.Context
	log       *zap.Logger
	websocket *server
	redis     *redis.Client
}

func newHandler(ctx context.Context, logger *zap.Logger, server *server, redis *redis.Client) *handler {
	return &handler{
		ctx:       ctx,
		log:       logger,
		websocket: server,
		redis:     redis,
	}
}

// handleResponse handles a single response message from redis.
// The function simply forwards the message to all connected clients.
func (h *handler) handleResponse(channel string, msg string) {
	h.log.Info("response received", zap.String("channel", channel), zap.String("msg", string(msg)))

	// Simply forward the message to all connected clients
	h.websocket.broadcast(msg)

	h.log.Info("response broadcasted")
}

// run will start a redis subscriber and listen for incoming messages
func (h *handler) run(channel string) {
	pubsub := h.redis.Subscribe(h.ctx, channel)

	// Subscribe to a redis channel
	rps := pubsub.Channel()
	// Wait for messages
	for {
		select {
		case <-h.ctx.Done():
			h.log.Info("closing context")
			break
		case n := <-rps:
			h.handleResponse(n.Channel, n.Payload)
		}
	}

}
