package main

import (
	"context"
	"encoding/json"
	"nanobit_assignment/internal/common"
	"sort"

	"github.com/go-redis/redis/v8"

	"go.uber.org/zap"
)

type handler struct {
	ctx   context.Context
	redis *redis.Client
	log   *zap.Logger
}

func newHandler(ctx context.Context, r *redis.Client, logger *zap.Logger) *handler {
	return &handler{
		ctx:   ctx,
		redis: r,
		log:   logger,
	}
}

func (h *handler) HandleMessage(channel string, msg string) {
	h.log.Info("messge received", zap.String("channel", channel), zap.String("msg", msg))

	// Decode message
	reqMsg := common.FavoriteNumber{}
	message := common.Command{
		Message: &reqMsg,
	}
	if err := json.Unmarshal([]byte(msg), &message); err != nil {
		h.log.Error("error unmarshaling message", zap.Error(err))
		return
	}

	// Save the number in redis
	if message.Action == "setNumber" {
		err := h.redis.Set(h.ctx, reqMsg.Username, reqMsg.Number, 0).Err()
		if err != nil {
			h.log.Error("error saving number", zap.Error(err))
		}

		h.log.Info("number saved")

		h.sendUsersList()
		return
	}

	// Send a list of all users
	if message.Action == "listUsers" {
		h.sendUsersList()
		return
	}

	h.log.Warn("invalid command")
}

func (h *handler) sendUsersList() {
	var cursor uint64

	keys, cursor, err := h.redis.Scan(h.ctx, cursor, "*", 0).Result()
	if err != nil {
		h.log.Error("error loading keys", zap.Error(err))
		return
	}

	// Sort alphabetically
	sort.Strings(keys)

	// Get numbers for all users and generate the response
	res := common.UsersList{}

	for _, key := range keys {
		// Load value from redis
		num, err := h.redis.Get(h.ctx, key).Int()
		if err != nil {
			h.log.Error("error loading key")
			return
		}

		res.Count++
		res.Users = append(res.Users, common.User{
			Username:       key,
			FavoriteNumber: num,
		})
	}

	// Encode the message
	msg, err := json.Marshal(res)
	if err != nil {
		h.log.Error("error marshaling messsage", zap.Error(err))
		return
	}

	// Send message to redis
	if err := h.redis.Publish(h.ctx, "responses", msg).Err(); err != nil {
		h.log.Error("error publishing message", zap.Error(err))
	}

	h.log.Info("response sent")
}
