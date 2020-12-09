package main

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// server is a websocket server handler
type server struct {
	upgrader websocket.Upgrader
	clients  *clients
	redis    *redis.Client
	channel  string
	ctx      context.Context
	log      *zap.Logger
}

func newServer(ctx context.Context, logger *zap.Logger, redis *redis.Client, channel string) *server {
	return &server{
		upgrader: websocket.Upgrader{
			// Allow all origins.
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		ctx:     ctx,
		log:     logger,
		clients: &clients{},
		redis:   redis,
		channel: channel,
	}
}

// wsHandler is the http request handler for websocket connections
func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection.
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("error upgrading connection")
		return
	}
	defer c.Close()

	// Client has connected, save its connection in clients struct.
	// Remove the connection once the client is disconected.
	s.clients.add(c)
	defer s.clients.remove(c)

	s.log.Info("client connected", zap.Any("address", c.RemoteAddr()))

	// Wait for messages
	for {
		select {
		case <-s.ctx.Done():
			s.log.Info("closing context")
			return
		default:
			_, msg, err := c.ReadMessage()
			if err != nil {
				s.log.Error("error reading message", zap.Error(err))
				break
			}

			s.log.Info("message received", zap.String("msg", string(msg)))

			// Publish the message to redis
			err = s.redis.Publish(s.ctx, s.channel, msg).Err()
			if err != nil {
				s.log.Error("error publishing message to redis", zap.Error(err))
			}

			s.log.Info("message published")
		}
	}
}

// broadcast is a wrapper function for clients.broadcast() func
func (s *server) broadcast(msg string) []error {
	return s.clients.broadcast([]byte(msg))
}

// start will run a websocket server
func (s *server) start(addr string) {
	// Set http handlers
	http.HandleFunc("/", s.wsHandler)

	// Start server
	if err := http.ListenAndServe(addr, nil); err != nil {
		s.log.Fatal("error starting http server", zap.Error(err))
	}
}
