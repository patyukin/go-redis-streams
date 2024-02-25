package redis

import (
	"context"
	"fmt"
	"github.com/patyukin/go-redis-streams/internal/streamer"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
)

var _ streamer.StreamerInterface = (*Streamer)(nil)

type Streamer struct {
	c *redis.Client
}

func NewRedisStreamer(ctx context.Context) *Streamer {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", "127.0.0.1", "6379"),
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unbale to connect to Redis", err)
	}

	return &Streamer{
		c: client,
	}
}

func (s *Streamer) Stream() {
	// TODO
}

func (s *Streamer) Consume() {
	// TODO
}

func (s *Streamer) Publish(ctx context.Context) error {
	err := s.c.XAdd(ctx, &redis.XAddArgs{
		Stream: "tickets",
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			"whatHappened": string("ticket received"),
			"ticketID":     int(rand.Intn(100000000)),
			"ticketData":   string("some ticket data"),
		},
	}).Err()

	return err
}
