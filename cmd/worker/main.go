package main

import (
	"context"
	"nanobit_assignment/internal/common"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func main() {
	// Create a new logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// Create new context with cancel func
	ctx, cancel := context.WithCancel(context.Background())

	// Load config
	redisAddr := common.GetEnv("REDIS_ADDR", "localhost:6379")
	redisChannel := common.GetEnv("REDIS_CHANNEL", "messages")

	// Create a new redis client
	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		panic(err)
	}
	r := redis.NewClient(opt)
	defer r.Close()

	// Subscribe to a channel
	rps := r.Subscribe(ctx, redisChannel)

	done := make(chan error, 1)

	// Message hander
	h := newHandler(ctx, r, logger)

	// Receive messages in a goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case n := <-rps.Channel():
				h.HandleMessage(n.Channel, n.Payload)
			}
		}
	}()

	// Wait for SIGINT
	select {
	case <-common.TrapSignals():
		cancel()
		return
	case <-done:
		return
	}
}
