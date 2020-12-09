package main

import (
	"context"
	"nanobit_assignment/internal/common"

	"github.com/go-redis/redis/v8"

	"go.uber.org/zap"
)

func main() {
	// Init logger
	logger, err := zap.NewProduction()
	if err != nil {
		// We don't want to continue with any processes
		// if the logger is not initiated correctly.
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	// create context with cancel func
	ctx, cancel := context.WithCancel(context.Background())

	// Load configuration
	addr := common.GetEnv("ADDR", ":8080")
	redisAddr := common.GetEnv("REDIS_ADDR", "localhost:6379")
	messageChan := common.GetEnv("REDIS_MESSAGES_CHANNEL", "messages")
	responseChan := common.GetEnv("REDIS_RESPONSES_CHANNEL", "responses")

	// Create a new redis client
	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		panic(err)
	}
	r := redis.NewClient(opt)
	defer r.Close()

	// Websocket server
	s := newServer(ctx, logger, r, messageChan)
	go s.start(addr)

	// Response handler
	h := newHandler(ctx, logger, s, r)
	go h.run(responseChan)

	// Wait for SIGINT
	select {
	case <-common.TrapSignals():
		cancel()
	}
}
