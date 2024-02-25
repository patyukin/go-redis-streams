package main

import (
	"context"
	"github.com/patyukin/go-redis-streams/internal/streamer/redis"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := redis.NewRedisStreamer(ctx)

	go client.Consume(ctx, "tickets")

	go client.LimitConsume(ctx, "tickets")
}
